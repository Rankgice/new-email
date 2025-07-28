package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// UserClaims 用户JWT声明
type UserClaims struct {
	UserId   uint   `json:"userId"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// AdminClaims 管理员JWT声明
type AdminClaims struct {
	AdminId  uint   `json:"adminId"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateUserToken 生成用户Token
func GenerateUserToken(userId uint, username string, secret string, expireHours int) (string, error) {
	claims := UserClaims{
		UserId:   userId,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "email-system",
			Subject:   "user",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// GenerateAdminToken 生成管理员Token
func GenerateAdminToken(adminId uint, username, role string, secret string, expireHours int) (string, error) {
	claims := AdminClaims{
		AdminId:  adminId,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "email-system",
			Subject:   "admin",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken 解析用户Token
func ParseToken(tokenString string, secret string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ParseAdminToken 解析管理员Token
func ParseAdminToken(tokenString string, secret string) (*AdminClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AdminClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshToken 刷新Token
func RefreshToken(tokenString string, secret string, expireHours int) (string, error) {
	// 解析旧token（忽略过期错误）
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		// 检查是否只是过期错误
		if errors.Is(err, jwt.ErrTokenExpired) {
			// Token过期，但其他部分有效，可以刷新
			if claims, ok := token.Claims.(*UserClaims); ok {
				return GenerateUserToken(claims.UserId, claims.Username, secret, expireHours)
			}
		}
		return "", err
	}

	// Token有效，生成新token
	if claims, ok := token.Claims.(*UserClaims); ok {
		return GenerateUserToken(claims.UserId, claims.Username, secret, expireHours)
	}

	return "", errors.New("invalid token claims")
}

// RefreshAdminToken 刷新管理员Token
func RefreshAdminToken(tokenString string, secret string, expireHours int) (string, error) {
	// 解析旧token（忽略过期错误）
	token, err := jwt.ParseWithClaims(tokenString, &AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		// 检查是否只是过期错误
		if errors.Is(err, jwt.ErrTokenExpired) {
			// Token过期，但其他部分有效，可以刷新
			if claims, ok := token.Claims.(*AdminClaims); ok {
				return GenerateAdminToken(claims.AdminId, claims.Username, claims.Role, secret, expireHours)
			}
		}
		return "", err
	}

	// Token有效，生成新token
	if claims, ok := token.Claims.(*AdminClaims); ok {
		return GenerateAdminToken(claims.AdminId, claims.Username, claims.Role, secret, expireHours)
	}

	return "", errors.New("invalid token claims")
}

// ValidateToken 验证Token是否有效
func ValidateToken(tokenString string, secret string) bool {
	_, err := ParseToken(tokenString, secret)
	return err == nil
}

// ValidateAdminToken 验证管理员Token是否有效
func ValidateAdminToken(tokenString string, secret string) bool {
	_, err := ParseAdminToken(tokenString, secret)
	return err == nil
}

// GetTokenExpireTime 获取Token过期时间
func GetTokenExpireTime(tokenString string, secret string) (*time.Time, error) {
	claims, err := ParseToken(tokenString, secret)
	if err != nil {
		return nil, err
	}

	if claims.ExpiresAt != nil {
		expireTime := claims.ExpiresAt.Time
		return &expireTime, nil
	}

	return nil, errors.New("no expire time found")
}

// IsTokenExpired 检查Token是否过期
func IsTokenExpired(tokenString string, secret string) bool {
	expireTime, err := GetTokenExpireTime(tokenString, secret)
	if err != nil {
		return true
	}

	return expireTime.Before(time.Now())
}
