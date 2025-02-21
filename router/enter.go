package router

import (
	"gin-gorilla/global"
	"gin-gorilla/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type RouterGroup struct {
	*gin.RouterGroup
}

func InitRouter() *gin.Engine {
	gin.SetMode(global.Config.System.Env)

	router := gin.Default()

	// 跨域
	allowedOrigins := []string{"http://localhost:63342"}
	router.Use(middleware.Cors(allowedOrigins))
	router.Use(cors.Default()) // 解决跨域问题 gin官方包 "github.com/gin-contrib/cors"

	// api group
	apiRouterGroup := router.Group("api")
	routerGroupApp := RouterGroup{apiRouterGroup}

	// users api
	routerGroupApp.UsersRouter()
	// chat api
	routerGroupApp.ChatRouter()
	// file api
	routerGroupApp.FilesRouter()

	return router
}
