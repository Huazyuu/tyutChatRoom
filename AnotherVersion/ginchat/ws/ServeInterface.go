package ws

import "github.com/gin-gonic/gin"

type ServeInterface interface {
	RunWs(c *gin.Context)
	GetOnlineUserCount() int
	GetOnlineRoomUserCount(roomId int) int
}
