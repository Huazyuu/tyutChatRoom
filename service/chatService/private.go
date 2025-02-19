package chatService

import (
	"gin-gorilla/common/chatComm"
	"gin-gorilla/global"
	"gin-gorilla/model"
	"gin-gorilla/model/ctype"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"time"
)

// ChatPrivateService 处理群聊
func ChatPrivateService(conn *websocket.Conn, cu chatComm.ChatUser, targetid string) {
	var req chatComm.PrivateRequest

	for {
		_, content, err := conn.ReadMessage()
		if err != nil {
			break
		}
		if err := json.Unmarshal(content, &req); err != nil {
			sendPrivateSystemMsg(conn, "消息解析失败")
			continue
		}
		// 消息类型处理
		switch req.MsgType {
		case ctype.TextMsg:
			handlePrivateTextMessage(conn, cu, targetid, req)
		case ctype.FileMsg:
			// handlePrivateFileMessage(conn, userid, msgReq)
		default:
			_ = sendPrivateSystemMsg(conn, "不支持的消息类型")
		}
	}
}

// 发送私聊系统消息
func sendPrivateSystemMsg(conn *websocket.Conn, content string) error {
	response := chatComm.PrivateResponse{
		MsgType: ctype.SystemMsg,
		Content: content,
		Date:    time.Now(),
	}
	return sendPrivateMessage(conn, response)
}

// 通用私聊消息发送
func sendPrivateMessage(conn *websocket.Conn, response chatComm.PrivateResponse) error {
	byteData, _ := json.Marshal(response)

	_addr := conn.RemoteAddr().String()
	ip, port := getIPAndAddr(_addr)
	global.DB.Create(&model.ChatModel{
		UserID:   response.UserId,
		TargetID: response.TargetID,
		Content:  response.Content,
		IP:       ip,
		Addr:     port,
		IsGroup:  true,
		MsgType:  response.MsgType,
	})

	return conn.WriteMessage(websocket.TextMessage, byteData)
}

// 发送私聊回执系统消息
func sendBackSystemMsg(conn *websocket.Conn, content string) error {
	response := chatComm.PrivateResponse{
		MsgType: ctype.SystemMsg,
		Content: content,
		Date:    time.Now(),
	}
	return sendBackMessage(conn, response)
}

// 通用私聊回执消息发送
func sendBackMessage(conn *websocket.Conn, response chatComm.PrivateResponse) error {
	byteData, _ := json.Marshal(response)
	return conn.WriteMessage(websocket.TextMessage, byteData)
}
