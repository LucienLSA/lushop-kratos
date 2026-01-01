#!/usr/bin/env bash
set -euo pipefail

# 镜像仓库配置
REGISTRY="${REGISTRY:-crpi-t2nj22cqo7df6n3a.cn-shanghai.personal.cr.aliyuncs.com/k8s_study6}"

# 版本和配置
GRAFANA_VERSION="${GRAFANA_VERSION:-10.3.3}"
GRAFANA_CONTAINER_NAME="${GRAFANA_CONTAINER_NAME:-lushop-grafana}"
GRAFANA_PORT="${GRAFANA_PORT:-3000}"
GRAFANA_ADMIN_PASSWORD="${GRAFANA_ADMIN_PASSWORD:-admin}"
GRAFANA_ADMIN_USER="${GRAFANA_ADMIN_USER:-admin}"

DATA_ROOT="${DATA_DIR:-$HOME/lushop-data}"
GRAFANA_DATA_DIR="${DATA_ROOT}/grafana"
mkdir -p "${GRAFANA_DATA_DIR}"

docker run -d \
  --name "${GRAFANA_CONTAINER_NAME}" \
  --restart unless-stopped \
  -p "${GRAFANA_PORT}:3000" \
  -v "${GRAFANA_DATA_DIR}":/var/lib/grafana \
  -e "GF_SECURITY_ADMIN_USER=${GRAFANA_ADMIN_USER}" \
  -e "GF_SECURITY_ADMIN_PASSWORD=${GRAFANA_ADMIN_PASSWORD}" \
  -e "GF_USERS_ALLOW_SIGN_UP=false" \
  -e "GF_INSTALL_PLUGINS=" \
  "${REGISTRY}/grafana:${GRAFANA_VERSION}"