#!/usr/bin/env bash
set -euo pipefail

# 镜像仓库配置
REGISTRY="${REGISTRY:-crpi-t2nj22cqo7df6n3a.cn-shanghai.personal.cr.aliyuncs.com/k8s_study6}"

CONSUL_VERSION="${CONSUL_VERSION:-1.16}"
CONSUL_CONTAINER_NAME="${CONSUL_CONTAINER_NAME:-lushop-consul}"
CONSUL_HTTP_PORT="${CONSUL_HTTP_PORT:-8500}"
CONSUL_SERF_LAN_PORT="${CONSUL_SERF_LAN_PORT:-8301}"
CONSUL_SERF_WAN_PORT="${CONSUL_SERF_WAN_PORT:-8302}"
CONSUL_SERVER_PORT="${CONSUL_SERVER_PORT:-8300}"
CONSUL_DNS_PORT="${CONSUL_DNS_PORT:-8600}"

docker run -d \
  --name "${CONSUL_CONTAINER_NAME}" \
  --restart unless-stopped \
  -p "${CONSUL_HTTP_PORT:-8500}:8500" \
  -p "${CONSUL_SERF_LAN_PORT:-8301}:8301" \
  -p "${CONSUL_SERF_WAN_PORT:-8302}:8302" \
  -p "${CONSUL_SERVER_PORT:-8300}:8300" \
  -p "${CONSUL_DNS_PORT:-8600}:8600/udp" \
  "${REGISTRY}/consul:${CONSUL_VERSION}" agent -dev -client=0.0.0.0