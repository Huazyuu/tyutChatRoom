package controller

import (
	"ginchat/models/res"
	"ginchat/service/user_service"
	"github.com/gin-gonic/gin"
)

func Home(c *gin.Context) {
	userInfo := user_service.GetUserInfo(c)
	rooms := []map[string]interface{}{
		{"id": 1, "num": OnlineRoomUserCount(1)},
		{"id": 2, "num": OnlineRoomUserCount(2)},
		{"id": 3, "num": OnlineRoomUserCount(3)},
		{"id": 4, "num": OnlineRoomUserCount(4)},
		{"id": 5, "num": OnlineRoomUserCount(5)},
		{"id": 6, "num": OnlineRoomUserCount(6)},
	}
	res.OkWithData(gin.H{
		"rooms":     rooms,
		"user_info": userInfo,
	}, c)
}
