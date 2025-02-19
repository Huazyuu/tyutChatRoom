package chatService

import (
	"fmt"
	"gin-gorilla/api/chat_api/chatComm"
	"gin-gorilla/global"
	"gin-gorilla/model"
	"gin-gorilla/model/ctype"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"strings"
	"time"
)

// ChatgroupService 处理群聊
func ChatgroupService(conn *websocket.Conn, chatUser chatComm.ChatUser) {
	for {
		// 读取聊天数据 json格式
		_, content, err := conn.ReadMessage()
		if err != nil {
			// 断开连接
			sendGroupMsg(conn, chatComm.GroupResponse{
				UserId:      chatUser.UserID,
				Username:    chatUser.Username,
				Avatar:      chatUser.Avatar,
				Content:     fmt.Sprintf("[%s]离开聊天室", chatUser.Username),
				Date:        time.Now(),
				MsgType:     ctype.OutRoomMsg,
				OnlineCount: len(chatComm.ConnGroupMap) - 1,
			})
			break
		}
		// json数据序列化request结构体
		var req chatComm.GroupRequest
		err = json.Unmarshal(content, &req)
		if err != nil {
			// 参数绑定失败
			sendMsg(conn, chatUser.UserID, chatComm.GroupResponse{
				UserId:      chatUser.UserID,
				Username:    chatUser.Username,
				Avatar:      chatUser.Avatar,
				MsgType:     ctype.SystemMsg,
				Date:        time.Now(),
				Content:     "参数绑定失败",
				OnlineCount: len(chatComm.ConnGroupMap),
			})
			continue
		}
		// 断言
		switch req.MsgType {
		case ctype.TextMsg:
			// 内容为空
			if strings.TrimSpace(req.Content) == "" {
				sendMsg(conn, chatUser.UserID, chatComm.GroupResponse{
					UserId:      chatUser.UserID,
					Username:    chatUser.Username,
					Avatar:      chatUser.Avatar,
					MsgType:     ctype.SystemMsg,
					Content:     "消息不能为空",
					OnlineCount: len(chatComm.ConnGroupMap),
				})
				continue
			}
			sendGroupMsg(conn, chatComm.GroupResponse{
				UserId:      chatUser.UserID,
				Username:    chatUser.Username,
				Avatar:      chatUser.Avatar,
				Content:     req.Content,
				MsgType:     ctype.TextMsg,
				Date:        time.Now(),
				OnlineCount: len(chatComm.ConnGroupMap),
			})
		case ctype.InRoomMsg:
			sendGroupMsg(conn, chatComm.GroupResponse{
				UserId:      chatUser.UserID,
				Username:    chatUser.Username,
				Avatar:      chatUser.Avatar,
				Content:     fmt.Sprintf("[%s]进入聊天室", chatUser.Username),
				Date:        time.Now(),
				OnlineCount: len(chatComm.ConnGroupMap),
			})
		default:
			sendMsg(conn, chatUser.UserID, chatComm.GroupResponse{
				UserId:      chatUser.UserID,
				Username:    chatUser.Username,
				Avatar:      chatUser.Avatar,
				MsgType:     ctype.SystemMsg,
				Content:     "消息类型错误",
				OnlineCount: len(chatComm.ConnGroupMap),
			})
		}

	}
}

// sendGroupMsg 群聊功能
func sendGroupMsg(conn *websocket.Conn, response chatComm.GroupResponse) {
	byteData, _ := json.Marshal(response)
	_addr := conn.RemoteAddr().String()
	ip, addr := getIPAndAddr(_addr)

	global.DB.Create(&model.ChatModel{
		UserID:   response.UserId,
		TargetID: response.UserId,
		Content:  response.Content,
		IP:       ip,
		Addr:     addr,
		IsGroup:  true,
		MsgType:  response.MsgType,
	})
	for _, chatUser := range chatComm.ConnGroupMap {
		chatUser.Conn.WriteMessage(websocket.TextMessage, byteData)
	}
}

// sendMsg 给某个用户发消息
func sendMsg(conn *websocket.Conn, userid string, response chatComm.GroupResponse) {
	byteData, _ := json.Marshal(response)
	chatUser := chatComm.ConnGroupMap[userid]

	_addr := conn.RemoteAddr().String()
	ip, port := getIPAndAddr(_addr)

	global.DB.Create(&model.ChatModel{
		UserID:  response.Username,
		Content: response.Content,
		IP:      ip,
		Addr:    port,
		IsGroup: false,
		MsgType: response.MsgType,
	})
	chatUser.Conn.WriteMessage(websocket.TextMessage, byteData)
}

func getIPAndAddr(_addr string) (ip string, addr string) {
	addrList := strings.Split(_addr, ":")
	addr = "内网地址"
	return addrList[0], addr
}
