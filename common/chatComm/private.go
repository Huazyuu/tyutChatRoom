package chatComm

import (
	"gin-gorilla/model/ctype"
	"time"
)

// ConnPrivateMap 用户map
var ConnPrivateMap = make(map[string][]ChatUser)

type PrivateRequest struct {
	MsgType ctype.MsgType `json:"msg_type"`       // 消息类型
	Content string        `json:"content"`        // 消息内容
	File    *File         `json:"file,omitempty"` // 文件上传
}

// PrivateResponse 响应
type PrivateResponse struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"` // 头像

	TargetID     string `json:"target_id"`
	TargetName   string `json:"target_name"`
	TargetAvatar string `json:"target_avatar"`

	MsgType     ctype.MsgType `json:"msg_type"`     // 聊天类型
	Content     string        `json:"content"`      // 聊天的内容
	Date        time.Time     `json:"date"`         // 消息的时间
	OnlineCount int           `json:"online_count"` // 在线人数
}
