#!/bin/bash

# Lushop 微服务一键部署脚本
# 使用 Docker Compose 部署完整的 Lushop 系统

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo "========================================"
echo "  Lushop 微服务系统一键部署"
echo "========================================"

# 检查 Docker 和 Docker Compose
check_requirements() {
    echo -e "${YELLOW}检查环境依赖...${NC}"
    
    if ! command -v docker &> /dev/null; then
        echo -e "${RED}❌ Docker 未安装${NC}"
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        echo -e "${RED}❌ Docker Compose 未安装${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}✅ 环境检查通过${NC}"
}

# 停止并清理旧容器
cleanup() {
    echo -e "${YELLOW}清理旧容器...${NC}"
    docker-compose down -v 2>/dev/null || docker compose down -v 2>/dev/null || true
    echo -e "${GREEN}✅ 清理完成${NC}"
}

# 构建镜像
build_images() {
    echo -e "${YELLOW}构建 Docker 镜像...${NC}"
    docker-compose build || docker compose build
    echo -e "${GREEN}✅ 镜像构建完成${NC}"
}

# 启动服务
start_services() {
    echo -e "${YELLOW}启动服务...${NC}"
    
    # 先启动基础设施
    echo "启动基础设施（MySQL, Redis, Consul, Jaeger, RocketMQ）..."
    docker-compose up -d mysql redis consul jaeger rocketmq-namesrv rocketmq-broker rocketmq-console || \
    docker compose up -d mysql redis consul jaeger rocketmq-namesrv rocketmq-broker rocketmq-console
    
    # 等待基础设施就绪
    echo "等待基础设施就绪..."
    sleep 30
    
    # 启动微服务
    echo "启动微服务..."
    docker-compose up -d goods-service inventory-service userauth-service userop-service order-service || \
    docker compose up -d goods-service inventory-service userauth-service userop-service order-service
    
    # 等待微服务就绪
    echo "等待微服务就绪..."
    sleep 20
    
    # 启动 API 网关
    echo "启动 API 网关..."
    docker-compose up -d api-gateway || docker compose up -d api-gateway
    
    echo -e "${GREEN}✅ 所有服务启动完成${NC}"
}

# 检查服务状态
check_status() {
    echo ""
    echo "========================================"
    echo "  服务状态检查"
    echo "========================================"
    
    docker-compose ps || docker compose ps
    
    echo ""
    echo "========================================"
    echo "  服务访问地址"
    echo "========================================"
    echo -e "${GREEN}API 网关:${NC}"
    echo "  HTTP: http://localhost:8001"
    echo "  gRPC: localhost:9001"
    echo ""
    echo -e "${GREEN}管理界面:${NC}"
    echo "  Consul: http://localhost:8500"
    echo "  Jaeger: http://localhost:16686"
    echo "  RocketMQ Console: http://localhost:8080"
    echo ""
    echo -e "${GREEN}微服务端口:${NC}"
    echo "  Goods Service: localhost:50052"
    echo "  Order Service: localhost:50053"
    echo "  Inventory Service: localhost:50054"
    echo "  UserOp Service: localhost:50055"
    echo "  UserAuth Service: localhost:50056"
    echo ""
    echo -e "${GREEN}基础设施:${NC}"
    echo "  MySQL: localhost:3306"
    echo "  Redis: localhost:6379"
    echo "  RocketMQ NameServer: localhost:9876"
}

# 查看日志
view_logs() {
    echo ""
    echo "查看服务日志（按 Ctrl+C 退出）..."
    docker-compose logs -f || docker compose logs -f
}

# 主流程
main() {
    check_requirements
    
    case "${1:-start}" in
        start)
            cleanup
            build_images
            start_services
            check_status
            ;;
        stop)
            echo -e "${YELLOW}停止所有服务...${NC}"
            docker-compose down || docker compose down
            echo -e "${GREEN}✅ 服务已停止${NC}"
            ;;
        restart)
            echo -e "${YELLOW}重启所有服务...${NC}"
            docker-compose restart || docker compose restart
            echo -e "${GREEN}✅ 服务已重启${NC}"
            ;;
        logs)
            view_logs
            ;;
        status)
            check_status
            ;;
        clean)
            echo -e "${YELLOW}清理所有容器和数据...${NC}"
            docker-compose down -v || docker compose down -v
            echo -e "${GREEN}✅ 清理完成${NC}"
            ;;
        *)
            echo "用法: $0 {start|stop|restart|logs|status|clean}"
            echo ""
            echo "命令说明:"
            echo "  start   - 构建并启动所有服务"
            echo "  stop    - 停止所有服务"
            echo "  restart - 重启所有服务"
            echo "  logs    - 查看服务日志"
            echo "  status  - 查看服务状态"
            echo "  clean   - 清理所有容器和数据"
            exit 1
            ;;
    esac
}

main "$@"
