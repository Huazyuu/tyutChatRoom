package jwt

import (
	"gin-gorilla/global"
	"github.com/dgrijalva/jwt-go/v4"

	"time"
)

// GenToken 创建 Token
func GenToken(user JwtPayLoad) (string, error) {
	MySecret := []byte(global.Config.Jwt.Secret)
	claim := CustomClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(time.Hour * time.Duration(global.Config.Jwt.Expires))),
			Issuer:    global.Config.Jwt.Issuer, // 签发人
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(MySecret)
}
