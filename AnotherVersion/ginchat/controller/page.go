package controller

import (
	"ginchat/models/res"
	"ginchat/service/message_service"
	"ginchat/service/user_service"
	"ginchat/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

func Pagination(c *gin.Context) {
	roomId := c.Query("room_id")
	toUid := c.Query("uid")
	offset := c.Query("offset")
	offsetInt, e := strconv.Atoi(offset)
	if e != nil || offsetInt <= 0 {
		offsetInt = 0
	}

	rooms := []string{"1", "2", "3", "4", "5", "6"}

	if !utils.InArray(roomId, rooms) {
		res.FailWithMessage("错误房间号", c)
		return
	}

	msgList := []map[string]interface{}{}
	if toUid != "" {
		userInfo := user_service.GetUserInfo(c)

		uid := strconv.Itoa(int(userInfo["uid"].(uint)))

		msgList = message_service.GetLimitPrivateMsg(uid, toUid, offsetInt)
	} else {
		msgList = message_service.GetLimitMsg(roomId, offsetInt)
	}

	res.OkWithList(msgList, int64(len(msgList)), c)
	// c.JSON(http.StatusOK, gin.H{
	// 	"code": 0,
	// 	"data": map[string]interface{}{
	// 		"list": msgList,
	// 	},
	// })
}
