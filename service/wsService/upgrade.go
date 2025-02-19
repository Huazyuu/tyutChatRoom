package wsService

import (
	"gin-gorilla/res"
	"gin-gorilla/utils/jwt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

// WSUpgarde http 升级 ws 协议
func WSUpgarde(claims *jwt.CustomClaims, c *gin.Context) (*websocket.Conn, error) {
	var upGrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// 鉴权
			if claims == nil {
				return false
			}
			return true
		},
	}
	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		res.FailWithCode(res.ArgumentError, c)
		return nil, err
	}
	return conn, nil
}
