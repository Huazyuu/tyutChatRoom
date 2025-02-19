package model

import (
	"crypto/rand"
	"encoding/hex"
	"gorm.io/gorm"

	"time"
)

type UserModel struct {
	ID        uint      `json:"id,omitempty" gorm:"column:id"`
	CreatedAt time.Time `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" gorm:"column:updated_at"`

	UserID   string `json:"user_id" gorm:"column:user_id"`        // 用户唯一标识
	UserName string `json:"user_name" gorm:"column:username"`     // 用户名
	Password string `json:"-" gorm:"column:password"`             // 密码
	Avatar   string `json:"avatar" gorm:"column:avatar"`          // 头像
	Email    string `json:"email"  gorm:"column:email"`           // 邮箱
	Token    string `json:"token,omitempty"  gorm:"column:token"` // 其他平台的唯一id
}

// BeforeCreate 在创建用户记录之前生成 UserID
func (u *UserModel) BeforeCreate(tx *gorm.DB) error {
	if u.UserID == "" {
		for {
			u.UserID = generateRandomID(10)
			var count int64
			tx.Model(&UserModel{}).Where("user_id = ?", u.UserID).Count(&count)
			if count == 0 {
				break
			}
		}
	}
	return nil
}

// generateRandomID 生成指定长度的随机字符串
func generateRandomID(length int) string {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}
