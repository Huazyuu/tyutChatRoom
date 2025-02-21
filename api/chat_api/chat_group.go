package chat_api

import (
	"gin-gorilla/common/chatComm"
	"gin-gorilla/global"
	"gin-gorilla/model"
	"gin-gorilla/res"
	"gin-gorilla/service/chatService"
	"gin-gorilla/service/wsService"
	"gin-gorilla/utils/jwt"
	"github.com/gin-gonic/gin"
)

// ChatGroupView 群聊接口
func (ChatApi) ChatGroupView(c *gin.Context) {
	_claims, exists := c.Get("claims")
	if !exists {
		res.FailWithMessage("未通过 JWT 验证", c)
		return
	}
	claims := _claims.(*jwt.CustomClaims)
	// http 升级 websockets protocol
	conn, err := wsService.WSUpgarde(claims, c)
	if err != nil {
		global.Log.Error("WS Upgarder 失败", err.Error())
		res.FailWithCode(res.ArgumentError, c)
		return
	}
	// JWT 获取 user info
	username := claims.Username
	userid := claims.UserID
	addr := conn.RemoteAddr().String()
	avatar := model.SelectAvatar(userid)

	chatUser := chatComm.ChatUser{
		Conn:     conn,
		UserID:   userid,
		Username: username,
		Avatar:   avatar,
	}
	mu.Lock()
	chatComm.ConnGroupMap[userid] = chatUser
	global.Log.Infof("[%s][%s]连接成功", addr, username)
	mu.Unlock()

	// 处理数据
	chatService.ChatGroupService(conn, chatUser)
	// 断开连接
	defer conn.Close()
	delete(chatComm.ConnGroupMap, userid)

}
