package model

import (
	"gin-gorilla/model/ctype"
	"time"
)

// ChatModel 表示聊天记录模型
type ChatModel struct {
	ID        uint      `json:"id" gorm:"column:id"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`

	UserID   string `json:"user_id" gorm:"column:user_id"`
	TargetID string `json:"target_id" gorm:"column:target_id"`

	// 聊天内容
	Content string        `json:"content" gorm:"column:content"`
	IP      string        `json:"ip" gorm:"column:ip"`
	Addr    string        `json:"addr" gorm:"column:addr"`
	IsGroup bool          `json:"is_group" gorm:"column:is_group"`
	MsgType ctype.MsgType `json:"msg_type" gorm:"column:msg_type"`
}
