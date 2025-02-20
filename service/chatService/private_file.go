package chatService

import (
	"fmt"
	"gin-gorilla/common/chatComm"
	"gin-gorilla/global"
	"gin-gorilla/model"
	"github.com/gorilla/websocket"
	"io"
	"os"
	"path/filepath"
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

	// 获取文件的详细信息
	file := req.File
	filePath, err := saveFile(file, cu.UserID)
	if err != nil {
		_ = sendPrivateSystemMsg(conn, "文件保存失败")
		return
	}

	// 获取目标用户的用户名和头像
	targetName := model.SelectUsername(targetid)
	targetAvatar := model.SelectAvatar(targetid)

	// 发送文件消息到目标用户
	err = sendPrivateMessage(receiver, chatComm.PrivateResponse{
		UserId:       cu.UserID,
		Username:     cu.Username,
		Avatar:       cu.Avatar,
		TargetID:     targetid,
		TargetName:   targetName,
		TargetAvatar: targetAvatar,
		MsgType:      req.MsgType,
		Content:      filePath, // 发送文件链接
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
		Content:      filePath, // 发送文件链接回执
		Date:         time.Now(),
		OnlineCount:  len(usersInRoom),
	})
	if err != nil {
		_ = sendBackSystemMsg(conn, "发送回执失败")
	}
}

// saveFile 保存文件并返回文件路径
func saveFile(file *chatComm.File, userID string) (string, error) {
	// 生成唯一的文件名（可以用时间戳或者UUID）
	uniqueFileName := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Name)
	// 定义文件存储路径：uploads/file/{userID}/
	storageDir := filepath.Join(global.Config.UploadPath, "file", userID)
	// 确保目录存在
	if err := os.MkdirAll(storageDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("创建文件夹失败: %v", err)
	}

	// 文件的完整路径
	filePath := filepath.Join(storageDir, uniqueFileName)

	// 保存文件到本地
	outFile, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("创建文件失败: %v", err)
	}
	defer outFile.Close()

	// 将接收到的文件内容写入磁盘
	_, err = io.Copy(outFile, file.Reader)
	if err != nil {
		return "", fmt.Errorf("保存文件失败: %v", err)
	}

	// 返回文件的存储路径
	return filePath, nil
}
