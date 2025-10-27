#!/bin/bash

# 微服务管理脚本

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

PROJECT_DIR="/home/zzx/GoProject/lushop-kratos-main"

usage() {
    echo "用法: $0 {start|stop|restart|status|logs}"
    echo ""
    echo "命令说明:"
    echo "  start      - 启动所有微服务"
    echo "  stop       - 停止所有微服务"
    echo "  restart    - 重启所有微服务"
    echo "  status     - 查看服务状态"
    echo "  logs <服务名> - 查看指定服务的日志"
    echo ""
    echo "可用服务名: user, userauth, goods, order, inventory, userop, gateway"
}

case "$1" in
    start)
        echo -e "${GREEN}正在启动所有微服务...${NC}"
        cd "$PROJECT_DIR" || exit 1
        if [ -f "./start-services.sh" ]; then
            ./start-services.sh all
        else
            echo -e "${RED}错误: 找不到 start-services.sh${NC}"
            exit 1
        fi
        ;;
    stop)
        echo -e "${YELLOW}正在停止所有微服务...${NC}"
        pkill -f "kratos run"
        pkill -f "go run.*lushop"
        sleep 2
        echo -e "${GREEN}所有微服务已停止${NC}"
        ;;
    restart)
        echo -e "${YELLOW}正在重启所有微服务...${NC}"
        pkill -f "kratos run"
        pkill -f "go run.*lushop"
        sleep 3
        cd "$PROJECT_DIR" || exit 1
        ./start-services.sh all
        ;;
    status)
        echo -e "${BLUE}=== 微服务运行状态 ===${NC}"
        echo ""
        echo -e "${GREEN}运行中的进程:${NC}"
        ps aux | grep -E "kratos run|go run.*lushop" | grep -v grep || echo "无"
        echo ""
        echo -e "${GREEN}Consul 注册的服务:${NC}"
        curl -s http://127.0.0.1:8500/v1/agent/services 2>/dev/null | python3 -m json.tool | grep '"Service":' | sort -u || echo "无法连接到 Consul"
        ;;
    logs)
        if [ -z "$2" ]; then
            echo -e "${RED}错误: 请指定服务名${NC}"
            usage
            exit 1
        fi
        LOG_FILE="/tmp/lushop-$2.log"
        if [ -f "$LOG_FILE" ]; then
            echo -e "${GREEN}查看 $2 服务日志 (Ctrl+C 退出)${NC}"
            tail -f "$LOG_FILE"
        else
            echo -e "${RED}错误: 找不到日志文件 $LOG_FILE${NC}"
            exit 1
        fi
        ;;
    *)
        usage
        exit 1
        ;;
esac
