package router

import (
	"gin-gorilla/api"
	"gin-gorilla/middleware"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
)

var store = cookie.NewStore([]byte("ZYUUFORYUCOOKIESECRET"))

func (router *RouterGroup) UsersRouter() {
	usersApi := api.ApiGroupApp.UserApi
	router.Use(sessions.Sessions("sessionid", store))
	router.POST("users/register", usersApi.UserRegisterView)
	router.POST("users/login", usersApi.UserEmailLoginView)
	// 鉴权
	router.POST("users/logout", middleware.JwtAuth(), usersApi.UserLogoutView)
	// 列表
	router.GET("users", middleware.JwtAuth(), usersApi.UserListView)

}
