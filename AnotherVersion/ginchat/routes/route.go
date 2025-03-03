package routes

import (
	"ginchat/controller"
	"ginchat/log"
	"ginchat/middleware/session"
	"github.com/gin-gonic/gin"
)

func InitRoute() *gin.Engine {
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.Use(log.GinLogger(), log.GinRecovery(true))
	routes := router.Group("/", session.EnableCookieSession())
	{
		routes.GET("/", controller.Index)
		routes.POST("/login", controller.Login)
		routes.GET("/logout", controller.Logout)
		// routes.GET("/ws", controller.Start)
		authorized := routes.Group("/", session.AuthSessionMiddle())
		{
			authorized.GET("/ws", controller.Start)
			authorized.GET("/home", controller.Home)
			authorized.GET("/room/:room_id", controller.Room)
			authorized.GET("/private-chat", controller.PrivateChat)
			authorized.POST("/img-upload", controller.ImgUpload)
			authorized.GET("/pagination", controller.Pagination)
		}

	}

	return router
}
