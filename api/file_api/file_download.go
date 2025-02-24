package file_api

import (
	"fmt"
	"gin-gorilla/global"
	"gin-gorilla/model"
	"gin-gorilla/res"
	"gin-gorilla/utils/jwt"
	"github.com/gin-gonic/gin"
	"os"
)

// FileDownloadView 处理文件下载的函数
func (FilesApi) FileDownloadView(c *gin.Context) {
	// 获取 claims
	_claims, exists := c.Get("claims")
	if !exists {
		global.Log.Error("无法获取用户认证信息")
		res.FailWithMessage("无法获取用户认证信息", c)
		return
	}
	claims := _claims.(*jwt.CustomClaims)

	// 获取要下载的文件名
	fileName := c.Query("file")
	if fileName == "" {
		global.Log.Error("文件名不能为空")
		res.FailWithMessage("文件名不能为空", c)
		return
	}

	// 从数据库中查询文件信息
	var fileModel model.FileModel
	// result := global.DB.Where("user_id =? AND file_name =?", claims.UserID, fileName).First(&fileModel)
	result := global.DB.Where("file_name = ? ", fileName).First(&fileModel)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			global.Log.Error("文件记录不存在")
			global.Log.Error(claims.UserID, fileName)
			res.FailWithMessage("文件记录不存在", c)
		} else {
			global.Log.Error("查询文件记录失败")
			res.FailWithMessage("查询文件记录失败", c)
		}
		return
	}
	if _, err := os.Stat(fileModel.Path); os.IsNotExist(err) {
		global.Log.Error("文件不存在")
		res.FailWithMessage("文件不存在", c)
		return
	}

	// header
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Content-Type", "application/octet-stream")

	// 发送文件给客户端
	c.File(fileModel.Path)
}
