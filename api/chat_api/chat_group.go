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
	_claims, _ := c.Get("claims")
	claims := _claims.(*jwt.CustomClaims)
	// http 升级 websockets protocol
	conn, err := wsService.WSUpgarde(claims, c)
	if err != nil {
		res.FailWithCode(res.ArgumentError, c)
		return
	}
	// JWT 获取username userid
	username := claims.Username
	userid := claims.UserID
	// addr
	addr := conn.RemoteAddr().String()
	// avatar
	avatar := model.SelectAvatar(userid)
	// add user
	chatUser := chatComm.ChatUser{
		Conn:     conn,
		UserID:   claims.UserID,
		Username: claims.Username,
		Avatar:   avatar,
	}
	chatComm.ConnGroupMap[userid] = chatUser
	global.Log.Infof("[%s][%s]连接成功", addr, username)
	// 处理数据
	chatService.ChatgroupService(conn, chatUser)
	// 断开连接
	defer conn.Close()
	delete(chatComm.ConnGroupMap, userid)
}
