package router

import (
	"gin-gorilla/api"
	"gin-gorilla/middleware"
)

func (router *RouterGroup) ChatRouter() {
	chatApi := api.ApiGroupApp.ChatApi
	router.GET("chat_groups", middleware.JwtAuth(), chatApi.ChatGroupView)
	router.GET("chat_private", middleware.JwtAuth(), chatApi.ChatPrivateView)
	router.GET("chat_groupList", middleware.JwtAuth(), chatApi.ChatGroupListView)
	router.GET("chat_privateList", middleware.JwtAuth(), chatApi.ChatPrivateListView)
}
