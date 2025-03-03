package ws

import (
	"encoding/json"
	"ginchat/models"
	"ginchat/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jianfengye/collection"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

type GoServe struct {
	ServeInterface
}

func (goServe *GoServe) RunWs(c *gin.Context) {
	run(c)
}
func (goServe *GoServe) GetOnlineUserCount() int {
	return GetOnlineUserCount()
}
func (goServe *GoServe) GetOnlineRoomUserCount(roomId int) int {
	return GetOnlineRoomUserCount(roomId)
}

func run(gin *gin.Context) {
	wsUpgrader.CheckOrigin = func(r *http.Request) bool { return true }
	c, _ := wsUpgrader.Upgrade(gin.Writer, gin.Request, nil)
	defer c.Close()

	done := make(chan struct{})
	go read(c, done)
	go write(done)
	select {}
}

// read 只写chan 持续读取客户端发送的消息,并根据消息的不同类型进行相应的处理,具备错误处理和心跳检测机制
func read(c *websocket.Conn, done chan<- struct{}) {
	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("read发生错误", zap.Any("err", err))
		}
	}()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			offline <- c
			zap.L().Error("ReadMessage error", zap.Any("err", err))
			c.Close()
			close(done)
			return
		}
		// heart
		if string(message) == "ping" {
			appendPing(c)
			chNotify <- 1
			c.WriteMessage(websocket.TextMessage, []byte(`{"status":0,"data":"heartbeat ok"}`))
			<-chNotify
			continue
		}

		json.Unmarshal(message, &clientMsgData)
		// 互斥锁
		clientMsgLock.Lock()
		clientMsg = clientMsgData
		clientMsgLock.Unlock()

		if clientMsg.Data.Uid != "" && clientMsg.State == msgTypeOnline {
			roomId, _ := getRoomId()
			enterRooms <- wsClients{
				Conn:       c,
				RemoteAddr: c.RemoteAddr().String(),
				Uid:        clientMsg.Data.Uid,
				Username:   clientMsg.Data.Username,
				RoomId:     roomId,
				AvatarId:   clientMsg.Data.AvatarId,
			}
		}
		_, serveMsg := formatServeMsgStr(clientMsg.State, c)
		sMsg <- serveMsg
	}
}

func write(done <-chan struct{}) {
	defer func() {
		// 捕获write抛出的panic
		if err := recover(); err != nil {
			zap.L().Error("write发生错误", zap.Any("err", err))
		}
	}()
	for {
		select {
		// 当 done 通道关闭时，退出 write 函数
		case <-done:
			return
		// 传递用户进入房间的信息
		case r := <-enterRooms:
			handleConnClient(r.Conn)
		// 传递服务器要处理的消息
		case cl := <-sMsg:
			serveMsgStr, _ := json.Marshal(cl)
			switch cl.State {
			case msgTypeOnline, msgTypeSend:
				// 通知其他用户
				notify(cl.Conn, string(serveMsgStr))
			case msgTypeGetOnlineUser:
				chNotify <- 1
				cl.Conn.WriteMessage(websocket.TextMessage, serveMsgStr)
				<-chNotify
			case msgTypePrivateChat:
				chNotify <- 1
				toC := findToUserCoonClient()
				if toC != nil {
					// 发送私聊
					toC.(wsClients).Conn.WriteMessage(websocket.TextMessage, serveMsgStr)
				}
				<-chNotify
			}
		// 传递用户离线的连接
		case ofl := <-offline:
			disconnect(ofl)
		}
	}

}

func getRoomId() (string, int) {
	roomId := clientMsg.Data.RoomId

	roomIdInt, _ := strconv.Atoi(roomId)
	return roomId, roomIdInt
}

func formatServeMsgStr(state State, conn *websocket.Conn) ([]byte, msg) {
	roomId, roomIdInt := getRoomId()
	data := msgData{
		Uid:      clientMsg.Data.Uid,
		Username: clientMsg.Data.Username,
		RoomId:   roomId,
		Time:     time.Now().UnixNano() / 1e6,
	}
	if state == msgTypeSend || state == msgTypePrivateChat {
		data.AvatarId = clientMsg.Data.AvatarId
		content := clientMsg.Data.Content
		if utils.MbStrLen(content) > 800 {
			// 大于800 截断
			data.Content = string([]rune(content)[:800])
		} else {
			data.Content = content
		}

		// save msg
		toUidStr := clientMsg.Data.ToUid
		stringUid := data.Uid
		toUid, _ := strconv.Atoi(toUidStr)
		intUid, _ := strconv.Atoi(stringUid)

		if clientMsg.Data.ImageUrl != "" {
			// 存在图片
			models.SaveContent(map[string]interface{}{
				"user_id":    intUid,
				"to_user_id": toUid,
				"content":    data.Content,
				"room_id":    data.RoomId,
				"image_url":  clientMsg.Data.ImageUrl,
			})
		} else {
			models.SaveContent(map[string]interface{}{
				"user_id":    intUid,
				"to_user_id": toUid,
				"content":    data.Content,
				"room_id":    data.RoomId,
			})
		}
	}
	if state == msgTypeGetOnlineUser {
		roomlist := rooms[roomIdInt]
		data.Count = len(roomlist)
		data.List = roomlist
	}
	jsonStrServeMsg := msg{
		State: state,
		Data:  data,
		Conn:  conn,
	}
	serveMsgStr, _ := json.Marshal(jsonStrServeMsg)
	return serveMsgStr, jsonStrServeMsg

}

