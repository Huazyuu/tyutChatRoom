package user_api

import (
	"gin-gorilla/global"
	"gin-gorilla/res"
	"gin-gorilla/service"
	"gin-gorilla/utils/jwt"
	"github.com/gin-gonic/gin"
)

func (UsersApi) UserLogoutView(c *gin.Context) {
	_claims, _ := c.Get("claims")
	claims := _claims.(*jwt.CustomClaims)

	token := c.Request.Header.Get("token")
	// global.Log.Debug("token", token)

	err := service.ServiceApp.UserService.Logout(claims, token)
	if err != nil {
		global.Log.Error(err)
		res.FailWithMessage("注销失败", c)
		return
	}

	res.OkWithMessage("注销成功", c)

}
