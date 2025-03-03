package ws

import (
	"github.com/gorilla/websocket"
	"sync"
)

type State int

const (
	msgTypeOnline State = iota + 1
	msgTypeOffline
	msgTypeSend
	msgTypeGetOnlineUser
	msgTypePrivateChat
)

var roomCount = 6

// wsClients 客户端连接详情
type wsClients struct {
	Conn       *websocket.Conn `json:"conn"`
	RemoteAddr string          `json:"remote_addr"`
	Uid        string          `json:"uid"`
	Username   string          `json:"username"`
	RoomId     string          `json:"room_id"`
	AvatarId   string          `json:"avatar_id"`
}

// msg 客户端和服务器之间的消息体
type msg struct {
	State State           `json:"state"`
	Data  msgData         `json:"data"`
	Conn  *websocket.Conn `json:"conn"`
}

// msgData 消息数据
type msgData struct {
	Uid      string `json:"uid"`
	Username string `json:"username"`
	AvatarId string `json:"avatar_id"`
	ToUid    string `json:"to_uid"`
	Content  string `json:"content"`
	ImageUrl string `json:"image_url"`
	RoomId   string `json:"room_id"`
	Count    int    `json:"count"`
	List     []any  `json:"list"`
	Time     int64  `json:"time"`
}

// pingStorage 存储心跳信息
type pingStorage struct {
	Conn       *websocket.Conn `json:"conn"`
	RemoteAddr string          `json:"remote_addr"`
	Time       int64           `json:"time"`
}

// 变量定义初始化
var (
	wsUpgrader = websocket.Upgrader{} // 用于将 HTTP 连接升级为 WebSocket 连接

	clientMsg = msg{} // 存储客户端消息

	mutex = sync.Mutex{} // 用于并发控制

	rooms = make(map[int][]any) // 存储每个房间的客户端连接信息

	enterRooms = make(chan wsClients)       // 用于处理用户进入房间的通道
	sMsg       = make(chan msg)             // 用于处理服务器消息的通道
	offline    = make(chan *websocket.Conn) // 用于处理用户离线的通道
	chNotify   = make(chan int, 1)          // 用于并发控制的通道

	pingMap []interface{} // 存储心跳信息

	clientMsgLock = sync.Mutex{} // 用于保护clientMsgData的互斥锁
	clientMsgData = clientMsg    // 临时存储客户端消息数据
)
