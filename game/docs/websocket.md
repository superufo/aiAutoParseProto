# 🔌 WebSocket通信规范

## 📋 概述

本文档详细说明了Crash游戏的WebSocket通信协议，包括连接建立、消息格式、消息类型和错误处理。

## 🔗 连接信息

- **WebSocket URL**: `ws://localhost:8080/ws` (开发环境)
- **WebSocket URL**: `wss://your-domain.com/ws` (生产环境)
- **协议版本**: `crash-protocol-v1`
- **心跳间隔**: 54秒
- **连接超时**: 60秒

## 🤝 连接建立

### 1. 建立WebSocket连接
```javascript
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onopen = function(event) {
    console.log('WebSocket连接已建立');
};
```

### 2. 发送握手请求
连接建立后，客户端必须发送握手请求：

```javascript
const handshakeRequest = {
    type: 'handshake',
    token: 'your-jwt-token-here',
    version: '1.0'
};

ws.send(JSON.stringify(handshakeRequest));
```

### 3. 接收握手响应
服务器会返回握手响应：

```json
{
    "type": "handshake_response",
    "status": "success",
    "user_id": 12345,
    "server_time": 1640995200,
    "message": ""
}
```

## 📦 消息格式

### 消息帧结构
所有WebSocket消息都使用以下格式：

```
[4字节长度][1字节类型][JSON数据]
```

- **长度**: 4字节大端序整数，表示JSON数据的长度
- **类型**: 1字节消息类型标识
- **数据**: JSON格式的消息内容

### 消息类型定义

| 类型码 | 消息类型 | 方向 | 描述 |
|--------|----------|------|------|
| 0x01 | GameStatusUpdate | 服务端→客户端 | 游戏状态更新 |
| 0x02 | PlayerBet | 客户端→服务端 | 玩家下注 |
| 0x03 | GameStart | 服务端→客户端 | 游戏开始 |
| 0x04 | GameEnd | 服务端→客户端 | 游戏结束 |
| 0x05 | PlayerCashout | 客户端→服务端 | 玩家止盈 |
| 0x06 | LeaderboardUpdate | 服务端→客户端 | 排行榜更新 |
| 0x07 | SystemNotification | 服务端→客户端 | 系统通知 |
| 0x08 | HandshakeRequest | 客户端→服务端 | 握手请求 |
| 0x09 | HandshakeResponse | 服务端→客户端 | 握手响应 |

## 📨 消息类型详解

### 1. 游戏状态更新 (0x01)

**服务端→客户端**

```json
{
    "game_id": "crash_001",
    "state": 1,
    "current_multiplier": 2.45,
    "players_count": 156,
    "next_round_in": 15,
    "server_time": 1640995200
}
```

**字段说明**:
- `game_id`: 游戏ID
- `state`: 游戏状态 (0:等待, 1:进行中, 2:已结束)
- `current_multiplier`: 当前倍数
- `players_count`: 玩家数量
- `next_round_in`: 下轮开始倒计时(秒)
- `server_time`: 服务器时间戳

### 2. 玩家下注 (0x02)

**客户端→服务端**

```json
{
    "amount": 10.50,
    "auto_cashout": 2.00
}
```

**字段说明**:
- `amount`: 下注金额
- `auto_cashout`: 自动止盈倍数(0表示手动止盈)

**服务端响应**:
```json
{
    "bet_id": "bet_12345_1640995200",
    "user_id": 12345,
    "amount": 10.50,
    "auto_cashout": 2.00,
    "timestamp": 1640995200
}
```

### 3. 游戏开始 (0x03)

**服务端→客户端**

```json
{
    "round_id": "round_1640995200",
    "players_count": 156,
    "total_bet_amount": 5000.00,
    "start_time": 1640995200
}
```

**字段说明**:
- `round_id`: 轮次ID
- `players_count`: 玩家数量
- `total_bet_amount`: 总下注金额
- `start_time`: 开始时间戳

### 4. 游戏结束 (0x04)

**服务端→客户端**

```json
{
    "round_id": "round_1640995200",
    "final_multiplier": 2.45,
    "winners_count": 89,
    "total_payout": 12250.00,
    "end_time": 1640995230
}
```

**字段说明**:
- `round_id`: 轮次ID
- `final_multiplier`: 最终倍数
- `winners_count`: 获胜者数量
- `total_payout`: 总赔付金额
- `end_time`: 结束时间戳

