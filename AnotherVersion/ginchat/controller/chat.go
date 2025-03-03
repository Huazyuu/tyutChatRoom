package controller

import (
	"ginchat/models/res"
	"ginchat/service/message_service"
	"ginchat/service/user_service"
	"github.com/gin-gonic/gin"
	"strconv"
)

func PrivateChat(c *gin.Context) {
	roomId := c.Query("room_id")
	toUid := c.Query("uid")

	userInfo := user_service.GetUserInfo(c)

	uid := strconv.Itoa(int(userInfo["uid"].(uint)))

	msgList := message_service.GetLimitPrivateMsg(uid, toUid, 0)
	res.OkWithData(gin.H{
		"user_info": userInfo,
		"msg_list":  msgList,
		"room_id":   roomId,
	}, c)
}
