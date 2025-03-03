package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type APIResponseUrl struct {
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
	imageURL := "https://images.pexels.com/photos/30707838/pexels-photo-30707838.jpeg" // 替换为实际的图片 URL

	// 构建请求参数
	params := url.Values{}
	params.Add("key", apiKey)
	params.Add("action", "upload")
	params.Add("source", imageURL)
	params.Add("format", "json")

	// 创建 POST 请求
	resp, err := http.PostForm("https://freeimage.host/api/1/upload", params)
	if err != nil {
		log.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("读取响应内容失败: %v", err)
	}

	// 解析 JSON 响应
	var response APIResponseUrl
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatalf("解析 JSON 响应失败: %v", err)
	}

	// 输出响应信息
	fmt.Printf("状态码: %d\n", response.StatusCode)
	fmt.Printf("状态文本: %s\n", response.StatusTxt)
	fmt.Printf("图片 URL: %s\n", response.Image.URL)
}
