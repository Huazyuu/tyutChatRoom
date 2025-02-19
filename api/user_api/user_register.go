package user_api

import (
	"fmt"
	"gin-gorilla/global"
	"gin-gorilla/res"
	"gin-gorilla/service/userServer"
	"github.com/gin-gonic/gin"
)

type UserRegisterRequest struct {
	UserName string `json:"username" binding:"required" msg:"请输入用户名"`
	Password string `json:"password" binding:"required" msg:"请输入密码"`
	Email    string `json:"email" binding:"required" msg:"请输入邮箱"`
}

// UserRegisterView 创建用户

func (UsersApi) UserRegisterView(c *gin.Context) {
	var cr UserRegisterRequest
	if err := c.ShouldBindJSON(&cr); err != nil {
		res.FailWithError(err, &cr, c)
		return
	}
	err := userServer.UserService{}.CreateUser(cr.UserName, cr.Password, cr.Email)
	if err != nil {
		global.Log.Error(err)
		res.FailWithMessage("邮箱已被使用", c)
		return
	}
	res.OkWithMessage(fmt.Sprintf("用户创建成功 %s ", cr.UserName), c)
}
