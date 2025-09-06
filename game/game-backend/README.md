# 🎮 Crash Game Backend

一个基于Go语言开发的Crash游戏后端服务，支持WebSocket实时通信、JWT认证和Protobuf数据序列化。

## ✨ 特性

- 🚀 **高性能**: 基于Gin框架，支持高并发
- 🔌 **实时通信**: WebSocket支持实时游戏状态更新
- 🔐 **安全认证**: JWT Token认证机制
- 📦 **数据序列化**: Protobuf3高效数据序列化
- 🗄️ **数据存储**: MySQL + Redis双存储
- 🐳 **容器化**: Docker + Docker Compose一键部署
- 📊 **监控**: 健康检查和性能监控
- 🛡️ **安全**: 速率限制、CORS、输入验证

## 🏗️ 技术栈

- **Web框架**: Gin
- **实时通信**: Gorilla WebSocket
- **数据序列化**: Protobuf3
- **数据库**: MySQL 8.0
- **缓存**: Redis 7
- **认证**: JWT
- **容器化**: Docker + Docker Compose
- **反向代理**: Nginx

## 📁 项目结构

```
game-backend/
├── cmd/server/          # 主程序入口
├── internal/
│   ├── handler/         # HTTP处理器
│   ├── service/         # 业务逻辑层
│   ├── model/           # 数据模型
│   ├── websocket/       # WebSocket处理
│   └── middleware/      # 中间件
├── proto/               # Protobuf定义
├── config/              # 配置文件
├── pkg/database/        # 数据库连接
├── docs/                # 文档
├── test/                # 测试文件
├── nginx/               # Nginx配置
├── scripts/             # 脚本文件
├── docker-compose.yml   # Docker Compose配置
├── Dockerfile          # Docker镜像构建
└── go.mod              # Go模块依赖
```

## 🚀 快速开始

### 1. 环境要求

- Go 1.21+
- Docker 20.10+
- Docker Compose 2.0+

### 2. 克隆项目

```bash
git clone <repository-url>
cd game-backend
```

### 3. 配置环境

```bash
# 复制配置文件
cp config/config.yaml.example config/config.yaml

# 编辑配置文件
vim config/config.yaml
```

### 4. 启动服务

```bash
# 使用Docker Compose启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f game-backend
```

### 5. 验证部署

```bash
# 检查健康状态
curl http://localhost:8080/health

# 检查API接口
curl http://localhost:8080/api/v1/game/status
```

## 📚 API文档

### 认证接口

- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/logout` - 用户登出
- `GET /api/v1/auth/profile` - 获取用户信息
- `PUT /api/v1/auth/profile` - 更新用户信息
- `POST /api/v1/auth/refresh` - 刷新Token

### 游戏接口

- `GET /api/v1/game/status` - 获取游戏状态
- `POST /api/v1/game/bet` - 下注
- `POST /api/v1/game/cashout` - 止盈
- `GET /api/v1/game/bet/history` - 获取下注历史
- `GET /api/v1/game/history` - 获取游戏历史
- `GET /api/v1/game/leaderboard` - 获取排行榜
- `GET /api/v1/game/stats` - 获取用户统计

### WebSocket接口

- `ws://localhost:8080/ws` - WebSocket连接

## 🔌 WebSocket通信

### 连接建立

```javascript
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onopen = function() {
    // 发送握手请求
    ws.send(JSON.stringify({
        type: 'handshake',
        token: 'your-jwt-token',
        version: '1.0'
    }));
};
```

### 消息类型

- `GameStatusUpdate` - 游戏状态更新
- `PlayerBet` - 玩家下注
- `GameStart` - 游戏开始
- `GameEnd` - 游戏结束
- `PlayerCashout` - 玩家止盈
- `LeaderboardUpdate` - 排行榜更新
- `SystemNotification` - 系统通知

## 🧪 测试

### 运行测试客户端

```bash
# 进入测试目录
cd test

# 运行测试客户端
go run test_client.go
```

### 测试命令

```
> bet 10.50 2.00    # 下注10.50，自动止盈2.00倍
> cashout bet_123   # 止盈指定下注
> status            # 获取游戏状态
> quit              # 退出
```

