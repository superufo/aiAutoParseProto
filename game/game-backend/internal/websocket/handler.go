package websocket

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocket升级器
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 在生产环境中应该检查Origin
		return true
	},
}

// ServeWS WebSocket处理器
func ServeWS(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 升级HTTP连接为WebSocket连接
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "WebSocket升级失败: " + err.Error(),
			})
			return
		}

		// 创建客户端
		client := &Client{
			conn:        conn,
			send:        make(chan []byte, 256),
			connectedAt: time.Now(),
			lastActive:  time.Now(),
		}

		// 注册客户端
		hub.register <- client

		// 启动读写协程
		go client.writePump(hub)
		go client.readPump(hub)
	}
}
