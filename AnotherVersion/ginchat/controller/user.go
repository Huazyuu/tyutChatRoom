package controller

import (
	"ginchat/service/user_service"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	user_service.Login(c)
}
func Logout(c *gin.Context) {
	user_service.Logout(c)
}
