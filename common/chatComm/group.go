package chatComm

import (
	"gin-gorilla/model/ctype"
	"time"
)

// ConnGroupMap 用户map
var ConnGroupMap = make(map[string]ChatUser)

// GroupRequest 请求
type GroupRequest struct {
	Content string        `json:"content"`  // 聊天的内容
	MsgType ctype.MsgType `json:"msg_type"` // 聊天类型
}

// GroupResponse 响应
type GroupResponse struct {
	UserId      string        `json:"user_id"`
	Username    string        `json:"username"`
	Avatar      string        `json:"avatar"`       // 头像
	MsgType     ctype.MsgType `json:"msg_type"`     // 聊天类型
	Content     string        `json:"content"`      // 聊天的内容
	Date        time.Time     `json:"date"`         // 消息的时间
	OnlineCount int           `json:"online_count"` // 在线人数
}

// ChatListRequest 列表
type ChatListRequest struct {
	Username string `form:"username"`
	Page     int    `form:"page"`
	Limit    int    `form:"limit"`
	Sort     string `form:"sort"`
}
