package websocket

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"game-backend/proto"
)

// Hub WebSocket连接中心
type Hub struct {
	// 注册的客户端连接
	clients map[*Client]bool

	// 从客户端接收消息的通道
	broadcast chan []byte

	// 注册客户端
	register chan *Client

	// 注销客户端
	unregister chan *Client

	// 互斥锁
	mutex sync.RWMutex

	// 游戏状态
	gameState *GameState
}

// GameState 游戏状态
type GameState struct {
	GameID           string  `json:"game_id"`
	Status           int     `json:"status"` // 0:等待 1:进行中 2:已结束
	CurrentMultiplier float64 `json:"current_multiplier"`
	PlayersCount     int32   `json:"players_count"`
	NextRoundIn      int32   `json:"next_round_in"`
	LastUpdate       int64   `json:"last_update"`
	mutex            sync.RWMutex
}


// MessageType 消息类型
type MessageType byte

const (
	GameStatusUpdate MessageType = 0x01
	PlayerBet        MessageType = 0x02
	GameStart        MessageType = 0x03
	GameEnd          MessageType = 0x04
	PlayerCashout    MessageType = 0x05
	LeaderboardUpdate MessageType = 0x06
	SystemNotification MessageType = 0x07
	HandshakeRequest MessageType = 0x08
	HandshakeResponse MessageType = 0x09
)

// NewHub 创建新的WebSocket中心
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		gameState: &GameState{
			GameID:           "crash_001",
			Status:           0,
			CurrentMultiplier: 1.0,
			PlayersCount:     0,
			NextRoundIn:      10,
			LastUpdate:       time.Now().Unix(),
		},
	}
}

// Run 运行WebSocket中心
func (h *Hub) Run() {
	// 启动游戏状态更新定时器
	go h.gameLoop()

	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastMessage(message)
		}
	}
}

// registerClient 注册客户端
func (h *Hub) registerClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.clients[client] = true
	h.gameState.mutex.Lock()
	h.gameState.PlayersCount++
	h.gameState.mutex.Unlock()

	log.Printf("客户端已连接: %s (用户ID: %d), 当前连接数: %d", 
		client.username, client.userID, len(h.clients))

	// 发送当前游戏状态给新连接的客户端
	h.sendGameStatusToClient(client)
}

// unregisterClient 注销客户端
func (h *Hub) unregisterClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)
		h.gameState.mutex.Lock()
		h.gameState.PlayersCount--
		h.gameState.mutex.Unlock()

		log.Printf("客户端已断开: %s (用户ID: %d), 当前连接数: %d", 
			client.username, client.userID, len(h.clients))
	}
}

// broadcastMessage 广播消息
func (h *Hub) broadcastMessage(message []byte) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	for client := range h.clients {
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}

// sendGameStatusToClient 发送游戏状态给指定客户端
func (h *Hub) sendGameStatusToClient(client *Client) {
	h.gameState.mutex.RLock()
	defer h.gameState.mutex.RUnlock()

	statusUpdate := &proto.GameStatusUpdate{
		GameId:            h.gameState.GameID,
		State:             proto.GameState(h.gameState.Status),
		CurrentMultiplier: h.gameState.CurrentMultiplier,
		PlayersCount:      h.gameState.PlayersCount,
		NextRoundIn:       h.gameState.NextRoundIn,
		ServerTime:        h.gameState.LastUpdate,
	}

	message, err := h.encodeMessage(GameStatusUpdate, statusUpdate)
	if err != nil {
		log.Printf("编码游戏状态消息失败: %v", err)
		return
	}

	select {
	case client.send <- message:
	default:
		close(client.send)
		delete(h.clients, client)
	}
}

// encodeMessage 编码消息
func (h *Hub) encodeMessage(msgType MessageType, data interface{}) ([]byte, error) {
	// 序列化数据
	protoData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("序列化数据失败: %v", err)
	}

	// 创建消息帧: [4字节长度][1字节类型][protobuf数据]
	message := make([]byte, 5+len(protoData))
	
	// 设置长度（大端序）
	binary.BigEndian.PutUint32(message[0:4], uint32(len(protoData)+1))
	
	// 设置消息类型
	message[4] = byte(msgType)
	
	// 设置数据
	copy(message[5:], protoData)

	return message, nil
}

