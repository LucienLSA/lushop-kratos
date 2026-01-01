#!/usr/bin/env bash
set -euo pipefail

# 镜像仓库配置
REGISTRY="${REGISTRY:-crpi-t2nj22cqo7df6n3a.cn-shanghai.personal.cr.aliyuncs.com/k8s_study6}"

JAEGER_VERSION="${JAEGER_VERSION:-1.52}"
JAEGER_CONTAINER_NAME="${JAEGER_CONTAINER_NAME:-lushop-jaeger}"
JAEGER_UI_PORT="${JAEGER_UI_PORT:-16686}"
JAEGER_COLLECTOR_PORT="${JAEGER_COLLECTOR_PORT:-14268}"
JAEGER_COLLECTOR_GRPC_PORT="${JAEGER_COLLECTOR_GRPC_PORT:-14250}"

docker run -d \
  --name "${JAEGER_CONTAINER_NAME}" \
  --restart unless-stopped \
  -p "${JAEGER_UI_PORT}:16686" \
  -p "${JAEGER_COLLECTOR_PORT}:14268" \
  -p "${JAEGER_COLLECTOR_GRPC_PORT}:14250" \
  "${REGISTRY}/jaeger:${JAEGER_VERSION}" \
  --memory.max-traces=50000 \
  --query.max-clock-skew-adjustment=0