### 5. 玩家止盈 (0x05)

**客户端→服务端**

```json
{
    "bet_id": "bet_12345_1640995200"
}
```

**字段说明**:
- `bet_id`: 下注ID

**服务端响应**:
```json
{
    "bet_id": "bet_12345_1640995200",
    "user_id": 12345,
    "multiplier": 2.45,
    "payout": 25.73,
    "timestamp": 1640995200
}
```

### 6. 排行榜更新 (0x06)

**服务端→客户端**

```json
{
    "entries": [
        {
            "user_id": 12345,
            "username": "player1",
            "total_winnings": 5000.00,
            "biggest_multiplier": 15.67,
            "rank": 1
        },
        {
            "user_id": 12346,
            "username": "player2",
            "total_winnings": 4500.00,
            "biggest_multiplier": 12.34,
            "rank": 2
        }
    ],
    "update_time": 1640995200
}
```

### 7. 系统通知 (0x07)

**服务端→客户端**

```json
{
    "type": "info",
    "message": "游戏即将开始",
    "timestamp": 1640995200
}
```

**字段说明**:
- `type`: 通知类型 ("info", "warning", "error")
- `message`: 通知内容
- `timestamp`: 时间戳

## 🔄 消息流示例

### 完整的游戏流程

```javascript
// 1. 建立连接
const ws = new WebSocket('ws://localhost:8080/ws');

// 2. 连接建立后发送握手
ws.onopen = function() {
    ws.send(JSON.stringify({
        type: 'handshake',
        token: 'your-jwt-token',
        version: '1.0'
    }));
};

// 3. 接收握手响应
ws.onmessage = function(event) {
    const message = JSON.parse(event.data);
    
    if (message.type === 'handshake_response') {
        if (message.status === 'success') {
            console.log('握手成功，用户ID:', message.user_id);
        } else {
            console.error('握手失败:', message.message);
        }
    }
};

// 4. 游戏进行中，接收状态更新
ws.onmessage = function(event) {
    const message = JSON.parse(event.data);
    
    switch (message.type) {
        case 'game_status_update':
            console.log('游戏状态:', message.current_multiplier);
            break;
            
        case 'game_start':
            console.log('游戏开始，轮次ID:', message.round_id);
            break;
            
        case 'game_end':
            console.log('游戏结束，最终倍数:', message.final_multiplier);
            break;
            
        case 'leaderboard_update':
            console.log('排行榜更新:', message.entries);
            break;
            
        case 'system_notification':
            console.log('系统通知:', message.message);
            break;
    }
};

// 5. 发送下注
function placeBet(amount, autoCashout) {
    ws.send(JSON.stringify({
        type: 'player_bet',
        amount: amount,
        auto_cashout: autoCashout
    }));
}

// 6. 发送止盈
function cashout(betId) {
    ws.send(JSON.stringify({
        type: 'player_cashout',
        bet_id: betId
    }));
}
```

## ❌ 错误处理

### 握手错误
```json
{
    "type": "handshake_response",
    "status": "error",
    "user_id": 0,
    "server_time": 1640995200,
    "message": "Token无效"
}
```

### 下注错误
```json
{
    "type": "system_notification",
    "type": "error",
    "message": "余额不足",
    "timestamp": 1640995200
}
```

### 止盈错误
```json
{
    "type": "system_notification",
    "type": "error",
    "message": "下注记录不存在",
    "timestamp": 1640995200
}
```

## 🔧 客户端实现建议

