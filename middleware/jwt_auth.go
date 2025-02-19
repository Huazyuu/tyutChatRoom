package middleware

import (
	"gin-gorilla/res"
	"gin-gorilla/service/redisService"
	"gin-gorilla/utils/jwt"
	"github.com/gin-gonic/gin"
)

// JwtAuth jwt auth权限
func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		// global.Log.Debug("token", token)
		if token == "" {
			res.FailWithMessage("未携带token", c)
			c.Abort()
			return
		}
		claims, err := jwt.ParseToken(token)
		if err != nil {
			res.FailWithMessage("token错误", c)
			c.Abort()
			return
		}
		// 是否在redis中
		if redisService.CheckLogout(token) {
			res.FailWithMessage("token已失效", c)
			c.Abort()
			return
		}
		c.Set("claims", claims)
	}
}
