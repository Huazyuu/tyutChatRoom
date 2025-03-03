package models

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	AvatarId  string    `json:"avatar_id"`
	CreatedAt time.Time `time_format:"2006-01-02 15:04:05" json:"created_at"`
	UpdatedAt time.Time `time_format:"2006-01-02 15:04:05" json:"updated_at"`
}

func AddUser(value any) (u User) {
	u.Username = value.(map[string]any)["username"].(string)
	u.Password = value.(map[string]any)["password"].(string)
	u.AvatarId = value.(map[string]any)["avatar_id"].(string)
	ChatDB.Create(&u)
	return u
}

func SaveAvatarId(AvatarId string, u User) User {
	u.AvatarId = AvatarId
	ChatDB.Save(&u)
	return u
}

func FindUserByField(field, value string) User {
	var u User

	if field == "id" || field == "username" {
		ChatDB.Where(fmt.Sprintf("%s = ?", field), value).First(&u)
	}
	return u
}
func GetOnlineUserList(uids []float64) []map[string]any {
	var results []map[string]any
	ChatDB.Where("id in ?", uids).Find(&results)
	return results
}
