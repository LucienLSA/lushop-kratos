#!/bin/bash

# 清理重复的 Lushop Docker 镜像
# 保留 lushop/* 系列，删除其他重复镜像

set -e

echo "=== Lushop Docker 镜像清理脚本 ==="
echo ""

# 显示当前镜像
echo "当前 lushop 相关镜像:"
docker images | grep lushop
echo ""

# 保留的镜像系列
KEEP_PATTERN="lushop/"
echo "将保留以下镜像系列: $KEEP_PATTERN"
echo ""

# 查找需要删除的镜像
IMAGES_TO_REMOVE=$(docker images | grep "lushop" | grep -v "$KEEP_PATTERN" | awk '{print $1":"$2}')

if [ -z "$IMAGES_TO_REMOVE" ]; then
    echo "没有找到重复镜像，无需清理"
    exit 0
fi

echo "将删除以下重复镜像:"
echo "$IMAGES_TO_REMOVE" | while read -r image; do
    echo "  - $image"
done
echo ""

read -p "确认删除这些重复镜像? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "已取消清理"
    exit 0
fi

echo ""
echo "开始删除镜像..."

# 删除镜像
echo "$IMAGES_TO_REMOVE" | while read -r image; do
    if [ -n "$image" ]; then
        echo "删除: $image"
        docker rmi "$image" 2>/dev/null || echo "  跳过: $image (可能已被删除)"
    fi
done

echo ""
echo "清理完成！剩余镜像:"
docker images | grep lushop
