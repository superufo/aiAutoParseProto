package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// 测试客户端
type TestClient struct {
	ws       *websocket.Conn
	token    string
	userID   uint
	username string
}

// 登录响应
type LoginResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Token    string `json:"token"`
		UserID   uint   `json:"user_id"`
		Username string `json:"username"`
		Balance  float64 `json:"balance"`
	} `json:"data"`
}

// WebSocket消息
type WSMessage struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data,omitempty"`
	Token     string      `json:"token,omitempty"`
	Version   string      `json:"version,omitempty"`
	Amount    float64     `json:"amount,omitempty"`
	AutoCashout float64   `json:"auto_cashout,omitempty"`
	BetID     string      `json:"bet_id,omitempty"`
}

func main() {
	client := &TestClient{}
	
	// 登录获取Token
	if err := client.login(); err != nil {
		log.Fatalf("登录失败: %v", err)
	}
	
	// 连接WebSocket
	if err := client.connectWebSocket(); err != nil {
		log.Fatalf("WebSocket连接失败: %v", err)
	}
	defer client.ws.Close()
	
	// 启动消息处理协程
	go client.handleMessages()
	
	// 启动命令行交互
	client.startCLI()
}

// 登录
func (c *TestClient) login() error {
	loginData := map[string]string{
		"username": "testuser1",
		"password": "password",
	}
	
	jsonData, _ := json.Marshal(loginData)
	
	resp, err := http.Post("http://localhost:8080/api/v1/auth/login", 
		"application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	var loginResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return err
	}
	
	if loginResp.Code != 200 {
		return fmt.Errorf("登录失败: %s", loginResp.Message)
	}
	
	c.token = loginResp.Data.Token
	c.userID = loginResp.Data.UserID
	c.username = loginResp.Data.Username
	
	fmt.Printf("登录成功: %s (ID: %d, 余额: %.2f)\n", 
		c.username, c.userID, loginResp.Data.Balance)
	
	return nil
}

// 连接WebSocket
func (c *TestClient) connectWebSocket() error {
	url := "ws://localhost:8080/ws"
	
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	
	c.ws = conn
	
	// 发送握手请求
	handshake := WSMessage{
		Type:    "handshake",
		Token:   c.token,
		Version: "1.0",
	}
	
	if err := c.ws.WriteJSON(handshake); err != nil {
		return err
	}
	
	fmt.Println("WebSocket连接已建立，握手请求已发送")
	return nil
}

// 处理WebSocket消息
func (c *TestClient) handleMessages() {
	for {
		var msg map[string]interface{}
		if err := c.ws.ReadJSON(&msg); err != nil {
			log.Printf("读取消息失败: %v", err)
			return
		}
		
		msgType, ok := msg["type"].(string)
		if !ok {
			continue
		}
		
		switch msgType {
		case "handshake_response":
			status, _ := msg["status"].(string)
			if status == "success" {
				fmt.Println("✅ 握手成功")
			} else {
				fmt.Printf("❌ 握手失败: %v\n", msg["message"])
			}
			
		case "game_status_update":
			multiplier, _ := msg["current_multiplier"].(float64)
			players, _ := msg["players_count"].(float64)
			nextRound, _ := msg["next_round_in"].(float64)
			fmt.Printf("🎮 游戏状态: 倍数=%.2f, 玩家=%d, 下轮=%ds\n", 
				multiplier, int(players), int(nextRound))
			
		case "game_start":
			roundID, _ := msg["round_id"].(string)
			players, _ := msg["players_count"].(float64)
			fmt.Printf("🚀 游戏开始: 轮次=%s, 玩家=%d\n", roundID, int(players))
			
		case "game_end":
			multiplier, _ := msg["final_multiplier"].(float64)
			winners, _ := msg["winners_count"].(float64)
			fmt.Printf("🏁 游戏结束: 最终倍数=%.2f, 获胜者=%d\n", 
				multiplier, int(winners))
			
		case "player_bet":
			betID, _ := msg["bet_id"].(string)
			amount, _ := msg["amount"].(float64)
			fmt.Printf("💰 下注成功: ID=%s, 金额=%.2f\n", betID, amount)
			
		case "player_cashout":
			betID, _ := msg["bet_id"].(string)
			multiplier, _ := msg["multiplier"].(float64)
			payout, _ := msg["payout"].(float64)
			fmt.Printf("💸 止盈成功: ID=%s, 倍数=%.2f, 赔付=%.2f\n", 
				betID, multiplier, payout)
			
		case "leaderboard_update":
			fmt.Println("📊 排行榜已更新")
			
		case "system_notification":
			notifType, _ := msg["type"].(string)
			message, _ := msg["message"].(string)
			fmt.Printf("📢 系统通知 [%s]: %s\n", notifType, message)
			
		default:
			fmt.Printf("📨 收到消息: %s\n", msgType)
		}
	}
}

