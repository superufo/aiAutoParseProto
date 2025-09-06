# 🚀 部署指南

## 📋 部署概述

本文档详细说明了如何部署Crash游戏后端服务，包括开发环境、测试环境和生产环境的部署方法。

## 🛠️ 环境要求

### 系统要求
- **操作系统**: Linux (Ubuntu 20.04+ / CentOS 8+)
- **内存**: 最少2GB，推荐4GB+
- **存储**: 最少10GB可用空间
- **网络**: 稳定的网络连接

### 软件依赖
- **Docker**: 20.10+
- **Docker Compose**: 2.0+
- **Go**: 1.21+ (仅开发环境需要)

## 📦 快速部署

### 1. 克隆项目
```bash
git clone <repository-url>
cd game-backend
```

### 2. 配置环境变量
```bash
# 复制配置文件
cp config/config.yaml.example config/config.yaml

# 编辑配置文件
vim config/config.yaml
```

### 3. 启动服务
```bash
# 使用Docker Compose启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f game-backend
```

### 4. 验证部署
```bash
# 检查健康状态
curl http://localhost:8080/health

# 检查API接口
curl http://localhost:8080/api/v1/game/status
```

## 🔧 详细配置

### 数据库配置

#### MySQL配置
```yaml
database:
  host: "localhost"
  port: 3306
  username: "crash_user"
  password: "crash_password"
  database: "crash_game"
  charset: "utf8mb4"
  max_idle: 10
  max_open: 100
```

#### Redis配置
```yaml
redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0
  pool_size: 10
```

### 服务器配置
```yaml
server:
  host: "0.0.0.0"
  port: 8080
  mode: "release"  # debug, release, test
  read_timeout: 30
  write_timeout: 30
```

### JWT配置
```yaml
jwt:
  secret: "your-secret-key-change-in-production"
  expire_time: 24  # 小时
  issuer: "crash-game"
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
  update_interval: 100
  max_players_per_game: 1000
```

## 🐳 Docker部署

### 构建镜像
```bash
# 构建应用镜像
docker build -t crash-game-backend .

# 查看镜像
docker images | grep crash-game-backend
```

### 运行容器
```bash
# 运行单个容器
docker run -d \
  --name crash-game-backend \
  -p 8080:8080 \
  -e DATABASE_HOST=mysql \
  -e DATABASE_USERNAME=crash_user \
  -e DATABASE_PASSWORD=crash_password \
  -e DATABASE_DATABASE=crash_game \
  -e REDIS_HOST=redis \
  crash-game-backend
```

### Docker Compose部署
```bash
# 启动所有服务
docker-compose up -d

# 停止所有服务
docker-compose down

# 重启服务
docker-compose restart

# 查看日志
docker-compose logs -f

# 扩展服务
docker-compose up -d --scale game-backend=3
```

## 🌐 Nginx配置

### 反向代理配置
```nginx
upstream game_backend {
    server game-backend:8080;
    keepalive 32;
}

server {
    listen 80;
    server_name your-domain.com;
    
    location /api/ {
        proxy_pass http://game_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
    
    location /ws {
        proxy_pass http://game_backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```

### SSL证书配置
```bash
# 生成自签名证书（仅用于测试）
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes

# 使用Let's Encrypt（生产环境推荐）
certbot --nginx -d your-domain.com
```

## 📊 监控和日志

### 日志配置
```yaml
log:
  level: "info"
  format: "json"
  output: "stdout"
  max_size: 100
  max_backups: 3
  max_age: 7
```

### 健康检查
```bash
# HTTP健康检查
curl http://localhost:8080/health

# Docker健康检查
docker inspect --format='{{.State.Health.Status}}' crash-game-backend
```

### 性能监控
```bash
# 查看容器资源使用
docker stats crash-game-backend

# 查看系统资源
htop
```

## 🔒 安全配置

### 防火墙设置
```bash
# Ubuntu/Debian
ufw allow 22
ufw allow 80
ufw allow 443
ufw enable

# CentOS/RHEL
firewall-cmd --permanent --add-service=ssh
firewall-cmd --permanent --add-service=http
firewall-cmd --permanent --add-service=https
firewall-cmd --reload
```

### 环境变量安全
```bash
# 使用环境变量文件
echo "JWT_SECRET=your-super-secret-key" > .env
echo "DATABASE_PASSWORD=your-db-password" >> .env

# 在Docker Compose中使用
docker-compose --env-file .env up -d
```

## 🚀 生产环境部署

### 1. 服务器准备
```bash
# 更新系统
sudo apt update && sudo apt upgrade -y

# 安装Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# 安装Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

### 2. 配置生产环境
```bash
# 创建生产配置
cp config/config.yaml config/config.prod.yaml

# 修改生产配置
vim config/config.prod.yaml
```

### 3. 部署服务
```bash
# 使用生产配置启动
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# 设置开机自启
sudo systemctl enable docker
```

### 4. 备份策略
```bash
# 数据库备份脚本
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
docker exec crash-game-mysql mysqldump -u root -p crash_game > backup_$DATE.sql
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

### 日志清理
```bash
# 清理Docker日志
docker system prune -f

# 清理应用日志
find /var/log -name "*.log" -mtime +7 -delete
```

## 🐛 故障排除

### 常见问题

#### 1. 服务无法启动
```bash
# 检查日志
docker-compose logs game-backend

# 检查端口占用
netstat -tlnp | grep 8080

# 检查配置文件
docker-compose config
```

#### 2. 数据库连接失败
```bash
# 检查数据库状态
docker-compose ps mysql

# 检查数据库日志
docker-compose logs mysql

# 测试数据库连接
docker exec crash-game-mysql mysql -u crash_user -p crash_game
```

#### 3. Redis连接失败
```bash
# 检查Redis状态
docker-compose ps redis

# 测试Redis连接
docker exec crash-game-redis redis-cli ping
```

#### 4. WebSocket连接问题
```bash
# 检查WebSocket端点
curl -i -N -H "Connection: Upgrade" -H "Upgrade: websocket" -H "Sec-WebSocket-Version: 13" -H "Sec-WebSocket-Key: test" http://localhost:8080/ws

# 检查Nginx配置
nginx -t
```

### 性能优化

#### 1. 数据库优化
```sql
-- 添加索引
CREATE INDEX idx_bets_user_created ON bets(user_id, created_at);
CREATE INDEX idx_bets_status ON bets(status);

-- 分析查询性能
EXPLAIN SELECT * FROM bets WHERE user_id = 1 ORDER BY created_at DESC LIMIT 20;
```

#### 2. Redis优化
```bash
# 配置Redis内存策略
echo "maxmemory-policy allkeys-lru" >> redis.conf

# 监控Redis性能
redis-cli --stat
```

#### 3. 应用优化
```yaml
# 调整连接池大小
database:
  max_idle: 20
  max_open: 200

redis:
  pool_size: 20
```

## 📞 技术支持

如果在部署过程中遇到问题，请：

1. 查看日志文件
2. 检查配置文件
3. 验证网络连接
4. 联系技术支持团队

---

**注意**: 生产环境部署前请务必：
- 修改默认密码
- 配置SSL证书
- 设置防火墙规则
- 配置监控告警
- 制定备份策略
