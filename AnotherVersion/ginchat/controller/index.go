package controller

import (
	"ginchat/models/res"
	"ginchat/service/user_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(c *gin.Context) {
	userInfo := user_service.GetUserInfo(c)
	if len(userInfo) > 0 {
		c.Redirect(http.StatusFound, "/home")
	}
	res.OkWithData(gin.H{"OnlineUserCount": OnlineUserCount()}, c)
}
