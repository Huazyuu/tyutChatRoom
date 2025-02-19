package jwt

import (
	"errors"
	"fmt"
	"gin-gorilla/global"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/sirupsen/logrus"
)

// ParseToken 解析 token
func ParseToken(tokenStr string) (*CustomClaims, error) {
	MySecret := []byte(global.Config.Jwt.Secret)
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return MySecret, nil
	})
	if err != nil {
		logrus.Error(fmt.Sprintf("token parse err: %s", err.Error()))
		return nil, errors.New(fmt.Sprintf("token parse err: %s", err.Error()))
	}
	// 类型断言
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
