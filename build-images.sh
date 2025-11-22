#!/usr/bin/env bash
set -euo pipefail

# 获取脚本所在目录（项目根目录）
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 如果 k8s/build-images.sh 存在，则调用它
if [ -f "${SCRIPT_DIR}/k8s/build-images.sh" ]; then
    exec "${SCRIPT_DIR}/k8s/build-images.sh" "$@"
else
    echo "错误: 未找到 k8s/build-images.sh"
    exit 1
fi

