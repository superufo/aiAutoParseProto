#!/bin/bash

# Crash Game Backend 启动脚本

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查依赖
check_dependencies() {
    log_info "检查依赖..."
    
    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安装，请先安装 Docker"
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose 未安装，请先安装 Docker Compose"
        exit 1
    fi
    
    log_success "依赖检查通过"
}

# 检查配置文件
check_config() {
    log_info "检查配置文件..."
    
    if [ ! -f "config/config.yaml" ]; then
        log_warning "配置文件不存在，创建默认配置..."
        cp config/config.yaml.example config/config.yaml
        log_success "默认配置文件已创建"
    fi
    
    log_success "配置文件检查通过"
}

# 启动服务
start_services() {
    log_info "启动服务..."
    
    # 构建镜像
    log_info "构建Docker镜像..."
    docker-compose build
    
    # 启动服务
    log_info "启动所有服务..."
    docker-compose up -d
    
    log_success "服务启动完成"
}

# 等待服务就绪
wait_for_services() {
    log_info "等待服务就绪..."
    
    # 等待MySQL
    log_info "等待MySQL启动..."
    timeout=60
    while [ $timeout -gt 0 ]; do
        if docker-compose exec -T mysql mysqladmin ping -h localhost --silent; then
            log_success "MySQL已就绪"
            break
        fi
        sleep 1
        timeout=$((timeout-1))
    done
    
    if [ $timeout -eq 0 ]; then
        log_error "MySQL启动超时"
        exit 1
    fi
    
    # 等待Redis
    log_info "等待Redis启动..."
    timeout=30
    while [ $timeout -gt 0 ]; do
        if docker-compose exec -T redis redis-cli ping | grep -q PONG; then
            log_success "Redis已就绪"
            break
        fi
        sleep 1
        timeout=$((timeout-1))
    done
    
    if [ $timeout -eq 0 ]; then
        log_error "Redis启动超时"
        exit 1
    fi
    
    # 等待应用
    log_info "等待应用启动..."
    timeout=60
    while [ $timeout -gt 0 ]; do
        if curl -f http://localhost:8080/health &> /dev/null; then
            log_success "应用已就绪"
            break
        fi
        sleep 1
        timeout=$((timeout-1))
    done
    
    if [ $timeout -eq 0 ]; then
        log_error "应用启动超时"
        exit 1
    fi
}

# 显示服务状态
show_status() {
    log_info "服务状态:"
    docker-compose ps
    
    echo ""
    log_info "服务地址:"
    echo "  应用服务: http://localhost:8080"
    echo "  健康检查: http://localhost:8080/health"
    echo "  API状态: http://localhost:8080/api/v1/game/status"
    echo "  WebSocket: ws://localhost:8080/ws"
    echo "  MySQL: localhost:3306"
    echo "  Redis: localhost:6379"
    
    echo ""
    log_info "测试命令:"
    echo "  健康检查: curl http://localhost:8080/health"
    echo "  游戏状态: curl http://localhost:8080/api/v1/game/status"
    echo "  查看日志: docker-compose logs -f"
    echo "  停止服务: docker-compose down"
}

# 主函数
main() {
    echo "🎮 Crash Game Backend 启动脚本"
    echo "================================"
    
    check_dependencies
    check_config
    start_services
    wait_for_services
    show_status
    
    echo ""
    log_success "🎉 所有服务启动成功!"
    log_info "使用 'docker-compose logs -f' 查看日志"
    log_info "使用 'docker-compose down' 停止服务"
}

# 错误处理
trap 'log_error "脚本执行失败，退出码: $?"' ERR

# 执行主函数
main "$@"
