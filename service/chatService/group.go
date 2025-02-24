package chatService

import (
	"encoding/json"
	"fmt"
	"gin-gorilla/common/chatComm"
	"gin-gorilla/global"
	"gin-gorilla/model"
	"gin-gorilla/model/ctype"
	"github.com/gorilla/websocket"
	"strings"
	"time"
)

// ChatGroupService 处理群聊
func ChatGroupService(conn *websocket.Conn, cu chatComm.ChatUser) {
	for {
		// 读取聊天数据 json格式
		_, content, err := conn.ReadMessage()
		if err != nil {
			// 断开连接
			sendGroupMsg(conn, chatComm.ChatMessage{
				MsgType: ctype.OutRoomMsg,
				Sender: chatComm.ChatUser{
					Conn:     nil,
					UserID:   cu.UserID,
					Username: cu.Username,
					Avatar:   cu.Avatar,
				},
				Target:      "global_user",
				Content:     fmt.Sprintf("[%s]离开聊天室", cu.Username),
				Date:        time.Now(),
				OnlineCount: len(chatComm.ConnGroupMap) - 1,
			})
			break
		}
		// json数据序列化request结构体
		var req chatComm.GroupRequest
		err = json.Unmarshal(content, &req)
		if err != nil {
			// 参数绑定失败
			sendMsg(conn, cu.UserID, chatComm.ChatMessage{
				MsgType: ctype.SystemMsg,
				Sender: chatComm.ChatUser{
					Conn:     nil,
					UserID:   cu.UserID,
					Username: cu.Username,
					Avatar:   cu.Avatar,
				},
				Content:     "参数绑定失败",
				Date:        time.Now(),
				OnlineCount: len(chatComm.ConnGroupMap),
			})
			continue
		}
		// 断言
		switch req.MsgType {
		case ctype.TextMsg:
			// 内容为空
			if strings.TrimSpace(req.Content) == "" {
				sendMsg(conn, cu.UserID, chatComm.ChatMessage{
					MsgType: ctype.SystemMsg,
					Sender: chatComm.ChatUser{
						Conn:     nil,
						UserID:   cu.UserID,
						Username: cu.Username,
						Avatar:   cu.Avatar,
					},
					Content:     "消息不能为空",
					Date:        time.Now(),
					OnlineCount: len(chatComm.ConnGroupMap),
				})
				continue
			}
			sendGroupMsg(conn, chatComm.ChatMessage{
				MsgType: ctype.SystemMsg,
				Sender: chatComm.ChatUser{
					Conn:     nil,
					UserID:   cu.UserID,
					Username: cu.Username,
					Avatar:   cu.Avatar,
				},
				Content:     req.Content,
				Date:        time.Now(),
				OnlineCount: len(chatComm.ConnGroupMap),
			})
		case ctype.InRoomMsg:
			sendGroupMsg(conn, chatComm.ChatMessage{
				MsgType: ctype.SystemMsg,
				Sender: chatComm.ChatUser{
					Conn:     nil,
					UserID:   cu.UserID,
					Username: cu.Username,
					Avatar:   cu.Avatar,
				},
				Target:      "global_user",
				Content:     fmt.Sprintf("[%s]进入聊天室", cu.Username),
				Date:        time.Now(),
				OnlineCount: len(chatComm.ConnGroupMap),
			})
		default:
			sendMsg(conn, cu.UserID, chatComm.ChatMessage{
				MsgType: ctype.SystemMsg,
				Sender: chatComm.ChatUser{
					Conn:     nil,
					UserID:   cu.UserID,
					Username: cu.Username,
					Avatar:   cu.Avatar,
				},
				Content:     "参数类型错误",
				Date:        time.Now(),
				OnlineCount: len(chatComm.ConnGroupMap),
			})
		}

	}
}

// sendGroupMsg 群聊功能
func sendGroupMsg(conn *websocket.Conn, response chatComm.ChatMessage) {
	byteData, _ := json.Marshal(response)
	_addr := conn.RemoteAddr().String()
	ip, port := getIPAndAddr(_addr)

	// global.Log.Debug(response)

	global.DB.Create(&model.ChatModel{
		UserID:   response.Sender.UserID,
		TargetID: "globaluser",
		Content:  response.Content.(string),
		IP:       ip,
		Addr:     port,
		IsGroup:  true,
		MsgType:  response.MsgType,
	})
	for _, chatUser := range chatComm.ConnGroupMap {
		chatUser.Conn.WriteMessage(websocket.TextMessage, byteData)
	}
}

// sendMsg sys消息
func sendMsg(conn *websocket.Conn, userid string, response chatComm.ChatMessage) {
	byteData, _ := json.Marshal(response)
	chatUser := chatComm.ConnGroupMap[userid]

	_addr := conn.RemoteAddr().String()
	ip, port := getIPAndAddr(_addr)

	global.DB.Create(&model.ChatModel{
		UserID:   response.Sender.UserID,
		TargetID: response.Sender.UserID,
		Content:  response.Content.(string),
		IP:       ip,
		Addr:     port,
		IsGroup:  false,
		MsgType:  response.MsgType,
	})
	chatUser.Conn.WriteMessage(websocket.TextMessage, byteData)
}

func getIPAndAddr(_addr string) (ip string, addr string) {
	addrList := strings.Split(_addr, ":")
	addr = "内网地址"
	return addrList[0], addr
}
