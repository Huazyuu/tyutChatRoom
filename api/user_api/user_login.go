package user_api

import (
	"fmt"
	"gin-gorilla/common/userCommon"
	"gin-gorilla/global"
	"gin-gorilla/model"
	"gin-gorilla/plugins/email"
	"gin-gorilla/res"
	"gin-gorilla/utils/jwt"
	"gin-gorilla/utils/pwd"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"math/rand"
)

var emailFirstReq string

func (UsersApi) UserEmailLoginView(c *gin.Context) {
	var cr userCommon.UserLoginRequest

	if err := c.ShouldBindJSON(&cr); err != nil {
		res.FailWithError(err, &cr, c)
		return
	}
	session := sessions.Default(c)
	// 第一次请求发送code
	if cr.Code == nil {
		emailFirstReq = cr.Email
		// 发送验证码存入session
		code := randCode(4)
		session.Set("email_code", code)
		err := session.Save()
		if err != nil {
			global.Log.Error(err)
			res.FailWithMessage("session错误", c)
			return
		}
		err = email.NewCode().Send(cr.Email, "你的验证码是"+code)
		if err != nil {
			global.Log.Error(err)
			return
		}
		res.OkWithMessage("验证码已发送", c)
		return
	}
	code := session.Get("email_code")
	if code != *cr.Code {
		res.FailWithMessage("验证码错误", c)
		return
	}

	// 登录
	// code
	if cr.Email != emailFirstReq {
		res.FailWithMessage("请求错误,非接受验证码邮箱", c)
		return
	}
	// email
	var userModel model.UserModel
	err := global.DB.Take(&userModel, "email = ? ", cr.Email).Error
	if err != nil {
		// 没找到
		global.Log.Warn("邮箱未注册")
		res.FailWithMessage("邮箱未注册", c)
		return
	}
	// password
	isCheck, _ := pwd.CheckPwd(userModel.Password, cr.Password)
	if !isCheck {
		global.Log.Warn("邮箱或密码错误")
		res.FailWithMessage("邮箱或密码错误", c)
		return
	}
	// token
	token, err := jwt.GenToken(jwt.JwtPayLoad{
		Username: userModel.UserName,
		UserID:   userModel.UserID,
	})
	if err != nil {
		global.Log.Error(err)
		res.FailWithMessage("token生成失败", c)
		return
	}
	global.DB.Model(&model.UserModel{}).Where("user_id = ?", userModel.UserID).Update("token", token)
	res.OkWithData(token, c)
}

func randCode(length int) string {
	return fmt.Sprintf("%04v", rand.Intn(10000))
}