### 使用curl测试

```bash
# 用户登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser1","password":"password"}'

# 获取游戏状态
curl http://localhost:8080/api/v1/game/status

# 下注（需要Token）
curl -X POST http://localhost:8080/api/v1/game/bet \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"amount":10.50,"auto_cashout":2.00}'
```

## 🔧 配置说明

### 服务器配置

```yaml
server:
  host: "0.0.0.0"
  port: 8080
  mode: "debug"  # debug, release, test
```

### 数据库配置

```yaml
database:
  host: "localhost"
  port: 3306
  username: "crash_user"
  password: "crash_password"
  database: "crash_game"
```

### Redis配置

```yaml
redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0
```

### 游戏配置

```yaml
game:
  min_bet_amount: 1.0
  max_bet_amount: 1000.0
  min_multiplier: 1.01
  max_multiplier: 1000.0
  round_duration: 30
  betting_duration: 15
  waiting_duration: 10
```

## 🐳 Docker部署

### 构建镜像

```bash
docker build -t crash-game-backend .
```

### 运行容器

```bash
docker run -d \
  --name crash-game-backend \
  -p 8080:8080 \
  -e DATABASE_HOST=mysql \
  -e DATABASE_USERNAME=crash_user \
  -e DATABASE_PASSWORD=crash_password \
  crash-game-backend
```

### Docker Compose

```bash
# 启动所有服务
docker-compose up -d

# 停止所有服务
docker-compose down

# 查看日志
docker-compose logs -f
```

## 📊 监控和日志

### 健康检查

```bash
curl http://localhost:8080/health
```

### 查看日志

```bash
# Docker日志
docker-compose logs -f game-backend

# 应用日志
tail -f /var/log/crash-game.log
```

### 性能监控

```bash
# 容器资源使用
docker stats crash-game-backend

# 系统资源
htop
```

## 🔒 安全考虑

- JWT Token认证
- 请求频率限制
- CORS跨域保护
- 输入参数验证
- SQL注入防护
- XSS攻击防护

## 🚀 生产环境部署

### 1. 服务器准备

```bash
# 更新系统
sudo apt update && sudo apt upgrade -y

# 安装Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
```

### 2. 配置生产环境

```bash
# 修改生产配置
vim config/config.prod.yaml

# 使用生产配置启动
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

### 3. 设置开机自启

```bash
sudo systemctl enable docker
```

## 🔄 更新和维护

### 应用更新

```bash
# 拉取最新代码
git pull origin main

# 重新构建镜像
docker-compose build game-backend

# 滚动更新
docker-compose up -d --no-deps game-backend
```

### 数据库维护

```bash
# 连接数据库
docker exec -it crash-game-mysql mysql -u root -p

# 清理过期会话
docker exec crash-game-mysql mysql -u root -p -e "DELETE FROM user_sessions WHERE expires_at < NOW();"
```

## 🐛 故障排除

### 常见问题

1. **服务无法启动**
   ```bash
   docker-compose logs game-backend
   ```

2. **数据库连接失败**
   ```bash
   docker-compose logs mysql
   ```

3. **WebSocket连接问题**
   ```bash
   curl -i -N -H "Connection: Upgrade" -H "Upgrade: websocket" http://localhost:8080/ws
   ```

## 📞 技术支持

如果在使用过程中遇到问题，请：

1. 查看日志文件
2. 检查配置文件
3. 验证网络连接
4. 联系技术支持团队

## 📄 许可证

本项目采用MIT许可证，详情请查看[LICENSE](LICENSE)文件。

## 🤝 贡献

欢迎提交Issue和Pull Request来帮助改进项目。

## 📝 更新日志

### v1.0.0 (2024-01-01)
- 初始版本发布
- 支持基本的Crash游戏功能
- WebSocket实时通信
- JWT认证系统
- MySQL + Redis数据存储
- Docker容器化部署

---

**注意**: 生产环境部署前请务必：
- 修改默认密码
- 配置SSL证书
- 设置防火墙规则
- 配置监控告警
- 制定备份策略
