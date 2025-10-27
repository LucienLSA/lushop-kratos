#!/bin/bash

# 批量重新生成所有服务的 wire 代码

echo "开始重新生成 wire 代码..."

services=("user" "userauth" "goods" "inventory" "order" "userop")

for service in "${services[@]}"; do
    echo "正在处理 $service 服务..."
    cd "/home/zzx/GoProject/lushop-kratos-main/service/$service/cmd/$service"
    if [ -f "wire.go" ]; then
        wire
        if [ $? -eq 0 ]; then
            echo "✅ $service 服务 wire 生成成功"
        else
            echo "❌ $service 服务 wire 生成失败"
        fi
    else
        echo "⚠️  $service 服务没有 wire.go 文件"
    fi
done

echo "wire 代码生成完成！"
