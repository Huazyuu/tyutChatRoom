package chatComm

import "github.com/gorilla/websocket"

type ChatUser struct {
	Conn     *websocket.Conn
	UserID   string `json:"userid"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}
