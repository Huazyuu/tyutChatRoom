package file_api

import (
	"errors"
	"fmt"
	"gin-gorilla/common/fileComm"
	"gin-gorilla/global"
	"gin-gorilla/model"
	"gin-gorilla/res"
	"gin-gorilla/utils/jwt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"time"
)

// FileUploadView 处理文件上传的函数
func (FilesApi) FileUploadView(c *gin.Context) {
	_claims, _ := c.Get("claims")
	claims := _claims.(*jwt.CustomClaims)
	userid := claims.UserID
	targetid := c.Query("target_id")
	// id是否存在
	result := global.DB.First(&model.UserModel{}, "user_id = ?", userid)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			global.Log.Error("目标用户不存在")
			res.FailWithMessage("目标用户不存在", c)
			return
		} else {
			global.Log.Error("查找错误")
			res.FailWithMessage("查找错误", c)
			return
		}
	} else {
		global.Log.Infof("目标%s存在", targetid)
	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		global.Log.Error("上传文件失败")
		res.FailWithMessage("上传文件失败", c)
		return
	}

	// 获取文件名 类型
	fileName := file.Filename
	uniqueFileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileName) // 使用纳秒生成唯一文件名
	fileExt := filepath.Ext(fileName)

	uploadsDir := global.Config.UploadPath
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		err := os.MkdirAll(uploadsDir, 0755)
		if err != nil {
			global.Log.Error("创建uploads路径失败")
			res.FailWithMessage("创建uploads路径失败", c)
			return
		}
	}

	userDir := filepath.Join(uploadsDir, claims.Username)
	if _, err := os.Stat(userDir); os.IsNotExist(err) {
		err := os.MkdirAll(userDir, 0755)
		if err != nil {
			global.Log.Error("创建用户路径失败")
			res.FailWithMessage("创建用户路径失败", c)
			return
		}
	}

	filePath := filepath.Join(userDir, uniqueFileName)
	// 保存文件
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		global.Log.Error("保存文件失败")
		res.FailWithMessage("保存文件失败", c)
		return
	}

	// 读取文件信息
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		global.Log.Error("读取文件信息失败")
		res.FailWithMessage("读取文件信息失败", c)
		return
	}

	// 插入文件信息到数据库
	fileModel := model.FileModel{
		UserID:   userid,
		TargetID: targetid,
		Path:     filePath,
		FileName: uniqueFileName,
		FileSize: fileInfo.Size(),
		FileType: fileExt[1:],
	}

	if err := global.DB.Create(&fileModel).Error; err != nil {
		global.Log.Error("插入文件信息到数据库失败")
		res.FailWithMessage("插入文件信息到数据库失败", c)
		return
	}

	resp := fileComm.FileUploadResponse{
		UserID:   userid,
		TargetID: targetid,
		Path:     filePath,
		FileName: uniqueFileName,
		FileSize: int(fileInfo.Size()),
		FileType: fileExt[1:],
	}
	/*	// 返回响应
		response := map[string]interface{}{
			"Content": filePath,
			"File": map[string]interface{}{
				"Name": uniqueFileName,
				"Type": fileExt[1:],
				"Size": fileInfo.Size(),
			},
		}
		jsonData, _ := json.Marshal(response)
		jsonString := string(jsonData)*/
	res.OkWithData(resp, c)
}
