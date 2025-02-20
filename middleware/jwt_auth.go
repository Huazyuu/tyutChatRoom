package middleware

import (
	"gin-gorilla/res"
	"gin-gorilla/service/redisService"
	"gin-gorilla/utils/jwt"
	"github.com/gin-gonic/gin"
	"strings"
)

// JwtAuth jwt auth权限
func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			res.FailWithMessage("未携带token", c)
			c.Abort()
			return
		}
		// 去掉Bearer prefix
		// 去掉Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			res.FailWithMessage("token格式错误", c)
			c.Abort()
			return
		}
		tokenString := parts[1]

		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			res.FailWithMessage("token错误", c)
			c.Abort()
			return
		}
		// 是否在redis中
		if redisService.CheckLogout(tokenString) {
			res.FailWithMessage("token已失效", c)
			c.Abort()
			return
		}
		c.Set("claims", claims)
	}
}
