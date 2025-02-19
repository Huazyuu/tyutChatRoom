package jwt

import (
	"github.com/dgrijalva/jwt-go/v4"
)

// JwtPayLoad jwt中payload数据
type JwtPayLoad struct {
	Username string `json:"username"` // 用户名
	UserID   string `json:"user_id"`  // 用户id
}

type CustomClaims struct {
	JwtPayLoad
	jwt.StandardClaims
}
