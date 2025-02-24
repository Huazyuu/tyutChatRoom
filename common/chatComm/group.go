package chatComm

import (
	"gin-gorilla/model/ctype"
)

// ConnGroupMap 用户map
var ConnGroupMap = make(map[string]ChatUser)

// GroupRequest 请求
type GroupRequest struct {
	Content string        `json:"content"`  // 聊天的内容
	MsgType ctype.MsgType `json:"msg_type"` // 聊天类型
}

// ChatListRequest 列表
type ChatListRequest struct {
	Username string `form:"username"`
	Page     int    `form:"page"`
	Limit    int    `form:"limit"`
	Sort     string `form:"sort"`
}
