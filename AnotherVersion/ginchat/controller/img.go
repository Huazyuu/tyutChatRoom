package controller

import (
	"fmt"
	"ginchat/conf"
	"ginchat/models/res"
	"ginchat/service/img_service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
)

func ImgUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		zap.L().Error("get form err", zap.Error(err))
		return
	}
	filepath := conf.GlobalConf.App.UploadFilePath
	if _, err = os.Stat(filepath); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(filepath, os.ModePerm); err != nil {
				c.String(http.StatusInternalServerError, fmt.Sprintf("create directory err: %s", err.Error()))
				zap.L().Error("create directory err", zap.Error(err))
				return
			}
		} else {
			c.String(http.StatusInternalServerError, fmt.Sprintf("check directory err: %s", err.Error()))
			zap.L().Error("check directory err", zap.Error(err))
			return
		}
	}
	filename := filepath + file.Filename
	fmt.Println(filename)
	// 保存到本地
	if err := c.SaveUploadedFile(file, filename); err != nil {
		res.FailWithMessage(fmt.Sprintf("upload file err: %s", err.Error()), c)
		zap.L().Error("upload file err", zap.Error(err))
		return
	}
	krUpload := img_service.ImgCreate().Upload(filename)

	// 删除临时图片
	// os.Remove(filename)

	res.OkWithData(gin.H{
		"code": 0,
		"data": map[string]interface{}{
			"url": krUpload,
		}}, c)
}
