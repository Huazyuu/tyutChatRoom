package session

import (
	"ginchat/conf"
	"ginchat/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// EnableCookieSession session中间件
func EnableCookieSession() gin.HandlerFunc {
	store := cookie.NewStore([]byte(conf.GlobalConf.App.CookieKey))
	return sessions.Sessions("go-gin-chat", store)
}

// SaveAuthSession 保存session信息
func SaveAuthSession(c *gin.Context, info any) {
	session := sessions.Default(c)
	session.Set("uid", info)
	session.Save()
}

// GetSessionUserInfo 用户存储在session中的信息
func GetSessionUserInfo(c *gin.Context) map[string]any {
	session := sessions.Default(c)
	uid := session.Get("uid")
	data := make(map[string]any)
	var uidStr string

	// 检查 uid 的类型
	switch v := uid.(type) {
	case string:
		uidStr = v
	case uint:
		uidStr = strconv.FormatUint(uint64(v), 10)
	default:
		return data
	}

	if uidStr != "" {
		user := models.FindUserByField("id", uidStr)
		data["uid"] = user.ID
		data["username"] = user.Username
		data["avatar_id"] = user.AvatarId
	}
	return data
}

// ClearAuthSession 清除session
func ClearAuthSession(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
}

// HasSession Session检查
func HasSession(c *gin.Context) bool {
	session := sessions.Default(c)
	if sessionValue := session.Get("uid"); sessionValue == nil {
		zap.L().Info("没有session")
		return false
	}
	return true
}

// AuthSessionMiddle 鉴权中间件
func AuthSessionMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		sessionValue := session.Get("uid")
		if sessionValue == nil {
			zap.L().Info("session缺失")
			c.Redirect(http.StatusFound, "/")
			return
		}

		var uidInt int
		var err error

		// 检查 sessionValue 的类型
		switch v := sessionValue.(type) {
		case string:
			uidInt, err = strconv.Atoi(v)
		case uint:
			uidInt = int(v)
		default:
			c.Redirect(http.StatusFound, "/")
			return
		}

		if err != nil || uidInt <= 0 {
			c.Redirect(http.StatusFound, "/")
			return
		}

		// 设置简单的变量
		c.Set("uid", sessionValue)

		c.Next()
		return
	}
}
