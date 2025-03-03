package controller

import (
	"ginchat/models/res"
	"ginchat/service/message_service"
	"ginchat/service/user_service"
	"ginchat/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func Room(c *gin.Context) {
	roomId := c.Param("room_id")
	rooms := []string{"1", "2", "3", "4", "5", "6"}

	if !utils.InArray(roomId, rooms) {
		// 默认转跳第一间
		zap.L().Info("默认跳转第一间")
		c.Redirect(http.StatusFound, "/room/1")
	}
	userInfo := user_service.GetUserInfo(c)
	msgList := message_service.GetLimitMsg(roomId, 0)
	res.OkWithData(gin.H{
		"user_info":      userInfo,
		"msg_list":       msgList,
		"msg_list_count": len(msgList),
		"room_id":        roomId,
	}, c)
}
