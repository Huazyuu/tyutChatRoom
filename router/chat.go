package router

import (
	"gin-gorilla/api"
	"gin-gorilla/middleware"
)

func (router *RouterGroup) ChatRouter() {
	chatApi := api.ApiGroupApp.ChatApi
	router.GET("chat_groups", middleware.JwtAuth(), chatApi.ChatGroupView)
	router.GET("chat_private", middleware.JwtAuth(), chatApi.ChatPrivateView)
}
