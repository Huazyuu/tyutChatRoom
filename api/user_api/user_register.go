package user_api

import (
	"fmt"
	"gin-gorilla/common/userComm"
	"gin-gorilla/global"
	"gin-gorilla/res"
	"gin-gorilla/service/userService"
	"github.com/gin-gonic/gin"
)

// UserRegisterView 创建用户

func (UsersApi) UserRegisterView(c *gin.Context) {
	var cr userComm.UserRegisterRequest
	if err := c.ShouldBindJSON(&cr); err != nil {
		res.FailWithError(err, &cr, c)
		return
	}
	err := userService.UserService{}.CreateUser(cr.UserName, cr.Password, cr.Email)
	if err != nil {
		global.Log.Error(err)
		res.FailWithMessage("邮箱已被使用", c)
		return
	}
	res.OkWithMessage(fmt.Sprintf("用户创建成功 %s ", cr.UserName), c)
}
