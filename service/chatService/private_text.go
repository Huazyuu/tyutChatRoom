package chatService

import (
	"gin-gorilla/common/chatComm"
	"gin-gorilla/global"
	"gin-gorilla/model"
	"github.com/gorilla/websocket"
	"strings"
	"sync"
	"time"
)

var mu sync.Mutex

// handlePrivateTextMessage 处理私聊文本消息
func handlePrivateTextMessage(conn *websocket.Conn, cu chatComm.ChatUser, targetid string, req chatComm.PrivateRequest) {
	if strings.TrimSpace(req.Content) == "" {
		global.Log.Info("请求内容为空")
		_ = sendBackSystemMsg(conn, "消息内容不能为空")
		return
	}
	// 根据当前用户和目标用户生成房间ID
	roomID := generateRoomID(cu.UserID, targetid)

	// 获取接收方房间内的用户列表
	mu.Lock()
	usersInRoom, exists := chatComm.ConnPrivateMap[roomID]
	mu.Unlock()
	// 如果房间不存在或者接收方不在线
	if !exists {
		_ = sendBackSystemMsg(conn, "目标用户不在线或房间不存在")
		return
	}
	// 在房间中找到目标用户的连接
	var receiver *websocket.Conn
	for _, user := range usersInRoom {
		if user.UserID == targetid {
			receiver = user.Conn
			break
		}
	}
	// 如果目标用户未找到
	if receiver == nil {
		_ = sendBackSystemMsg(conn, "目标用户不在线")
		return
	}

	targetName := model.SelectUsername(targetid)
	targetAvatar := model.SelectAvatar(targetid)
	// 发送私聊消息到目标用户
	err := sendPrivateMessage(receiver, chatComm.PrivateResponse{
		UserId:       cu.UserID,
		Username:     cu.Username,
		Avatar:       cu.Avatar,
		TargetID:     targetid,
		TargetName:   targetName,
		TargetAvatar: targetAvatar,
		MsgType:      req.MsgType,
		Content:      req.Content,
		Date:         time.Now(),
		OnlineCount:  len(usersInRoom),
	})
	if err != nil {
		_ = sendBackSystemMsg(conn, "发送消息失败")
	}
	// 回执
	err = sendBackMessage(conn, chatComm.PrivateResponse{
		UserId:       cu.UserID,
		Username:     cu.Username,
		Avatar:       cu.Avatar,
		TargetID:     targetid,
		TargetName:   targetName,
		TargetAvatar: targetAvatar,
		MsgType:      req.MsgType,
		Content:      req.Content,
		Date:         time.Now(),
		OnlineCount:  len(usersInRoom),
	})
	if err != nil {
		_ = sendBackSystemMsg(conn, "发送消息失败")
	}

}

// generateRoomID 生成唯一的房间ID，确保用户1和用户2无论连接顺序如何，房间ID一致
func generateRoomID(user1, user2 string) string {
	if user1 < user2 {
		return user1 + "_" + user2
	}
	return user2 + "_" + user1
}
