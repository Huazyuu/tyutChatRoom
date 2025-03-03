package img_service

import (
	"bytes"
	"encoding/json"
	"ginchat/conf"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"os"
	"path"
)

type FreeImageService struct {
	ImgUploadInterface
}

// Upload https://freeimage.host/page/api
func (f *FreeImageService) Upload(filename string) string {
	return freeImgUpload(filename)
}

// freeImgUpload 函数用于将本地文件上传到 Freeimage.host 并返回图片的显示 URL
func freeImgUpload(uploadFile string) string {
	apiKey := conf.GlobalConf.PictureBed.ApiKey
	// 创建一个字节缓冲区用于存储请求体
	bodyBuffer := &bytes.Buffer{}
	// 创建一个 multipart 文件写入器，方便按照 HTTP 规定格式写入内容
	bodyWriter := multipart.NewWriter(bodyBuffer)

	// 添加必要的表单字段
	if err := bodyWriter.WriteField("key", apiKey); err != nil {
		zap.L().Error("写入 API 密钥字段失败:", zap.Error(err))
		return ""
	}
	if err := bodyWriter.WriteField("type", "file"); err != nil {
		zap.L().Error("写入 TYPE 密钥字段失败:", zap.Error(err))
		return ""
	}
	if err := bodyWriter.WriteField("action", "upload"); err != nil {
		zap.L().Error("写入 action 密钥字段失败:", zap.Error(err))
		return ""
	}

	// 从 bodyWriter 生成 fileWriter，并将文件内容写入 fileWriter
	fileWriter, err := bodyWriter.CreateFormFile("source", path.Base(uploadFile))
	if err != nil {
		zap.L().Error("创建文件字段失败:", zap.Error(err))
		return ""
	}

	// 打开本地文件
	file, err := os.Open(uploadFile)
	if err != nil {
		zap.L().Error("打开本地文件失败:", zap.Error(err))
		return ""
	}
	defer file.Close()

	// 将文件内容复制到请求体中
	if _, err := io.Copy(fileWriter, file); err != nil {
		zap.L().Error("复制文件内容失败:", zap.Error(err))
		return ""
	}

	// 关闭 multipart 写入器
	if err := bodyWriter.Close(); err != nil {
		zap.L().Error("关闭 multipart 写入器失败:", zap.Error(err))
		return ""
	}

	// 获取请求体的 Content-Type
	contentType := bodyWriter.FormDataContentType()

	// 构建请求
	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()
	defer func() {
		// 用完需要释放资源
		fasthttp.ReleaseResponse(response)
		fasthttp.ReleaseRequest(request)
	}()

	// 设置请求头和请求方法
	request.Header.SetContentType(contentType)
	request.Header.SetMethod("POST")
	request.SetRequestURI("https://freeimage.host/api/1/upload")
	request.SetBody(bodyBuffer.Bytes())

	// 发送请求
	if err := fasthttp.Do(request, response); err != nil {
		zap.L().Error("发送请求失败", zap.Error(err))
		return ""
	}

	// 解析响应
	var res map[string]any
	if err := json.Unmarshal(response.Body(), &res); err != nil {
		zap.L().Error("解析响应 JSON 失败", zap.Error(err))
		return ""
	}
	// 检查响应中是否包含 image 字段
	imageData, ok := res["image"].(map[string]any)
	if !ok {
		zap.L().Error("响应缺少image字段", zap.Error(err))
		return ""
	}

	// 检查 image 字段中是否包含 display_url 字段
	displayURL, ok := imageData["display_url"].(string)
	if !ok {
		zap.L().Error("响应缺少display_url字段", zap.Error(err))
		return ""
	}

	return displayURL
}
