package wsService

import (
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
			return true
		},
	}
	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
