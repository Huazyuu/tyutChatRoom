package chatService

import (
	"bytes"
	"fmt"
	"gin-gorilla/model/ctype"

	"os"

	"gin-gorilla/common/chatComm"
	"gin-gorilla/global"
	"gin-gorilla/model"
	"github.com/gorilla/websocket"

	"io"
	"mime/multipart"
	"net/http"

	"time"
)

func handlePrivateFileMessage(conn *websocket.Conn, cu chatComm.ChatUser, targetid string, req chatComm.PrivateRequest) {
	roomID := generateRoomID(cu.UserID, targetid)
	mu.Lock()
	usersInRoom, exists := chatComm.ConnPrivateMap[roomID]
	mu.Unlock()
	if !exists {
		_ = sendBackSystemMsg(conn, "目标用户不在线或房间不存在")
		return
	}
	var receiver *websocket.Conn
	// 获取对方连接conn
	for _, user := range usersInRoom {
		if user.UserID == targetid {
			receiver = user.Conn
			break
		}
	}
	if receiver == nil {
		_ = sendBackSystemMsg(conn, "目标用户不在线")
		return
	}

	if req.MsgType == ctype.FileUploadMsg {
		err := requestFileUpload(conn, cu, req, targetid)
		if err != nil {
			global.Log.Error(err.Error())
			return
		}
	} else if req.MsgType == ctype.FileDownloadMsg {
		err := requestFileDownload(conn, cu, req, targetid)
		if err != nil {
			global.Log.Error(err.Error())
			return
		}
	}

	// 获取目标用户的用户名和头像
	targetName := model.SelectUsername(targetid)
	targetAvatar := model.SelectAvatar(targetid)

	// 发送文件消息到目标用户
	err := sendPrivateMessage(receiver, chatComm.PrivateResponse{
		UserId:       cu.UserID,
		Username:     cu.Username,
		Avatar:       cu.Avatar,
		TargetID:     targetid,
		TargetName:   targetName,
		TargetAvatar: targetAvatar,
		MsgType:      req.MsgType,
		Content:      req.File.Path, // 发送文件链接
		Date:         time.Now(),
		OnlineCount:  len(usersInRoom),
	})
	if err != nil {
		_ = sendBackSystemMsg(conn, "发送文件消息失败")
		return
	}
	// 回执
	err = sendBackMessage(conn, chatComm.PrivateResponse{
		UserId:       cu.UserID,
		Username:     cu.Username,
		Avatar:       cu.Avatar,
		TargetID:     targetid,
		TargetName:   targetName,
		TargetAvatar: targetAvatar,
		MsgType:      req.MsgType,
		Content:      req.File.Path, // 发送文件链接回执
		Date:         time.Now(),
		OnlineCount:  len(usersInRoom),
	})
	if err != nil {
		_ = sendBackSystemMsg(conn, "发送回执失败")
	}
}

// requestFileUpload 接收前端文件并转发到另一个后端接口
func requestFileUpload(conn *websocket.Conn, cu chatComm.ChatUser, req chatComm.PrivateRequest, targetid string) error {
	// 打开文件
	file, err := os.Open(req.File.Path)
	if err != nil {
		global.Log.Error(err.Error())
		sendBackSystemMsg(conn, "打开文件失败")
		return err
	}
	defer file.Close()

	// 构造 HTTP 请求
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加文件
	fileField, err := writer.CreateFormFile("file", req.File.Name)
	if err != nil {
		global.Log.Error(err.Error())
		sendBackSystemMsg(conn, "创建文件字段失败")
		return err
	}
	_, err = io.Copy(fileField, file)
	if err != nil {
		global.Log.Error(err)
		sendBackSystemMsg(conn, "写入数据失败")
		return err
	}
	writer.Close()
	// 构造 HTTP 请求
	url := fmt.Sprintf("http://127.0.0.1:8080/api/files/upload?target_id=%s", targetid)
	reqHTTP, err := http.NewRequest("POST", url, body)
	if err != nil {
		global.Log.Error(err.Error())
		sendBackSystemMsg(conn, "请求失败")
		return err
	}
	reqHTTP.Header.Set("Content-Type", writer.FormDataContentType())
	// 假设添加 JWT 认证头
	var user model.UserModel
	global.DB.Select("token").Model(&model.UserModel{}).Where("user_id = ?", cu.UserID).First(&user)
	reqHTTP.Header.Set("Authorization", fmt.Sprintf("Bearer %v", user.Token))

	// 发送 HTTP 请求
	client := &http.Client{}
	resp, err := client.Do(reqHTTP)
	if err != nil {
		global.Log.Error(err.Error())
		sendBackSystemMsg(conn, "请求失败")
		return err
	}
	defer resp.Body.Close()
	return nil
}

// requestFileDownload 接收前端文件并转发到另一个后端接口
func requestFileDownload(conn *websocket.Conn, cu chatComm.ChatUser, req chatComm.PrivateRequest, targetid string) error {
	// 构造 HTTP 请求
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	url := fmt.Sprintf("http://127.0.0.1:8080/api/files/download?file=%s", req.File.Name)
	reqHTTP, err := http.NewRequest("GET", url, body)
	if err != nil {
		global.Log.Error(err.Error())
		sendBackSystemMsg(conn, "请求失败")
		return err
	}
	reqHTTP.Header.Set("Content-Type", writer.FormDataContentType())
	// 假设添加 JWT 认证头
	var user model.UserModel
	global.DB.Select("token").Model(&model.UserModel{}).Where("user_id = ?", cu.UserID).First(&user)
	reqHTTP.Header.Set("Authorization", fmt.Sprintf("Bearer %v", user.Token))
	// 发送 HTTP 请求
	client := &http.Client{}
	resp, err := client.Do(reqHTTP)
	if err != nil {
		global.Log.Error(err.Error())
		sendBackSystemMsg(conn, "请求失败")
		return err
	}
	fmt.Println(resp)
	defer resp.Body.Close()
	return nil
}
