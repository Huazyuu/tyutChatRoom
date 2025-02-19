package chat_api

import (
	"gin-gorilla/api/chat_api/chatComm"
	"gin-gorilla/global"
	"gin-gorilla/model"
	"gin-gorilla/res"
	"gin-gorilla/service/chatService"
	"gin-gorilla/utils/jwt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

// ChatGroupView 群聊接口
func (ChatApi) ChatGroupView(c *gin.Context) {
	_claims, _ := c.Get("claims")
	claims := _claims.(*jwt.CustomClaims)
	// http 升级 websockets protocol
	var upGrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// 鉴权
			if claims == nil {
				return false
			}
			return true
		},
	}
	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
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