// decodeMessage 解码消息
func (h *Hub) decodeMessage(data []byte) (MessageType, []byte, error) {
	if len(data) < 5 {
		return 0, nil, fmt.Errorf("消息长度不足")
	}

	// 读取长度
	length := binary.BigEndian.Uint32(data[0:4])
	
	// 读取消息类型
	msgType := MessageType(data[4])
	
	// 读取数据
	payload := data[5:]

	if len(payload) != int(length-1) {
		return 0, nil, fmt.Errorf("消息长度不匹配")
	}

	return msgType, payload, nil
}

// gameLoop 游戏循环
func (h *Hub) gameLoop() {
	ticker := time.NewTicker(100 * time.Millisecond) // 100ms更新一次
	defer ticker.Stop()

	for range ticker.C {
		h.updateGameState()
	}
}

// updateGameState 更新游戏状态
func (h *Hub) updateGameState() {
	h.gameState.mutex.Lock()
	defer h.gameState.mutex.Unlock()

	now := time.Now().Unix()
	h.gameState.LastUpdate = now

	// 模拟游戏逻辑
	switch h.gameState.Status {
	case 0: // 等待状态
		h.gameState.NextRoundIn--
		if h.gameState.NextRoundIn <= 0 {
			h.gameState.Status = 1 // 开始游戏
			h.gameState.CurrentMultiplier = 1.0
			h.gameState.NextRoundIn = 30 // 30秒游戏时间
			h.broadcastGameStart()
		}
	case 1: // 游戏进行中
		h.gameState.CurrentMultiplier += 0.01 // 每秒增加0.01倍
		h.gameState.NextRoundIn--
		if h.gameState.NextRoundIn <= 0 {
			h.gameState.Status = 2 // 游戏结束
			h.broadcastGameEnd()
		}
	case 2: // 游戏结束
		h.gameState.Status = 0 // 重置为等待状态
		h.gameState.CurrentMultiplier = 1.0
		h.gameState.NextRoundIn = 10 // 10秒等待时间
	}

	// 广播游戏状态更新
	h.broadcastGameStatusUpdate()
}

// broadcastGameStatusUpdate 广播游戏状态更新
func (h *Hub) broadcastGameStatusUpdate() {
	statusUpdate := &proto.GameStatusUpdate{
		GameId:            h.gameState.GameID,
		State:             proto.GameState(h.gameState.Status),
		CurrentMultiplier: h.gameState.CurrentMultiplier,
		PlayersCount:      h.gameState.PlayersCount,
		NextRoundIn:       h.gameState.NextRoundIn,
		ServerTime:        h.gameState.LastUpdate,
	}

	message, err := h.encodeMessage(GameStatusUpdate, statusUpdate)
	if err != nil {
		log.Printf("编码游戏状态更新消息失败: %v", err)
		return
	}

	h.broadcastMessage(message)
}

// broadcastGameStart 广播游戏开始
func (h *Hub) broadcastGameStart() {
	gameStart := &proto.GameStart{
		RoundId:        fmt.Sprintf("round_%d", time.Now().Unix()),
		PlayersCount:   h.gameState.PlayersCount,
		TotalBetAmount: 0, // 这里应该从数据库获取
		StartTime:      h.gameState.LastUpdate,
	}

	message, err := h.encodeMessage(GameStart, gameStart)
	if err != nil {
		log.Printf("编码游戏开始消息失败: %v", err)
		return
	}

	h.broadcastMessage(message)
}

// broadcastGameEnd 广播游戏结束
func (h *Hub) broadcastGameEnd() {
	gameEnd := &proto.GameEnd{
		RoundId:         fmt.Sprintf("round_%d", time.Now().Unix()),
		FinalMultiplier: h.gameState.CurrentMultiplier,
		WinnersCount:    0, // 这里应该从数据库获取
		TotalPayout:     0, // 这里应该从数据库获取
		EndTime:         h.gameState.LastUpdate,
	}

	message, err := h.encodeMessage(GameEnd, gameEnd)
	if err != nil {
		log.Printf("编码游戏结束消息失败: %v", err)
		return
	}

	h.broadcastMessage(message)
}

// GetGameState 获取当前游戏状态
func (h *Hub) GetGameState() *GameState {
	h.gameState.mutex.RLock()
	defer h.gameState.mutex.RUnlock()

	// 返回副本以避免竞态条件
	return &GameState{
		GameID:            h.gameState.GameID,
		Status:            h.gameState.Status,
		CurrentMultiplier: h.gameState.CurrentMultiplier,
		PlayersCount:      h.gameState.PlayersCount,
		NextRoundIn:       h.gameState.NextRoundIn,
		LastUpdate:        h.gameState.LastUpdate,
	}
}

// GetClientsCount 获取客户端连接数
func (h *Hub) GetClientsCount() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return len(h.clients)
}