### 1. 连接管理
```javascript
class GameWebSocket {
    constructor(url, token) {
        this.url = url;
        this.token = token;
        this.ws = null;
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 5;
        this.reconnectInterval = 1000;
    }
    
    connect() {
        this.ws = new WebSocket(this.url);
        
        this.ws.onopen = () => {
            console.log('WebSocket连接已建立');
            this.reconnectAttempts = 0;
            this.sendHandshake();
        };
        
        this.ws.onmessage = (event) => {
            this.handleMessage(event.data);
        };
        
        this.ws.onclose = () => {
            console.log('WebSocket连接已关闭');
            this.reconnect();
        };
        
        this.ws.onerror = (error) => {
            console.error('WebSocket错误:', error);
        };
    }
    
    reconnect() {
        if (this.reconnectAttempts < this.maxReconnectAttempts) {
            this.reconnectAttempts++;
            console.log(`尝试重连 (${this.reconnectAttempts}/${this.maxReconnectAttempts})`);
            
            setTimeout(() => {
                this.connect();
            }, this.reconnectInterval * this.reconnectAttempts);
        } else {
            console.error('重连失败，已达到最大重试次数');
        }
    }
    
    sendHandshake() {
        this.send({
            type: 'handshake',
            token: this.token,
            version: '1.0'
        });
    }
    
    send(data) {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify(data));
        } else {
            console.error('WebSocket未连接');
        }
    }
    
    handleMessage(data) {
        try {
            const message = JSON.parse(data);
            // 处理消息
            this.onMessage(message);
        } catch (error) {
            console.error('解析消息失败:', error);
        }
    }
    
    onMessage(message) {
        // 子类实现
    }
}
```

### 2. 心跳机制
```javascript
class HeartbeatManager {
    constructor(ws, interval = 30000) {
        this.ws = ws;
        this.interval = interval;
        this.timer = null;
    }
    
    start() {
        this.timer = setInterval(() => {
            if (this.ws.readyState === WebSocket.OPEN) {
                this.ws.ping();
            }
        }, this.interval);
    }
    
    stop() {
        if (this.timer) {
            clearInterval(this.timer);
            this.timer = null;
        }
    }
}
```

### 3. 消息队列
```javascript
class MessageQueue {
    constructor() {
        this.queue = [];
        this.isProcessing = false;
    }
    
    add(message) {
        this.queue.push(message);
        this.process();
    }
    
    async process() {
        if (this.isProcessing || this.queue.length === 0) {
            return;
        }
        
        this.isProcessing = true;
        
        while (this.queue.length > 0) {
            const message = this.queue.shift();
            await this.handleMessage(message);
        }
        
        this.isProcessing = false;
    }
    
    async handleMessage(message) {
        // 处理消息逻辑
        console.log('处理消息:', message);
    }
}
```

## 🧪 测试工具

### 使用wscat测试
```bash
# 安装wscat
npm install -g wscat

# 连接WebSocket
wscat -c ws://localhost:8080/ws

# 发送握手请求
{"type":"handshake","token":"your-token","version":"1.0"}

# 发送下注请求
{"type":"player_bet","amount":10.50,"auto_cashout":2.00}
```

### 使用浏览器测试
```html
<!DOCTYPE html>
<html>
<head>
    <title>WebSocket测试</title>
</head>
<body>
    <div id="messages"></div>
    <input type="text" id="messageInput" placeholder="输入消息">
    <button onclick="sendMessage()">发送</button>
    
    <script>
        const ws = new WebSocket('ws://localhost:8080/ws');
        const messagesDiv = document.getElementById('messages');
        const messageInput = document.getElementById('messageInput');
        
        ws.onopen = function() {
            addMessage('连接已建立');
            // 发送握手请求
            ws.send(JSON.stringify({
                type: 'handshake',
                token: 'your-token',
                version: '1.0'
            }));
        };
        
        ws.onmessage = function(event) {
            addMessage('收到: ' + event.data);
        };
        
        ws.onclose = function() {
            addMessage('连接已关闭');
        };
        
        ws.onerror = function(error) {
            addMessage('错误: ' + error);
        };
        
        function addMessage(message) {
            const div = document.createElement('div');
            div.textContent = message;
            messagesDiv.appendChild(div);
        }
        
        function sendMessage() {
            const message = messageInput.value;
            if (message) {
                ws.send(message);
                messageInput.value = '';
            }
        }
        
        messageInput.addEventListener('keypress', function(e) {
            if (e.key === 'Enter') {
                sendMessage();
            }
        });
    </script>
</body>
</html>
```

## 📝 注意事项

1. **连接超时**: WebSocket连接超时时间为60秒
2. **心跳机制**: 建议实现心跳机制保持连接活跃
3. **断线重连**: 客户端应实现自动重连机制
4. **消息顺序**: 消息可能乱序到达，需要根据时间戳排序
5. **错误处理**: 所有错误都会通过系统通知消息发送
6. **性能优化**: 避免频繁发送消息，建议使用消息队列
7. **安全考虑**: 生产环境请使用WSS协议
8. **日志记录**: 建议记录所有WebSocket消息用于调试
