# Makefile for Crash Game Backend

.PHONY: help build run test clean docker-build docker-run docker-stop proto

# 默认目标
help:
	@echo "Crash Game Backend - 可用命令:"
	@echo "  build         - 构建应用"
	@echo "  run           - 运行应用"
	@echo "  test          - 运行测试"
	@echo "  clean         - 清理构建文件"
	@echo "  docker-build  - 构建Docker镜像"
	@echo "  docker-run    - 使用Docker Compose运行"
	@echo "  docker-stop   - 停止Docker服务"
	@echo "  proto         - 生成Protobuf代码"
	@echo "  deps          - 下载依赖"
	@echo "  fmt           - 格式化代码"
	@echo "  lint          - 代码检查"

# 构建应用
build:
	@echo "构建应用..."
	go build -o bin/server ./cmd/server

# 运行应用
run:
	@echo "运行应用..."
	go run ./cmd/server

# 运行测试
test:
	@echo "运行测试..."
	go test ./...

# 清理构建文件
clean:
	@echo "清理构建文件..."
	rm -rf bin/
	go clean

# 下载依赖
deps:
	@echo "下载依赖..."
	go mod download
	go mod tidy

# 格式化代码
fmt:
	@echo "格式化代码..."
	go fmt ./...

# 代码检查
lint:
	@echo "代码检查..."
	golangci-lint run

# 生成Protobuf代码
proto:
	@echo "生成Protobuf代码..."
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/crash.proto

# 构建Docker镜像
docker-build:
	@echo "构建Docker镜像..."
	docker build -t crash-game-backend .

# 使用Docker Compose运行
docker-run:
	@echo "使用Docker Compose运行..."
	docker-compose up -d

# 停止Docker服务
docker-stop:
	@echo "停止Docker服务..."
	docker-compose down

# 查看Docker日志
docker-logs:
	@echo "查看Docker日志..."
	docker-compose logs -f

# 重启服务
restart:
	@echo "重启服务..."
	docker-compose restart game-backend

# 进入容器
shell:
	@echo "进入容器..."
	docker-compose exec game-backend sh

# 数据库迁移
migrate:
	@echo "数据库迁移..."
	docker-compose exec game-backend ./server migrate

# 创建测试数据
seed:
	@echo "创建测试数据..."
	docker-compose exec game-backend ./server seed

# 健康检查
health:
	@echo "健康检查..."
	curl -f http://localhost:8080/health || exit 1

# 完整部署
deploy: docker-build docker-run
	@echo "部署完成!"
	@echo "服务地址: http://localhost:8080"
	@echo "健康检查: http://localhost:8080/health"
	@echo "API文档: http://localhost:8080/api/v1/game/status"

# 开发环境设置
dev-setup: deps proto
	@echo "开发环境设置完成!"

# 生产环境部署
prod-deploy:
	@echo "生产环境部署..."
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# 备份数据库
backup:
	@echo "备份数据库..."
	docker-compose exec mysql mysqldump -u root -p crash_game > backup_$(shell date +%Y%m%d_%H%M%S).sql

# 恢复数据库
restore:
	@echo "恢复数据库..."
	@read -p "请输入备份文件名: " file; \
	docker-compose exec -T mysql mysql -u root -p crash_game < $$file

# 监控
monitor:
	@echo "监控服务状态..."
	watch -n 1 'docker-compose ps && echo "" && curl -s http://localhost:8080/health | jq .'

# 性能测试
benchmark:
	@echo "性能测试..."
	go test -bench=. ./...

# 覆盖率测试
coverage:
	@echo "覆盖率测试..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"
