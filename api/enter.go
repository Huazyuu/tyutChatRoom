package api

import (
	"gin-gorilla/api/chat_api"
	"gin-gorilla/api/file_api"
	"gin-gorilla/api/user_api"
)

type ApiGroup struct {
	UserApi  user_api.UsersApi
	ChatApi  chat_api.ChatApi
	FilesApi file_api.FilesApi
}

var ApiGroupApp = new(ApiGroup)
