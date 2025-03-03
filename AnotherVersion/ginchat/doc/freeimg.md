https://freeimage.host/page/api

### 函数功能概述
`freeImgUpload` 函数的主要功能是将本地的一个文件（通常是图片）上传到 Freeimage.host 这个图片托管服务平台，并返回该图片的显示 URL。如果上传过程中出现任何错误，函数将记录错误信息并返回空字符串。

### 代码思路及原理详细分析

#### 1. 获取 API 密钥
```go
apiKey := conf.GlobalConf.PictureBed.ApiKey
```
- 从配置文件中获取 Freeimage.host 的 API 密钥。API 密钥用于身份验证，确保只有拥有有效密钥的用户才能使用该服务。

#### 2. 创建请求体缓冲区和 multipart 写入器
```go
bodyBuffer := &bytes.Buffer{}
bodyWriter := multipart.NewWriter(bodyBuffer)
```
- `bytes.Buffer` 是一个字节缓冲区，用于存储即将发送的 HTTP 请求体。
- `multipart.NewWriter` 创建一个 `multipart/form-data` 格式的写入器，这种格式常用于上传文件，它允许在一个请求中包含多个字段和文件。

#### 3. 添加表单字段
```go
if err := bodyWriter.WriteField("key", apiKey); err != nil {
    log.Printf("写入 API 密钥字段失败: %v", err)
    return ""
}
if err := bodyWriter.WriteField("type", "file"); err != nil {
    log.Printf("写入 type 字段失败: %v", err)
    return ""
}
if err := bodyWriter.WriteField("action", "upload"); err != nil {
    log.Printf("写入 action 字段失败: %v", err)
    return ""
}
```
- `bodyWriter.WriteField` 方法用于向请求体中添加普通的表单字段。这里添加了三个字段：
  - `key`：API 密钥，用于身份验证。
  - `type`：指定上传的类型为文件。
  - `action`：指定操作类型为上传。

#### 4. 创建文件字段并写入文件内容
```go
fileWriter, err := bodyWriter.CreateFormFile("source", path.Base(uploadFile))
if err != nil {
    log.Printf("创建文件字段失败: %v", err)
    return ""
}

file, err := os.Open(uploadFile)
if err != nil {
    log.Printf("打开本地文件失败: %v", err)
    return ""
}
defer file.Close()

if _, err := io.Copy(fileWriter, file); err != nil {
    log.Printf("复制文件内容失败: %v", err)
    return ""
}
```
- `bodyWriter.CreateFormFile` 方法创建一个文件字段，字段名为 `source`，文件名使用 `path.Base(uploadFile)` 获取。
- `os.Open` 打开本地文件。
- `io.Copy` 将打开的本地文件内容复制到 `fileWriter` 中，从而将文件内容添加到请求体中。

#### 5. 关闭 multipart 写入器
```go
if err := bodyWriter.Close(); err != nil {
    log.Printf("关闭 multipart 写入器失败: %v", err)
    return ""
}
```
- 调用 `bodyWriter.Close` 方法关闭写入器，确保请求体的格式正确。

#### 6. 获取请求体的 Content-Type
```go
contentType := bodyWriter.FormDataContentType()
```
- `bodyWriter.FormDataContentType` 方法返回 `multipart/form-data` 格式的请求体的 `Content-Type` 头部信息，该信息用于告诉服务器请求体的格式。

#### 7. 构建 HTTP 请求
```go
request := fasthttp.AcquireRequest()
response := fasthttp.AcquireResponse()
defer func() {
    fasthttp.ReleaseResponse(response)
    fasthttp.ReleaseRequest(request)
}()

request.Header.SetContentType(contentType)
request.Header.SetMethod("POST")
request.SetRequestURI("https://freeimage.host/api/1/upload") 
request.SetBody(bodyBuffer.Bytes())
```
- `fasthttp.AcquireRequest` 和 `fasthttp.AcquireResponse` 分别获取一个请求对象和一个响应对象。
- `defer` 语句确保在函数结束时释放请求和响应对象，避免资源泄漏。
- 设置请求头的 `Content-Type` 为之前获取的 `contentType`。
- 设置请求方法为 `POST`，因为上传文件通常使用 `POST` 请求。
- 设置请求的 URI 为 Freeimage.host 的上传 API 地址。
- 将之前构建好的请求体（存储在 `bodyBuffer` 中）设置到请求对象中。

#### 8. 发送 HTTP 请求
```go
if err := fasthttp.Do(request, response); err != nil {
    log.Printf("发送请求失败: %v", err)
    return ""
}
```
- `fasthttp.Do` 方法发送 HTTP 请求并等待响应。如果发送过程中出现错误，记录错误信息并返回空字符串。

#### 9. 解析响应
```go
var res map[string]any
if err := json.Unmarshal(response.Body(), &res); err != nil {
    log.Printf("解析响应 JSON 失败: %v, 响应内容: %s", err, string(response.Body()))
    return ""
}
```
- `json.Unmarshal` 方法将响应体的 JSON 数据解析为一个 `map[string]any` 类型的变量 `res`。如果解析失败，记录错误信息并返回空字符串。

#### 10. 提取图片显示 URL
```go
imageData, ok := res["image"].(map[string]any)
if!ok {
    log.Printf("响应中缺少 image 字段: %v", res)
    return ""
}

displayURL, ok := imageData["display_url"].(string)
if!ok {
    log.Printf("响应中缺少 display_url 字段: %v", imageData)
    return ""
}

return displayURL
```
- 首先检查响应中是否包含 `image` 字段，并将其转换为 `map[string]any` 类型。如果不存在，记录错误信息并返回空字符串。
- 然后检查 `image` 字段中是否包含 `display_url` 字段，并将其转换为字符串类型。如果不存在，记录错误信息并返回空字符串。
- 如果以上检查都通过，返回 `display_url` 作为图片的显示 URL。

综上所述，该函数通过构建 `multipart/form-data` 格式的 HTTP 请求，将本地文件上传到 Freeimage.host 平台，并解析响应获取图片的显示 URL。在整个过程中，对每个可能出错的步骤都进行了错误处理，确保函数的健壮性。