// 启动命令行交互
func (c *TestClient) startCLI() {
	scanner := bufio.NewScanner(os.Stdin)
	
	fmt.Println("\n🎮 Crash游戏测试客户端")
	fmt.Println("命令:")
	fmt.Println("  bet <金额> [自动止盈倍数] - 下注")
	fmt.Println("  cashout <下注ID> - 止盈")
	fmt.Println("  status - 获取游戏状态")
	fmt.Println("  quit - 退出")
	fmt.Println()
	
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}
		
		switch parts[0] {
		case "bet":
			c.handleBet(parts)
		case "cashout":
			c.handleCashout(parts)
		case "status":
			c.handleStatus()
		case "quit", "exit":
			fmt.Println("👋 再见!")
			return
		default:
			fmt.Println("❓ 未知命令，请输入 help 查看帮助")
		}
	}
}

// 处理下注命令
func (c *TestClient) handleBet(parts []string) {
	if len(parts) < 2 {
		fmt.Println("❌ 用法: bet <金额> [自动止盈倍数]")
		return
	}
	
	var amount float64
	if _, err := fmt.Sscanf(parts[1], "%f", &amount); err != nil {
		fmt.Println("❌ 金额格式错误")
		return
	}
	
	var autoCashout float64
	if len(parts) > 2 {
		if _, err := fmt.Sscanf(parts[2], "%f", &autoCashout); err != nil {
			fmt.Println("❌ 自动止盈倍数格式错误")
			return
		}
	}
	
	betMsg := WSMessage{
		Type:        "player_bet",
		Amount:      amount,
		AutoCashout: autoCashout,
	}
	
	if err := c.ws.WriteJSON(betMsg); err != nil {
		fmt.Printf("❌ 发送下注消息失败: %v\n", err)
		return
	}
	
	fmt.Printf("💰 下注请求已发送: 金额=%.2f", amount)
	if autoCashout > 0 {
		fmt.Printf(", 自动止盈=%.2f", autoCashout)
	}
	fmt.Println()
}

// 处理止盈命令
func (c *TestClient) handleCashout(parts []string) {
	if len(parts) < 2 {
		fmt.Println("❌ 用法: cashout <下注ID>")
		return
	}
	
	betID := parts[1]
	
	cashoutMsg := WSMessage{
		Type:  "player_cashout",
		BetID: betID,
	}
	
	if err := c.ws.WriteJSON(cashoutMsg); err != nil {
		fmt.Printf("❌ 发送止盈消息失败: %v\n", err)
		return
	}
	
	fmt.Printf("💸 止盈请求已发送: ID=%s\n", betID)
}

// 处理状态命令
func (c *TestClient) handleStatus() {
	resp, err := http.Get("http://localhost:8080/api/v1/game/status")
	if err != nil {
		fmt.Printf("❌ 获取游戏状态失败: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	var statusResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&statusResp); err != nil {
		fmt.Printf("❌ 解析状态响应失败: %v\n", err)
		return
	}
	
	if code, ok := statusResp["code"].(float64); ok && code == 200 {
		data, _ := statusResp["data"].(map[string]interface{})
		gameID, _ := data["game_id"].(string)
		status, _ := data["status"].(float64)
		multiplier, _ := data["current_multiplier"].(float64)
		players, _ := data["players_count"].(float64)
		nextRound, _ := data["next_round_in"].(float64)
		
		fmt.Printf("🎮 游戏状态:\n")
		fmt.Printf("  游戏ID: %s\n", gameID)
		fmt.Printf("  状态: %d (0:等待 1:进行中 2:已结束)\n", int(status))
		fmt.Printf("  当前倍数: %.2f\n", multiplier)
		fmt.Printf("  玩家数量: %d\n", int(players))
		fmt.Printf("  下轮开始: %ds\n", int(nextRound))
	} else {
		fmt.Printf("❌ 获取状态失败: %v\n", statusResp["message"])
	}
}
