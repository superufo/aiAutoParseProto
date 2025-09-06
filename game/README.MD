## 🎯 项目目标

请分析一个在线游戏平台的通信协议，并使用 Golang 开发完整的后端服务，包含 Gin Web 框架、WebSocket 实时通信和 Protobuf3 数据序列化。

## 🔍 协议分析目标网站

**目标URL**: `https://mini.game/gamedetail?gamekey=miniA_crash` 

## 📋 要求 Cursor AI 完成的任务

### 第一阶段：协议分析

1. **分析通信协议**：
   - 识别 HTTP RESTful API 端点
   - 分析 WebSocket 连接参数和消息格式
   - 解析数据序列化格式（JSON/Protobuf等）
   - 识别认证和授权机制
2. **生成协议文档**：
   - 创建 ` docs/protocol.md` 详细说明通信协议
   - 包含请求/响应示例
   - 说明消息流和数据格式
   - 消息帧格式（长度前缀、类型码、proto 名称、负载）
   - 所有消息类型的 proto3 定义（含编号、字段含义、示例十六进制 dump）
     时序图（客户端↔服务端）

### 第二阶段：Golang 后端开发

1. **技术栈要求**：
   - **Web框架**: Gin
   - **实时通信**: Gorilla WebSocket
   - **数据序列化**: Protobuf3
   - **数据库**: Redis + Mysql
   - **认证**: JWT
2. **项目结构**：

go

```
game-backend/
├── proto/                 # .proto 文件定义
├── internal/
│   ├── handler/          # HTTP 处理器
│   ├── service/          # 业务逻辑层
│   ├── model/            # 数据模型
│   ├── websocket/        # WebSocket 处理
│   └── middleware/       # 中间件
├── config/               # 配置文件
├── pkg/                 # 公共包
└── cmd/server/          # 主程序入口
```

<svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 12 12" fill="none" class="_9bc997d _33882ae"><path d="M-5.24537e-07 0C-2.34843e-07 6.62742 5.37258 12 12 12L0 12L-5.24537e-07 0Z" fill="currentColor"></path></svg>

<svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 12 12" fill="none" class="_9bc997d _28d7e84"><path d="M-5.24537e-07 0C-2.34843e-07 6.62742 5.37258 12 12 12L0 12L-5.24537e-07 0Z" fill="currentColor"></path></svg>

1. **必须实现的功能**：
   - WebSocket 连接管理
   - Protobuf 消息编码/解码
   - 玩家状态同步
   - 游戏事件处理
   - 实时排行榜
   - 用户认证和授权

## 📊 预期输出

### 文档输出：

1. **PROTOCOL_ANALYSIS.md** - 完整的协议分析文档
2. **API_DOCUMENTATION.md** - RESTful API 文档
3. **WEBSOCKET_SPEC.md** - WebSocket 消息规范
4. **DEPLOYMENT.md** - 部署指南

### 代码输出：

1. 完整的 Golang 项目源代码
2. Protobuf 定义文件（.proto）
3. Dockerfile 和 docker-compose.yml
4. 配置文件和环境变量示例
5. 单元测试和集成测试

## 🛠️ 开发约束

- 代码必须包含详细注释
- 遵循 Go 语言最佳实践
- 实现错误处理和日志记录
- 支持配置化部署
- 包含性能监控指标

## 📝 特别说明

**致 Cursor AI**：请根据对目标网站的分析，模拟实现一个完整的游戏后端服务。如果无法直接访问目标网站，可以基于常见的游戏通信模式进行合理设计。

## 🔄 开发流程

1. 首先分析目标网站的通信模式
2. 设计 Protobuf 消息结构
3. 实现 WebSocket 通信层
4. 开发 RESTful API
5. 实现业务逻辑和数据处理
6. 编写测试和文档