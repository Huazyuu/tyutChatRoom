package models

import (
	"gorm.io/gorm"
	"sort"
	"strconv"
	"time"
)

type Message struct {
	gorm.Model
	ID        uint      `json:"id"`
	UserId    int       `json:"user_id"`
	ToUserId  int       `json:"to_user_id"`
	RoomId    int       `json:"room_id"`
	Content   string    `json:"content"`
	ImageUrl  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SaveContent 保存消息信息
func SaveContent(value any) Message {
	var message Message
	message.UserId = value.(map[string]any)["user_id"].(int)
	message.ToUserId = value.(map[string]any)["to_user_id"].(int)
	message.Content = value.(map[string]any)["content"].(string)
	roomIdInt, _ := strconv.Atoi(value.(map[string]any)["room_id"].(string))
	message.RoomId = roomIdInt
	if _, ok := value.(map[string]any)["image_url"]; ok {
		message.ImageUrl = value.(map[string]any)["image_url"].(string)
	}
	ChatDB.Create(&message)
	return message
}

// GetLimitMsg 获取分页信息
func GetLimitMsg(roomId string, offset int) (results []map[string]any) {
	ChatDB.Model(&Message{}).
		Select("messages.*, users.username ,users.avatar_id").
		Joins("inner join users on users.id = messages.user_id").
		Where("messages.room_id = ? ", roomId).
		Where("messages.to_user_id = 0").
		Order("messages.id desc").
		Offset(offset).
		Limit(100).
		Scan(&results)
	if offset == 0 {
		sort.Slice(results, func(i, j int) bool {
			return results[i]["id"].(uint32) < results[j]["id"].(uint32)
		})
	}
	return results
}

// GetLimitPrivateMsg 获取私聊分页信息
func GetLimitPrivateMsg(uid, toUId string, offset int) (results []map[string]any) {
	ChatDB.Model(&Message{}).
		Select("messages.*, users.username ,users.avatar_id").
		Joins("inner join users on users.id = messages.user_id").
		Where("(messages.user_id = ? and messages.to_user_id= ? ) "+
			"or ( messages.user_id = ? and messages.to_user_id = ?)",
			uid, toUId, toUId, uid).
		Order("messages.id desc").
		Offset(offset).
		Limit(100).
		Scan(&results)
	if offset == 0 {
		sort.Slice(results, func(i, j int) bool {
			return results[i]["id"].(uint32) < results[j]["id"].(uint32)
		})
	}
	return results
}
