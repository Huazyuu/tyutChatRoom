package model

import (
	"time"
)

type FileModel struct {
	ID        uint      `json:"id" gorm:"column:id"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`

	UserID   string `json:"user_id" gorm:"column:user_id"`
	TargetID string `json:"target_id" gorm:"column:target_id"`

	Path     string `json:"path" gorm:"column:path"`
	FileName string `json:"file_name" gorm:"column:file_name"`
	FileType string `json:"file_type" gorm:"column:file_type"`
	FileSize int64  `json:"file_size" gorm:"column:file_size"`
}
