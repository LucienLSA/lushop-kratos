#!/bin/bash

# 启动单个微服务的脚本
# 用法: ./start-services.sh [service_name]

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

# 服务目录映射
declare -A SERVICES=(
    ["user"]="service/user"
    ["userauth"]="service/userauth"
    ["goods"]="service/goods"
    ["order"]="service/order"
    ["inventory"]="service/inventory"
    ["userop"]="service/userop"
    ["gateway"]="lushop"
)

print_usage() {
    echo "用法: $0 [service_name]"
    echo ""
    echo "可用服务:"
    for key in "${!SERVICES[@]}"; do
        echo "  - $key -> ${SERVICES[$key]}"
    done
    echo ""
    echo "示例:"
    echo "  $0 user        # 启动 user 服务"
    echo "  $0 all         # 启动所有服务（在后台）"
}

start_service() {
    local service_name=$1
    local service_path=$2
    
    if [ ! -d "$service_path" ]; then
        echo -e "${RED}错误: 目录不存在 $service_path${NC}"
        return 1
    fi
    
    echo -e "${BLUE}启动服务: $service_name${NC}"
    echo -e "${YELLOW}目录: $service_path${NC}"
    
    cd "$service_path"
    
    # 检查是否存在 kratos.yaml 或 main.go
    if [ -f "cmd/${service_name}/main.go" ] || [ -f "main.go" ]; then
        echo -e "${GREEN}使用 kratos run 启动...${NC}"
        kratos run
    else
        echo -e "${YELLOW}未找到 main.go，尝试直接运行已编译的程序...${NC}"
        # 尝试运行编译好的程序
        if [ -f "bin/${service_name}" ]; then
            ./bin/${service_name}
        else
            echo -e "${RED}错误: 未找到可执行文件${NC}"
            return 1
        fi
    fi
}

if [ $# -eq 0 ]; then
    print_usage
    exit 1
fi

SERVICE_NAME=$1

if [ "$SERVICE_NAME" = "all" ]; then
    echo -e "${YELLOW}将在后台启动所有服务...${NC}"
    for key in "${!SERVICES[@]}"; do
        if [ "$key" != "gateway" ]; then
            echo -e "${BLUE}启动 $key 服务...${NC}"
            bash -c "cd ${SERVICES[$key]} && kratos run > /tmp/lushop-${key}.log 2>&1" &
            sleep 2
        fi
    done
    # 最后启动网关
    echo -e "${BLUE}启动 gateway 服务...${NC}"
    bash -c "cd lushop && kratos run > /tmp/lushop-gateway.log 2>&1" &
    
    echo -e "${GREEN}所有服务已在后台启动${NC}"
    echo -e "${YELLOW}日志位置: /tmp/lushop-*.log${NC}"
    
elif [ "${SERVICES[$SERVICE_NAME]+isset}" ]; then
    start_service "$SERVICE_NAME" "${SERVICES[$SERVICE_NAME]}"
else
    echo -e "${RED}未知服务: $SERVICE_NAME${NC}"
    echo ""
    print_usage
    exit 1
fi
