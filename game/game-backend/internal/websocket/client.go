package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// Client WebSocket客户端处理
type Client struct {
	conn *websocket.Conn
	send chan []byte

	// 用户信息
	userID   uint
	username string

	// 连接信息
	connectedAt time.Time
	lastActive  time.Time
}

// readPump 读取客户端消息
func (c *Client) readPump(hub *Hub) {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()

	// 设置读取超时
	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket错误: %v", err)
			}
			break
		}

		c.lastActive = time.Now()
		c.handleMessage(message, hub)
	}
}

// writePump 向客户端发送消息
func (c *Client) writePump(hub *Hub) {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 批量发送队列中的消息
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage 处理客户端消息
func (c *Client) handleMessage(data []byte, hub *Hub) {
	// 解码消息
	msgType, payload, err := hub.decodeMessage(data)
	if err != nil {
		log.Printf("解码消息失败: %v", err)
		c.sendErrorMessage("消息格式错误", hub)
		return
	}

	switch msgType {
	case HandshakeRequest:
		c.handleHandshake(payload, hub)
	case PlayerBet:
		c.handlePlayerBet(payload, hub)
	case PlayerCashout:
		c.handlePlayerCashout(payload, hub)
	default:
		log.Printf("未知消息类型: %d", msgType)
		c.sendErrorMessage("未知消息类型", hub)
	}
}

// handleHandshake 处理握手请求
func (c *Client) handleHandshake(payload []byte, hub *Hub) {
	var handshakeReq struct {
		Token   string `json:"token"`
		Version string `json:"version"`
	}

	if err := json.Unmarshal(payload, &handshakeReq); err != nil {
		log.Printf("解析握手请求失败: %v", err)
		c.sendHandshakeResponse("error", 0, "握手请求格式错误", hub)
		return
	}

	// 这里应该验证JWT Token并获取用户信息
	// 为了演示，我们假设Token验证成功
	if handshakeReq.Token == "" {
		c.sendHandshakeResponse("error", 0, "Token不能为空", hub)
		return
	}

	// 模拟从Token中解析用户信息
	c.userID = 12345 // 这里应该从JWT中解析
	c.username = "player1" // 这里应该从JWT中解析

	// 发送握手响应
	c.sendHandshakeResponse("success", c.userID, "", hub)
	log.Printf("用户握手成功: %s (ID: %d)", c.username, c.userID)
}

// handlePlayerBet 处理玩家下注
func (c *Client) handlePlayerBet(payload []byte, hub *Hub) {
	var betReq struct {
		Amount      float64 `json:"amount"`
		AutoCashout float64 `json:"auto_cashout"`
	}

	if err := json.Unmarshal(payload, &betReq); err != nil {
		log.Printf("解析下注请求失败: %v", err)
		c.sendErrorMessage("下注请求格式错误", hub)
		return
	}

	// 验证下注金额
	if betReq.Amount <= 0 {
		c.sendErrorMessage("下注金额必须大于0", hub)
		return
	}

	// 这里应该调用业务逻辑处理下注
	betID := fmt.Sprintf("bet_%d_%d", c.userID, time.Now().Unix())
	
	playerBet := &proto.PlayerBet{
		BetId:       betID,
		UserId:      int64(c.userID),
		Amount:      betReq.Amount,
		AutoCashout: betReq.AutoCashout,
		Timestamp:   time.Now().Unix(),
	}

	// 广播下注消息
	message, err := hub.encodeMessage(PlayerBet, playerBet)
	if err != nil {
		log.Printf("编码下注消息失败: %v", err)
		c.sendErrorMessage("下注失败", hub)
		return
	}

	hub.broadcastMessage(message)
	log.Printf("用户 %s 下注: %.2f, 自动止盈: %.2f", c.username, betReq.Amount, betReq.AutoCashout)
}

// handlePlayerCashout 处理玩家止盈
func (c *Client) handlePlayerCashout(payload []byte, hub *Hub) {
	var cashoutReq struct {
		BetID string `json:"bet_id"`
	}

	if err := json.Unmarshal(payload, &cashoutReq); err != nil {
		log.Printf("解析止盈请求失败: %v", err)
		c.sendErrorMessage("止盈请求格式错误", hub)
		return
	}

	if cashoutReq.BetID == "" {
		c.sendErrorMessage("下注ID不能为空", hub)
		return
	}

	// 这里应该调用业务逻辑处理止盈
	// 模拟止盈成功
	multiplier := 2.5 // 这里应该从游戏状态获取当前倍数
	payout := 25.0    // 这里应该计算实际赔付

	playerCashout := &proto.PlayerCashout{
		BetId:     cashoutReq.BetID,
		UserId:    int64(c.userID),
		Multiplier: multiplier,
		Payout:    payout,
		Timestamp: time.Now().Unix(),
	}

	// 广播止盈消息
	message, err := hub.encodeMessage(PlayerCashout, playerCashout)
	if err != nil {
		log.Printf("编码止盈消息失败: %v", err)
		c.sendErrorMessage("止盈失败", hub)
		return
	}

	hub.broadcastMessage(message)
	log.Printf("用户 %s 止盈: 倍数 %.2f, 赔付 %.2f", c.username, multiplier, payout)
}

// sendHandshakeResponse 发送握手响应
func (c *Client) sendHandshakeResponse(status string, userID uint, message string, hub *Hub) {
	response := &proto.HandshakeResponse{
		Status:     status,
		UserId:     int64(userID),
		ServerTime: time.Now().Unix(),
		Message:    message,
	}

	msg, err := hub.encodeMessage(HandshakeResponse, response)
	if err != nil {
		log.Printf("编码握手响应失败: %v", err)
		return
	}

	select {
	case c.send <- msg:
	default:
		close(c.send)
	}
}

// sendErrorMessage 发送错误消息
func (c *Client) sendErrorMessage(message string, hub *Hub) {
	notification := &proto.SystemNotification{
		Type:      "error",
		Message:   message,
		Timestamp: time.Now().Unix(),
	}

	msg, err := hub.encodeMessage(SystemNotification, notification)
	if err != nil {
		log.Printf("编码错误消息失败: %v", err)
		return
	}

	select {
	case c.send <- msg:
	default:
		close(c.send)
	}
}

// sendInfoMessage 发送信息消息
func (c *Client) sendInfoMessage(message string, hub *Hub) {
	notification := &proto.SystemNotification{
		Type:      "info",
		Message:   message,
		Timestamp: time.Now().Unix(),
	}

	msg, err := hub.encodeMessage(SystemNotification, notification)
	if err != nil {
		log.Printf("编码信息消息失败: %v", err)
		return
	}

	select {
	case c.send <- msg:
	default:
		close(c.send)
	}
}
