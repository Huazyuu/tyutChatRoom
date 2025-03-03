package user_service

import (
	"ginchat/middleware/session"
	"ginchat/models"
	"ginchat/models/req"
	"ginchat/models/res"
	"ginchat/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetUserInfo 获取session中的信息
func GetUserInfo(c *gin.Context) map[string]any {
	return session.GetSessionUserInfo(c)
}

// Login 登录
func Login(c *gin.Context) {
	var request req.User
	if err := c.ShouldBindJSON(&request); err != nil {
		zap.L().Error("请求解析失败", zap.Any("err", err))
		res.FailWithError(err, &request, c)
		return
	}
	md5Pwd := utils.Md5Encrypt(request.Password)

	user := models.FindUserByField("username", request.Username)
	userInfo := user

	if userInfo.ID > 0 {
		// 用户存在 check pwd
		if userInfo.Password != md5Pwd {
			res.FailWithMessage("密码错误", c)
			return
		}
		models.SaveAvatarId(request.AvatarId, user)
	} else {
		// 新用户
		userInfo = models.AddUser(map[string]any{
			"username":  request.Username,
			"password":  md5Pwd,
			"avatar_id": request.AvatarId,
		})
	}
	if userInfo.ID > 0 {
		session.SaveAuthSession(c, userInfo.ID)
		res.OkWithMessage("登陆成功", c)
		return
	} else {
		res.FailWithCode(res.SystemError, c)
		return
	}
}

// Logout 推出登录
func Logout(c *gin.Context) {
	session.ClearAuthSession(c)
	res.OkWithMessage("退出登录", c)
	// c.Redirect(http.StatusFound, "/")
	return
}