func appendPing(c *websocket.Conn) {
	objColl := collection.NewObjCollection(pingMap)
	// 删除相同的
	retColl := utils.Safety.Do(func() interface{} {
		return objColl.Reject(func(obj interface{}, index int) bool {
			// Reject 按照某个方法进行过滤，去掉符合的
			if obj.(pingStorage).RemoteAddr == c.RemoteAddr().String() {
				return true
			}
			return false
		})
	}).(collection.ICollection)
	// 追加
	retColl = utils.Safety.Do(func() interface{} {
		return retColl.Append(pingStorage{
			Conn:       c,
			RemoteAddr: c.RemoteAddr().String(),
			Time:       time.Now().Unix(),
		})
	}).(collection.ICollection)

	interfaces, _ := retColl.ToInterfaces()
	pingMap = interfaces
}

func handleConnClient(c *websocket.Conn) {
	roomId, roomIdInt := getRoomId()
	objColl := collection.NewObjCollection(rooms[roomIdInt])
	retColl := utils.Safety.Do(func() interface{} {
		return objColl.Reject(func(item interface{}, key int) bool {
			if item.(wsClients).Uid == clientMsg.Data.Uid {
				chNotify <- 1
				item.(wsClients).Conn.WriteMessage(websocket.TextMessage, []byte(`{"status":-1,"data":[]}`))
				<-chNotify
				return true
			}
			return false
		})
	}).(collection.ICollection)
	retColl = utils.Safety.Do(func() interface{} {
		return retColl.Append(wsClients{
			Conn:       c,
			RemoteAddr: c.RemoteAddr().String(),
			Uid:        clientMsg.Data.Uid,
			Username:   clientMsg.Data.Username,
			RoomId:     roomId,
			AvatarId:   clientMsg.Data.AvatarId,
		})
	}).(collection.ICollection)

	interfaces, _ := retColl.ToInterfaces()
	rooms[roomIdInt] = interfaces
}

func notify(conn *websocket.Conn, msg string) {
	chNotify <- 1
	_, roomIdInt := getRoomId()
	assignRoom := rooms[roomIdInt]
	for _, con := range assignRoom {
		if con.(wsClients).RemoteAddr != conn.RemoteAddr().String() {
			con.(wsClients).Conn.WriteMessage(websocket.TextMessage, []byte(msg))
		}
	}
	<-chNotify
}

func disconnect(conn *websocket.Conn) {
	_, roomIdInt := getRoomId()
	objColl := collection.NewObjCollection(rooms[roomIdInt])
	retColl := utils.Safety.Do(func() interface{} {
		return objColl.Reject(func(item interface{}, key int) bool {
			if item.(wsClients).RemoteAddr == conn.RemoteAddr().String() {
				data := msgData{
					Username: item.(wsClients).Username,
					Uid:      item.(wsClients).Uid,
					Time:     time.Now().UnixNano() / 1e6, // 13位  10位 => now.Unix()
				}
				jsonStrServeMsg := msg{
					State: msgTypeOffline,
					Data:  data,
				}
				serveMsgStr, _ := json.Marshal(jsonStrServeMsg)

				disMsg := string(serveMsgStr)

				item.(wsClients).Conn.Close()
				notify(conn, disMsg)
				return true
			}
			return false
		})
	}).(collection.ICollection)

	interfaces, _ := retColl.ToInterfaces()
	rooms[roomIdInt] = interfaces
}

// 获取私聊的用户连接
func findToUserCoonClient() interface{} {
	_, roomIdInt := getRoomId()

	toUserUid := clientMsg.Data.ToUid
	assignRoom := rooms[roomIdInt]
	for _, c := range assignRoom {
		stringUid := c.(wsClients).Uid
		if stringUid == toUserUid {
			return c
		}
	}

	return nil
}

func GetOnlineUserCount() int {
	num := 0
	for i := 1; i <= roomCount; i++ {
		num = num + GetOnlineRoomUserCount(i)
	}
	return num
}
func GetOnlineRoomUserCount(roomId int) int {
	return len(rooms[roomId])
}
