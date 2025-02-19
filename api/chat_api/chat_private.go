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
	"sync"
)

// 使用Mutex来保护ConnPrivateMap
var mu sync.Mutex

// ChatPrivateView 私聊接口
func (ChatApi) ChatPrivateView(c *gin.Context) {
	_claims, _ := c.Get("claims")
	claims := _claims.(*jwt.CustomClaims)
	// http 升级 websockets protocol
	conn, err := wsService.WSUpgarde(claims, c)
	if err != nil {
		global.Log.Errorf("WebSocket升级失败: %v", err)
		res.FailWithCode(res.ArgumentError, c)
		return
	}

	// user info
	username := claims.Username
	userid := claims.UserID
	addr := conn.RemoteAddr().String()
	avatar := model.SelectAvatar(userid)
	// add user
	chatUser := chatComm.ChatUser{
		Conn:     conn,
		UserID:   userid,
		Username: username,
		Avatar:   avatar,
	}

	// target info
	targetid := c.Query("target_id")
	if targetid == "" {
		global.Log.Info("对方id不正确")
		res.FailWithMessage("对方id不正确", c)
		return
	}

	roomID := generateRoomID(userid, targetid)

	// 保护ConnPrivateMap的并发写入
	mu.Lock()
	chatComm.ConnPrivateMap[roomID] = append(chatComm.ConnPrivateMap[roomID], chatUser)
	global.Log.Infof("[%s][%s]连接成功", addr, username)
	mu.Unlock()

	// 处理私聊消息
	go func() {
		defer func() {
			// 断开连接时清理
			mu.Lock()
			users := chatComm.ConnPrivateMap[roomID]
			for i, u := range users {
				if u.UserID == userid {
					// 删除当前用户
					chatComm.ConnPrivateMap[roomID] = append(users[:i], users[i+1:]...)
					break
				}
			}
			if len(chatComm.ConnPrivateMap[roomID]) == 0 {
				delete(chatComm.ConnPrivateMap, roomID)
			}
			mu.Unlock()

			conn.Close()
		}()
		// 调用聊天服务处理私聊消息
		chatService.ChatPrivateService(conn, chatUser, targetid)
	}()
}

// generateRoomID 生成唯一的房间ID，确保用户1和用户2无论连接顺序如何，房间ID一致
func generateRoomID(user1, user2 string) string {
	if user1 < user2 {
		return user1 + "_" + user2
	}
	return user2 + "_" + user1
}
