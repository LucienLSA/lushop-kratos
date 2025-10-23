#!/bin/bash

# UserAuth Service 启动脚本

set -e

echo "========================================="
echo "  UserAuth Service Startup"
echo "========================================="

# 检查 Redis
echo "检查 Redis 连接..."
if ! redis-cli ping > /dev/null 2>&1; then
    echo "❌ Redis 未运行，请先启动 Redis"
    echo "   sudo systemctl start redis"
    exit 1
fi
echo "✅ Redis 连接正常"

# 检查 Consul
echo "检查 Consul 连接..."
if ! curl -s http://127.0.0.1:8500/v1/status/leader > /dev/null 2>&1; then
    echo "❌ Consul 未运行，请先启动 Consul"
    echo "   consul agent -dev"
    exit 1
fi
echo "✅ Consul 连接正常"

# 检查端口
echo "检查端口 50055..."
if lsof -i:50055 > /dev/null 2>&1; then
    echo "⚠️  端口 50055 已被占用"
    read -p "是否杀死占用进程？(y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        lsof -ti:50055 | xargs kill -9
        echo "✅ 已清理端口"
    else
        exit 1
    fi
fi

# 编译（如果需要）
if [ ! -f "./bin/userauth" ]; then
    echo "编译服务..."
    make build
fi

# 启动服务
echo ""
echo "========================================="
echo "  启动 UserAuth Service"
echo "========================================="
echo "  gRPC 端口: 50055"
echo "  服务名: user-auth-service"
echo "  配置文件: configs/config.yaml"
echo "========================================="
echo ""

./bin/userauth -conf ./configs/config.yaml
