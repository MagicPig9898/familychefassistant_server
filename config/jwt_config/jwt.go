package jwt_config

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var salt = []byte("abcd")

// Claims 自定义载荷
type Claims struct {
	OpenID string `json:"open_id"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 token，传入 openid 和过期时长
func GenerateToken(openID string, expireDuration time.Duration) (string, error) {
	claims := Claims{
		OpenID: openID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)), // token 的过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                     //token 的签发时间
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(salt)
}

// ParseToken 解析 token，返回 Claims
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return salt, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
