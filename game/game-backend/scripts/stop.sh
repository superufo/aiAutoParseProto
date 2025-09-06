#!/bin/bash

# Crash Game Backend 停止脚本

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

# 停止服务
stop_services() {
    log_info "停止服务..."
    
    if docker-compose ps | grep -q "Up"; then
        docker-compose down
        log_success "服务已停止"
    else
        log_warning "没有运行中的服务"
    fi
}

# 清理资源
cleanup() {
    log_info "清理资源..."
    
    # 清理未使用的镜像
    log_info "清理未使用的Docker镜像..."
    docker image prune -f
    
    # 清理未使用的容器
    log_info "清理未使用的Docker容器..."
    docker container prune -f
    
    # 清理未使用的网络
    log_info "清理未使用的Docker网络..."
    docker network prune -f
    
    log_success "资源清理完成"
}

# 备份数据
backup_data() {
    log_info "备份数据..."
    
    # 创建备份目录
    mkdir -p backups
    
    # 备份数据库
    if docker-compose ps mysql | grep -q "Up"; then
        log_info "备份MySQL数据库..."
        docker-compose exec -T mysql mysqldump -u root -p crash_game > backups/crash_game_$(date +%Y%m%d_%H%M%S).sql
        log_success "数据库备份完成"
    else
        log_warning "MySQL未运行，跳过数据库备份"
    fi
    
    # 备份Redis数据
    if docker-compose ps redis | grep -q "Up"; then
        log_info "备份Redis数据..."
        docker-compose exec -T redis redis-cli --rdb /data/dump.rdb
        docker cp $(docker-compose ps -q redis):/data/dump.rdb backups/redis_$(date +%Y%m%d_%H%M%S).rdb
        log_success "Redis备份完成"
    else
        log_warning "Redis未运行，跳过Redis备份"
    fi
}

# 显示帮助
show_help() {
    echo "Crash Game Backend 停止脚本"
    echo "=========================="
    echo ""
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -h, --help     显示帮助信息"
    echo "  -c, --cleanup  停止服务并清理资源"
    echo "  -b, --backup   停止服务前备份数据"
    echo "  -a, --all      停止服务、备份数据并清理资源"
    echo ""
    echo "示例:"
    echo "  $0              # 仅停止服务"
    echo "  $0 --cleanup    # 停止服务并清理资源"
    echo "  $0 --backup     # 停止服务前备份数据"
    echo "  $0 --all        # 完整停止流程"
}

# 主函数
main() {
    local cleanup_flag=false
    local backup_flag=false
    
    # 解析参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -c|--cleanup)
                cleanup_flag=true
                shift
                ;;
            -b|--backup)
                backup_flag=true
                shift
                ;;
            -a|--all)
                cleanup_flag=true
                backup_flag=true
                shift
                ;;
            *)
                log_error "未知选项: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    echo "🛑 Crash Game Backend 停止脚本"
    echo "==============================="
    
    # 备份数据
    if [ "$backup_flag" = true ]; then
        backup_data
    fi
    
    # 停止服务
    stop_services
    
    # 清理资源
    if [ "$cleanup_flag" = true ]; then
        cleanup
    fi
    
    echo ""
    log_success "🎉 停止流程完成!"
    
    if [ "$backup_flag" = true ]; then
        log_info "备份文件保存在 backups/ 目录"
    fi
    
    if [ "$cleanup_flag" = true ]; then
        log_info "资源清理完成，可以安全删除项目目录"
    fi
}

# 错误处理
trap 'log_error "脚本执行失败，退出码: $?"' ERR

# 执行主函数
main "$@"
