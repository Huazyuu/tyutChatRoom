package controller

import (
	"ginchat/conf"
	"ginchat/ws"
	"github.com/gin-gonic/gin"
)

// serve 映射关系
var serveMap = map[string]ws.ServeInterface{
	"GoServe": &ws.GoServe{},
}

func Create() ws.ServeInterface {
	return serveMap[conf.GlobalConf.App.ServeType]
}

func Start(gin *gin.Context) {
	Create().RunWs(gin)
}

func OnlineUserCount() int {
	return Create().GetOnlineUserCount()
}

func OnlineRoomUserCount(roomId int) int {
	return Create().GetOnlineRoomUserCount(roomId)
}
