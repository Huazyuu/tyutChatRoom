package router

import (
	"gin-gorilla/global"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type RouterGroup struct {
	*gin.RouterGroup
}

func InitRouter() *gin.Engine {
	gin.SetMode(global.Config.System.Env)

	router := gin.Default()

	router.Use(cors.Default()) // 解决跨域问题 gin官方包 "github.com/gin-contrib/cors"

	// api group
	apiRouterGroup := router.Group("api")
	routerGroupApp := RouterGroup{apiRouterGroup}

	// users api
	routerGroupApp.UsersRouter()
	// chat api
	routerGroupApp.ChatRouter()

	return router
}
