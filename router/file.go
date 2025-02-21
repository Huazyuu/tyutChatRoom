package router

import (
	"gin-gorilla/api"
	"gin-gorilla/middleware"
)

func (router *RouterGroup) FilesRouter() {
	fileApi := api.ApiGroupApp.FilesApi
	router.POST("files/upload", middleware.JwtAuth(), fileApi.FileUploadView)
	router.GET("files/download", middleware.JwtAuth(), fileApi.FileDownloadView)
}
