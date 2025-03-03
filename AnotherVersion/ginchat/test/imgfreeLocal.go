package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

type APIResponseLocal struct {
	StatusCode int `json:"status_code"`
	Success    struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	} `json:"success"`
	Image struct {
		Name             string  `json:"name"`
		Extension        string  `json:"extension"`
		Size             int     `json:"size"`
		Width            int     `json:"width"`
		Height           int     `json:"height"`
		Date             string  `json:"date"`
		DateGMT          string  `json:"date_gmt"`
		StorageID        string  `json:"storage_id"`
		Description      string  `json:"description"`
		NSFW             string  `json:"nsfw"`
		MD5              string  `json:"md5"`
		Storage          string  `json:"storage"`
		OriginalFilename string  `json:"original_filename"`
		OriginalExifdata string  `json:"original_exifdata"`
		Views            string  `json:"views"`
		IDEncoded        string  `json:"id_encoded"`
		Filename         string  `json:"filename"`
		Ratio            float64 `json:"ratio"`
		SizeFormatted    string  `json:"size_formatted"`
		MIME             string  `json:"mime"`
		Bits             int     `json:"bits"`
		Channels         string  `json:"channels"`
		URL              string  `json:"url"`
		URLViewer        string  `json:"url_viewer"`
		Thumb            struct {
			Filename      string  `json:"filename"`
			Name          string  `json:"name"`
			Width         int     `json:"width"`
			Height        int     `json:"height"`
			Ratio         float64 `json:"ratio"`
			Size          int     `json:"size"`
			SizeFormatted string  `json:"size_formatted"`
			MIME          string  `json:"mime"`
			Extension     string  `json:"extension"`
			Bits          int     `json:"bits"`
			Channels      string  `json:"channels"`
			URL           string  `json:"url"`
		} `json:"thumb"`
		Medium struct {
			Filename      string  `json:"filename"`
			Name          string  `json:"name"`
			Width         int     `json:"width"`
			Height        int     `json:"height"`
			Ratio         float64 `json:"ratio"`
			Size          int     `json:"size"`
			SizeFormatted string  `json:"size_formatted"`
			MIME          string  `json:"mime"`
			Extension     string  `json:"extension"`
			Bits          int     `json:"bits"`
			Channels      string  `json:"channels"`
			URL           string  `json:"url"`
		} `json:"medium"`
		ViewsLabel string `json:"views_label"`
		DisplayURL string `json:"display_url"`
		HowLongAgo string `json:"how_long_ago"`
	} `json:"image"`
	StatusTxt string `json:"status_txt"`
}

func main() {
	apiKey := "6d207e02198a847aa98d0a2a901485a5"
	filePath := "D:/goProject/ginchat/uploads/屏幕截图_20241229_212523.png"

	// 创建一个字节缓冲区和一个 multipart 写入器
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加 API 密钥字段
	err := writer.WriteField("key", apiKey)
	if err != nil {
		log.Fatalf("写入 API 密钥字段失败: %v", err)
	}

	// 添加操作类型字段
	err = writer.WriteField("action", "upload")
	if err != nil {
		log.Fatalf("写入操作类型字段失败: %v", err)
	}

	// 添加返回格式字段
	err = writer.WriteField("format", "json")
	if err != nil {
		log.Fatalf("写入返回格式字段失败: %v", err)
	}

	// 打开本地文件
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("打开本地文件失败: %v", err)
	}
	defer file.Close()

	// 创建一个文件字段
	part, err := writer.CreateFormFile("source", filePath)
	if err != nil {
		log.Fatalf("创建文件字段失败: %v", err)
	}

	// 将文件内容复制到请求体中
	_, err = io.Copy(part, file)
	if err != nil {
		log.Fatalf("复制文件内容失败: %v", err)
	}

	// 关闭 multipart 写入器
	err = writer.Close()
	if err != nil {
		log.Fatalf("关闭 multipart 写入器失败: %v", err)
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", "https://freeimage.host/api/1/upload", body)
	if err != nil {
		log.Fatalf("创建 HTTP 请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("发送 HTTP 请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("读取响应内容失败: %v", err)
	}

	// 解析 JSON 响应
	var response APIResponseLocal
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		log.Fatalf("解析 JSON 响应失败: %v", err)
	}

	// 输出响应信息
	fmt.Printf("状态码: %d\n", response.StatusCode)
	fmt.Printf("状态文本: %s\n", response.StatusTxt)
	fmt.Printf("图片 URL: %s\n", response.Image.URL)
}
