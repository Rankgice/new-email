package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// PasswordConfig 密码配置
type PasswordConfig struct {
	Memory      uint32 // 内存使用量（KB）
	Iterations  uint32 // 迭代次数
	Parallelism uint8  // 并行度
	SaltLength  uint32 // 盐长度
	KeyLength   uint32 // 密钥长度
}

// DefaultPasswordConfig 默认密码配置
var DefaultPasswordConfig = &PasswordConfig{
	Memory:      64 * 1024, // 64MB
	Iterations:  3,
	Parallelism: 2,
	SaltLength:  16,
	KeyLength:   32,
}

// HashPassword 加密密码
func HashPassword(password string) (string, error) {
	return HashPasswordWithConfig(password, DefaultPasswordConfig)
}

// HashPasswordWithConfig 使用指定配置加密密码
func HashPasswordWithConfig(password string, config *PasswordConfig) (string, error) {
	// 生成随机盐
	salt, err := generateRandomBytes(config.SaltLength)
	if err != nil {
		return "", err
	}

	// 使用Argon2id生成哈希
	hash := argon2.IDKey([]byte(password), salt, config.Iterations, config.Memory, config.Parallelism, config.KeyLength)

	// 编码为base64
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// 格式：$argon2id$v=19$m=65536,t=3,p=2$salt$hash
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, config.Memory, config.Iterations, config.Parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

// VerifyPassword 验证密码
func VerifyPassword(password, encodedHash string) (bool, error) {
	// 解析编码的哈希
	config, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// 使用相同参数计算哈希
	otherHash := argon2.IDKey([]byte(password), salt, config.Iterations, config.Memory, config.Parallelism, config.KeyLength)

	// 使用恒定时间比较防止时序攻击
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}

	return false, nil
}

// generateRandomBytes 生成随机字节
func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// decodeHash 解码哈希字符串
func decodeHash(encodedHash string) (*PasswordConfig, []byte, []byte, error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, errors.New("invalid hash format")
	}

	var version int
	_, err := fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, errors.New("incompatible version of argon2")
	}

	config := &PasswordConfig{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &config.Memory, &config.Iterations, &config.Parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	config.SaltLength = uint32(len(salt))

	hash, err := base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	config.KeyLength = uint32(len(hash))

	return config, salt, hash, nil
}

// GenerateRandomPassword 生成随机密码
func GenerateRandomPassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"

	b := make([]byte, length)
	for i := range b {
		randomBytes := make([]byte, 1)
		_, err := rand.Read(randomBytes)
		if err != nil {
			return "", err
		}
		b[i] = charset[randomBytes[0]%byte(len(charset))]
	}

	return string(b), nil
}

// ValidatePasswordStrength 验证密码强度
func ValidatePasswordStrength(password string) error {
	if len(password) < 6 {
		return errors.New("密码长度至少6位")
	}

	if len(password) > 128 {
		return errors.New("密码长度不能超过128位")
	}

	// 可以添加更多密码强度检查
	// 例如：必须包含大小写字母、数字、特殊字符等

	return nil
}

// SimpleHashPassword 简单密码加密（用于兼容旧系统）
func SimpleHashPassword(password string) string {
	// 注意：这是一个简化版本，生产环境建议使用上面的HashPassword
	return fmt.Sprintf("%x", []byte(password))
}

// SimpleVerifyPassword 简单密码验证（用于兼容旧系统）
func SimpleVerifyPassword(password, hash string) bool {
	return SimpleHashPassword(password) == hash
}
