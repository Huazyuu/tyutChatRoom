package chat_api

import (
	"errors"
	"fmt"
	"gin-gorilla/common/chatComm"
	"gin-gorilla/global"
	"gin-gorilla/model"
	"gin-gorilla/res"
	"gin-gorilla/utils/jwt"
	"github.com/gin-gonic/gin"
	"github.com/liu-cn/json-filter/filter"
	"gorm.io/gorm"
)

func (ChatApi) ChatPrivateListView(c *gin.Context) {
	_claims, _ := c.Get("claims")
	claims := _claims.(*jwt.CustomClaims)
	if claims == nil {
		global.Log.Error("请先登录")
		res.FailWithMessage("请先登录", c)
		return
	}
	var req chatComm.ChatListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		res.FailWithCode(res.ArgumentError, c)
		return
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	offset := (req.Page - 1) * req.Limit

	var chats []model.ChatModel
	// 目标 id
	var userModel model.UserModel
	err := global.DB.Select("user_id").Where("username = ?", req.Username).First(&userModel).Error
	searchid := userModel.UserID
	// 目标对象不存在
	if errors.Is(err, gorm.ErrRecordNotFound) {
		global.Log.Error(err.Error())
		res.FailWithMessage(fmt.Sprintf("与 %s 没有通话记录", req.Username), c)
		return
	}
	query := global.DB.
		Select("user_id", "target_id", "content", "created_at", "ip", "addr").
		Where("chat_models.is_group = 0  AND msg_type = 4 AND user_id = ? AND target_id = ?", claims.UserID, searchid)
	result := query.Offset(offset).Limit(req.Limit).Find(&chats)
	if result.Error != nil || result.RowsAffected == 0 {
		global.Log.Error(result.Error)
		res.FailWithMessage(fmt.Sprintf("与 %s 没有通话记录", req.Username), c)
		return
	}
	res.OkWithList(filter.Select("privatelist", chats), result.RowsAffected, c)
}
