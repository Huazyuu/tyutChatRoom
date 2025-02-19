package redisService

import (
	"gin-gorilla/global"
	"gin-gorilla/utils"
	"time"
)

const prefix = "logout_"

func Logout(token string, diff time.Duration) error {
	err := global.Redis.Set(prefix+token, "", diff).Err()
	return err
}
func CheckLogout(token string) bool {
	keys := global.Redis.Keys(prefix + "*").Val()
	if utils.InList(prefix+token, keys) {
		return true
	}
	return false
}
