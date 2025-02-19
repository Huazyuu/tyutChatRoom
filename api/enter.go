package api

import (
	"gin-gorilla/api/chat_api"
	"gin-gorilla/api/user_api"
)

type ApiGroup struct {
	UserApi user_api.UsersApi
	ChatApi chat_api.ChatApi
}

var ApiGroupApp = new(ApiGroup)
