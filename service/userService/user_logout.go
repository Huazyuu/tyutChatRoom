package userService

import (
	"gin-gorilla/service/redisService"
	"gin-gorilla/utils/jwt"
	"time"
)

func (UserService) Logout(claims *jwt.CustomClaims, token string) error {
	exp := claims.ExpiresAt
	now := time.Now()
	diff := exp.Time.Sub(now)
	return redisService.Logout(token, diff)
}
