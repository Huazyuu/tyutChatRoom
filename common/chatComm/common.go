package chatComm

import (
	"gin-gorilla/model/ctype"
	"github.com/gorilla/websocket"
	"time"
)

type ChatUser struct {
	Conn     *websocket.Conn
	UserID   string `json:"userid"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}
type ChatMessage struct {
	MsgType     ctype.MsgType `json:"msg_type"`
	Sender      ChatUser      `json:"sender"`
	Target      string        `json:"target"`
	Content     interface{}   `json:"content"`
	Date        time.Time     `json:"date"`         // 消息的时间
	OnlineCount int           `json:"online_count"` // 在线人数
}
