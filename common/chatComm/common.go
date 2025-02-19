package chatComm

import "github.com/gorilla/websocket"

type ChatUser struct {
	Conn     *websocket.Conn
	UserID   string `json:"userid"`
	Username string `json:"username"` // 前端自己生成
	Avatar   string `json:"avatar"`   // 头像 查表
}
