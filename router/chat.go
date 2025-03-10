package router

import (
	"gin-gorilla/api"
	"gin-gorilla/middleware"
)

func (router *RouterGroup) ChatRouter() {
	chatApi := api.ApiGroupApp.ChatApi
	// WS
	router.GET("chat_groups", middleware.JwtParamsAuth(), chatApi.ChatGroupView)
	router.GET("chat_private", middleware.JwtParamsAuth(), chatApi.ChatPrivateView)
	// HTTP
	router.GET("chat_groupList", middleware.JwtAuth(), chatApi.ChatGroupListView)
	router.GET("chat_privateList", middleware.JwtAuth(), chatApi.ChatPrivateListView)
}
