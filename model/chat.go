package model

import (
	"errors"
	"gin-gorilla/global"
	"gin-gorilla/model/ctype"
	"gorm.io/gorm"
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

func SelectUsername(userid string) (username string) {
	var userModel UserModel
	err := global.DB.Select("username").Where("user_id = ?", userid).First(&userModel).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		global.Log.Error(err.Error())
		return
	}
	username = userModel.UserName
	return
}
func SelectAvatar(userid string) (avatar string) {
	var userModel UserModel
	err := global.DB.Select("avatar").Where("user_id = ?", userid).First(&userModel).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		global.Log.Info("没有头像,将使用默认头像")
	}
	avatar = userModel.Avatar
	return
}
