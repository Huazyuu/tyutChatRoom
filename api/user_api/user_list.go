package user_api

import (
	"gin-gorilla/global"
	"gin-gorilla/model"
	"gin-gorilla/res"
	"gin-gorilla/utils/jwt"
	"github.com/gin-gonic/gin"
)

type UserListRequest struct {
	Page  int    `form:"page"`
	Limit int    `form:"limit"`
	Sort  string `form:"sort"`
}

func (UsersApi) UserListView(c *gin.Context) {
	_claims, _ := c.Get("claims")
	claims := _claims.(*jwt.CustomClaims)
	if claims == nil {
		global.Log.Error("请先登录")
		res.FailWithMessage("请先登录", c)
		return
	}
	// 分页查询
	var req UserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		res.FailWithCode(res.ArgumentError, c)
		return
	}
	var users []model.UserModel
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	offset := (req.Page - 1) * req.Limit
	query := global.DB.Select("user_id", "username", "avatar", "email")
	if req.Sort != "" {
		query = query.Order(req.Sort)
	}
	result := query.Offset(offset).Limit(req.Limit).Find(&users)
	if result.Error != nil {
		global.Log.Error(result.Error)
		res.FailWithMessage(result.Error.Error(), c)
		return
	}
	res.OkWithList(users, result.RowsAffected, c)
}
