package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Cors 自定义的跨域中间件
func Cors(allowedOrigins []string) gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		origin := context.Request.Header.Get("Origin")

		// 检查请求的源是否在允许的源列表中
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if allowedOrigin == origin {
				allowed = true
				break
			}
		}

		if allowed {
			context.Header("Access-Control-Allow-Origin", origin)
			context.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
			context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
			context.Header("Access-Control-Allow-Credentials", "true")
		}

		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
			return
		}

		context.Next()
	}
}
