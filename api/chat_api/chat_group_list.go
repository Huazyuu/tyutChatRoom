package chat_api

import (
	"errors"
	"gin-gorilla/common/chatComm"
	"gin-gorilla/global"
	"gin-gorilla/model"
	"gin-gorilla/res"
	"gin-gorilla/utils/jwt"
	"github.com/gin-gonic/gin"
	"github.com/liu-cn/json-filter/filter"
	"gorm.io/gorm"
)

func (ChatApi) ChatGroupListView(c *gin.Context) {
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
	if errors.Is(err, gorm.ErrRecordNotFound) || (req.Username == "") {
		// 目标名称错误,将返回所有记录
		query := global.DB.
			Select("user_id", "content", "created_at", "ip", "addr").
			Where("chat_models.is_group = 1  AND msg_type = 4")
		result := query.Offset(offset).Limit(req.Limit).Order(req.Sort).Find(&chats)
		if result.Error != nil {
			global.Log.Error(result.Error)
			res.FailWithMessage(result.Error.Error(), c)
			return
		}
		// 过滤不需要的key
		res.OkWithList(filter.Select("grouplist", chats), result.RowsAffected, c)
	} else {
		query := global.DB.
			Select("user_id", "content", "created_at", "ip", "addr").
			Where("chat_models.is_group = 1 AND user_id = ? AND msg_type = 4", searchid)
		result := query.Offset(offset).Limit(req.Limit).Find(&chats)
		if result.Error != nil {
			global.Log.Error(result.Error)
			res.FailWithMessage(result.Error.Error(), c)
			return
		}
		// 过滤不需要的key
		res.OkWithList(filter.Select("grouplist", chats), result.RowsAffected, c)
	}

}
