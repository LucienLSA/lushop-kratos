#!/bin/bash

# 简单的压力测试脚本
set -e

BASE_URL="http://localhost:8001"

# 检查 ab 命令
if ! command -v ab &> /dev/null; then
    echo "❌ 请先安装 Apache Bench:"
    echo "   sudo apt-get install apache2-utils"
    exit 1
fi

# 检查服务
echo "检查服务状态..."
if ! curl -s "$BASE_URL/metrics" > /dev/null 2>&1; then
    echo "❌ API 网关未运行"
    echo "请先运行: cd .. && ./deploy.sh start"
    exit 1
fi
echo "✅ API 网关运行正常"
echo ""

mkdir -p reports/logs

# 测试1: 商品列表
echo "=========================================="
echo "测试 1/3: 商品列表"
echo "=========================================="
ab -n 1000 -c 100 "$BASE_URL/api/goods/list?page=1&pageSize=20" > reports/logs/goods_list.txt 2>&1
cat reports/logs/goods_list.txt | grep -E "(Requests per second|Time per request|Failed)"
echo ""

sleep 1

# 测试2: 商品详情
echo "=========================================="
echo "测试 2/3: 商品详情"
echo "=========================================="
ab -n 2000 -c 200 "$BASE_URL/api/goods/1" > reports/logs/goods_detail.txt 2>&1
cat reports/logs/goods_detail.txt | grep -E "(Requests per second|Time per request|Failed)"
echo ""

sleep 1

# 测试3: 库存查询
echo "=========================================="
echo "测试 3/3: 库存查询"
echo "=========================================="
ab -n 3000 -c 300 "$BASE_URL/api/inventory/1" > reports/logs/inventory.txt 2>&1
cat reports/logs/inventory.txt | grep -E "(Requests per second|Time per request|Failed)"
echo ""

echo "=========================================="
echo "✅ 压力测试完成"
echo "=========================================="
echo ""
echo "详细报告请查看: reports/logs/"

