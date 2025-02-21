package file_api

import (
	"encoding/json"
	"gin-gorilla/global"
	"gin-gorilla/model"
	"gin-gorilla/res"
	"gin-gorilla/utils/jwt"
	"github.com/gin-gonic/gin"
	"os"
)

// FileInfoView 处理获取文件信息的接口
func (FilesApi) FileInfoView(c *gin.Context) {
	_claims, exists := c.Get("claims")
	if !exists {
		global.Log.Error("无法获取用户认证信息")
		res.FailWithMessage("无法获取用户认证信息", c)
		return
	}
	claims := _claims.(*jwt.CustomClaims)

	// 获取要查询的文件名
	fileName := c.Query("file")
	if fileName == "" {
		global.Log.Error("文件名不能为空")
		res.FailWithMessage("文件名不能为空", c)
		return
	}

	// 从数据库中查询文件信息
	var fileModel model.FileModel
	result := global.DB.Where("user_id =? AND file_name =?", claims.UserID, fileName).First(&fileModel)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			global.Log.Error("文件记录不存在")
			res.FailWithMessage("文件记录不存在", c)
		} else {
			global.Log.Error("查询文件记录失败")
			res.FailWithMessage("查询文件记录失败", c)
		}
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(fileModel.Path); os.IsNotExist(err) {
		global.Log.Error("文件不存在")
		res.FailWithMessage("文件不存在", c)
		return
	}

	// 构建响应数据
	response := map[string]interface{}{
		"Content": fileModel.Path,
		"File": map[string]interface{}{
			"user_id":   claims.UserID,
			"Name":      fileModel.FileName,
			"Type":      fileModel.FileType,
			"Size":      fileModel.FileSize,
			"CreatedAt": fileModel.CreatedAt,
			"UpdatedAt": fileModel.UpdatedAt,
		},
	}

	// 返回响应
	jsonData, err := json.Marshal(response)
	if err != nil {
		global.Log.Error("JSON 编码失败")
		res.FailWithMessage("JSON 编码失败", c)
		return
	}
	jsonString := string(jsonData)
	res.OkWithData(jsonString, c)
}
