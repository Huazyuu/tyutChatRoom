package message_service

import "ginchat/models"

func GetLimitMsg(roomId string, offset int) []map[string]any {
	return models.GetLimitMsg(roomId, offset)
}

func GetLimitPrivateMsg(uid, toUId string, offset int) []map[string]any {
	return models.GetLimitPrivateMsg(uid, toUId, offset)
}
