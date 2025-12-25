#!/usr/bin/env bash
set -euo pipefail

# 配置 - 如需修改请在脚本顶部调整
REGISTRY="crpi-t2nj22cqo7df6n3a.cn-shanghai.personal.cr.aliyuncs.com"
NAMESPACE="aliyun1123466419"
TAG="${TAG:-$(git rev-parse --short HEAD)}"

# 服务列表（按仓库实际服务名调整）
services=(lushop goods inventory order user userauth userop)

read -s -p "输入阿里云仓库密码并回车: " ACR_PWD
echo
echo "登录 ${REGISTRY} ..."
echo "$ACR_PWD" | docker login --username="aliyun1123466419" --password-stdin "$REGISTRY"

for svc in "${services[@]}"; do
  if [ -f "service/${svc}/Dockerfile" ]; then
    DOCKERFILE="service/${svc}/Dockerfile"
    CONTEXT="service/${svc}"
  elif [ -f "${svc}/Dockerfile" ]; then
    DOCKERFILE="${svc}/Dockerfile"
    CONTEXT="${svc}"
  elif [ "${svc}" = "lushop" ] && [ -f "lushop/Dockerfile" ]; then
    DOCKERFILE="lushop/Dockerfile"
    CONTEXT="lushop"
  else
    echo "跳过 ${svc}：未找到 Dockerfile"
    continue
  fi

  IMAGE="${REGISTRY}/${NAMESPACE}/${svc}:${TAG}"
  echo "构建 ${svc} -> ${IMAGE}"
  docker build -t "${IMAGE}" -f "${DOCKERFILE}" "${CONTEXT}"

  echo "推送 ${IMAGE}"
  docker push "${IMAGE}"
done

echo "构建与推送完成。TAG=${TAG}"


