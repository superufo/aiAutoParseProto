# 📚 API文档

## 🔗 基础信息

- **基础URL**: `http://localhost:8080/api/v1`
- **认证方式**: JWT Bearer Token
- **数据格式**: JSON
- **字符编码**: UTF-8

## 🔐 认证接口

### 用户登录
```http
POST /auth/login
```

**请求参数**:
```json
{
  "username": "string",
  "password": "string"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user_id": 12345,
    "username": "player1",
    "balance": 1000.50,
    "user": {
      "id": 12345,
      "username": "player1",
      "email": "player1@example.com",
      "balance": 1000.50,
      "avatar": "",
      "status": 1,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  }
}
```

### 用户注册
```http
POST /auth/register
```

**请求参数**:
```json
{
  "username": "string",
  "password": "string",
  "email": "string"
}
```

**响应示例**:
```json
{
  "code": 201,
  "message": "注册成功",
  "data": {
    "user_id": 12345,
    "username": "player1",
    "email": "player1@example.com"
  }
}
```

### 用户登出
```http
POST /auth/logout
```

**请求头**:
```
Authorization: Bearer <token>
```

**响应示例**:
```json
{
  "code": 200,
  "message": "登出成功"
}
```

### 获取用户信息
```http
GET /auth/profile
```

**请求头**:
```
Authorization: Bearer <token>
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "id": 12345,
    "username": "player1",
    "email": "player1@example.com",
    "balance": 1000.50,
    "avatar": "",
    "status": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### 更新用户信息
```http
PUT /auth/profile
```

**请求头**:
```
Authorization: Bearer <token>
```

**请求参数**:
```json
{
  "email": "string",
  "avatar": "string"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "更新成功"
}
```

### 刷新Token
```http
POST /auth/refresh
```

**请求头**:
```
Authorization: Bearer <token>
```

**响应示例**:
```json
{
  "code": 200,
  "message": "Token刷新成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

## 🎮 游戏接口

### 获取游戏状态
```http
GET /game/status
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "game_id": "crash_001",
    "status": 1,
    "current_multiplier": 2.45,
    "players_count": 156,
    "next_round_in": 15,
    "server_time": 1640995200
  }
}
```

### 下注
```http
POST /game/bet
```

**请求头**:
```
Authorization: Bearer <token>
```

**请求参数**:
```json
{
  "amount": 10.50,
  "auto_cashout": 2.00
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "下注成功",
  "data": {
    "bet_id": "bet_12345_1640995200",
    "amount": 10.50,
    "auto_cashout": 2.00,
    "status": 0
  }
}
```

### 止盈
```http
POST /game/cashout
```

**请求头**:
```
Authorization: Bearer <token>
```

**请求参数**:
```json
{
  "bet_id": "bet_12345_1640995200"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "止盈成功",
  "data": {
    "bet_id": "bet_12345_1640995200",
    "multiplier": 2.45,
    "payout": 25.73,
    "profit": 15.23
  }
}
```

### 获取下注历史
```http
GET /game/bet/history?page=1&page_size=20
```

**请求头**:
```
Authorization: Bearer <token>
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "bets": [
      {
        "id": 1,
        "bet_id": "bet_12345_1640995200",
        "user_id": 12345,
        "game_id": "crash_001",
        "round_id": "round_1640995200",
        "amount": 10.50,
        "auto_cashout": 2.00,
        "multiplier": 2.45,
        "payout": 25.73,
        "status": 1,
        "cashout_time": "2024-01-01T00:00:00Z",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

### 获取游戏历史
```http
GET /game/history?page=1&page_size=50
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "games": [
      {
        "id": 1,
        "round_id": "round_1640995200",
        "game_id": "crash_001",
        "final_multiplier": 2.45,
        "players_count": 156,
        "total_bets": 5000.00,
        "total_payout": 12250.00,
        "winners_count": 89,
        "start_time": "2024-01-01T00:00:00Z",
        "end_time": "2024-01-01T00:00:30Z",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 1000,
    "page": 1,
    "page_size": 50
  }
}
```

### 获取排行榜
```http
GET /game/leaderboard
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": [
    {
      "id": 1,
      "user_id": 12345,
      "username": "player1",
      "total_winnings": 5000.00,
      "biggest_multiplier": 15.67,
      "rank": 1,
      "updated_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": 2,
      "user_id": 12346,
      "username": "player2",
      "total_winnings": 4500.00,
      "biggest_multiplier": 12.34,
      "rank": 2,
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### 获取用户统计
```http
GET /game/stats
```

**请求头**:
```
Authorization: Bearer <token>
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "id": 1,
    "user_id": 12345,
    "total_bets": 100,
    "total_winnings": 5000.00,
    "biggest_multiplier": 15.67,
    "games_played": 100,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

## 🔌 WebSocket接口

### 连接WebSocket
```javascript
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onopen = function(event) {
    console.log('WebSocket连接已建立');
    
    // 发送握手请求
    const handshake = {
        type: 'handshake',
        token: 'your-jwt-token',
        version: '1.0'
    };
    ws.send(JSON.stringify(handshake));
};

ws.onmessage = function(event) {
    const message = JSON.parse(event.data);
    console.log('收到消息:', message);
};

ws.onclose = function(event) {
    console.log('WebSocket连接已关闭');
};

ws.onerror = function(error) {
    console.error('WebSocket错误:', error);
};
```

### 发送下注消息
```javascript
const betMessage = {
    type: 'player_bet',
    amount: 10.50,
    auto_cashout: 2.00
};
ws.send(JSON.stringify(betMessage));
```

### 发送止盈消息
```javascript
const cashoutMessage = {
    type: 'player_cashout',
    bet_id: 'bet_12345_1640995200'
};
ws.send(JSON.stringify(cashoutMessage));
```

## 📊 状态码说明

| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 201 | 创建成功 |
| 400 | 请求参数错误 |
| 401 | 未授权/Token无效 |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 409 | 资源冲突 |
| 429 | 请求过于频繁 |
| 500 | 服务器内部错误 |

## 🔄 错误处理

### 错误响应格式
```json
{
  "code": 400,
  "message": "错误描述信息"
}
```

### 常见错误

#### 认证错误
```json
{
  "code": 401,
  "message": "认证令牌无效"
}
```

#### 参数错误
```json
{
  "code": 400,
  "message": "请求参数错误: amount字段不能为空"
}
```

#### 余额不足
```json
{
  "code": 400,
  "message": "余额不足"
}
```

#### 游戏状态错误
```json
{
  "code": 400,
  "message": "游戏未进行中"
}
```

## 🧪 测试示例

### 使用curl测试

#### 1. 用户注册
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123",
    "email": "test@example.com"
  }'
```

#### 2. 用户登录
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

#### 3. 获取游戏状态
```bash
curl -X GET http://localhost:8080/api/v1/game/status
```

#### 4. 下注（需要Token）
```bash
curl -X POST http://localhost:8080/api/v1/game/bet \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "amount": 10.50,
    "auto_cashout": 2.00
  }'
```

### 使用JavaScript测试

```javascript
// 登录并获取Token
async function login(username, password) {
    const response = await fetch('/api/v1/auth/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ username, password })
    });
    
    const data = await response.json();
    return data.data.token;
}

// 下注
async function placeBet(token, amount, autoCashout) {
    const response = await fetch('/api/v1/game/bet', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({ amount, auto_cashout: autoCashout })
    });
    
    return await response.json();
}

// 使用示例
(async () => {
    try {
        const token = await login('testuser', 'password123');
        console.log('登录成功，Token:', token);
        
        const betResult = await placeBet(token, 10.50, 2.00);
        console.log('下注结果:', betResult);
    } catch (error) {
        console.error('操作失败:', error);
    }
})();
```

## 📝 注意事项

1. **Token有效期**: JWT Token默认有效期为24小时
2. **请求频率限制**: API请求限制为每秒10次，登录接口限制为每秒1次
3. **下注限制**: 最小下注金额1元，最大下注金额1000元
4. **止盈限制**: 最小止盈倍数1.01倍，最大止盈倍数1000倍
5. **WebSocket连接**: 支持断线重连，建议实现心跳机制
6. **错误处理**: 所有接口都返回统一的错误格式
7. **数据验证**: 所有输入参数都会进行严格验证
8. **安全考虑**: 生产环境请使用HTTPS和WSS协议
