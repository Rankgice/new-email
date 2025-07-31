package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims 统一JWT声明结构体
type Claims struct {
	UserId   int64  `json:"userId"`   // 用户ID，当isAdmin=true时代表adminId
	Username string `json:"username"` // 用户名
	IsAdmin  bool   `json:"isAdmin"`  // 是否为管理员
	Role     string `json:"role"`     // 角色（仅管理员使用）
	jwt.RegisteredClaims
}

// GenerateToken 生成Token（用户版本）
func GenerateToken(userId int64, userType string, secret string, expireHours int) (string, error) {
	isAdmin := userType == "admin"
	username := ""
	role := ""

	return GenerateTokenFull(userId, username, isAdmin, role, secret, expireHours)
}

// GenerateTokenFull 生成完整Token
func GenerateTokenFull(userId int64, username string, isAdmin bool, role string, secret string, expireHours int) (string, error) {
	subject := "user"
	if isAdmin {
		subject = "admin"
	}

	claims := Claims{
		UserId:   userId,
		Username: username,
		IsAdmin:  isAdmin,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "email-system",
			Subject:   subject,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken 解析Token（统一解析用户和管理员Token）
// 注意：此方法会忽略过期错误，需要在业务层处理过期逻辑
func ParseToken(tokenString string, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	// 如果是过期错误，但token其他部分有效，仍然返回claims
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			if claims, ok := token.Claims.(*Claims); ok {
				return claims, err // 返回claims和过期错误
			}
		}
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
