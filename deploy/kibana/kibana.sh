#!/usr/bin/env bash
set -euo pipefail

# 镜像仓库配置
REGISTRY="${REGISTRY:-crpi-t2nj22cqo7df6n3a.cn-shanghai.personal.cr.aliyuncs.com/k8s_study6}"

KIBANA_VERSION="${KIBANA_VERSION:-7.10.1}"
KIBANA_CONTAINER_NAME="${KIBANA_CONTAINER_NAME:-lushop-kibana}"
KIBANA_PORT="${KIBANA_PORT:-5601}"
ELASTICSEARCH_HOST="${ELASTICSEARCH_HOST:-http://localhost:9200}"
KIBANA_INDEX="${KIBANA_INDEX:-.kibana}"

docker run -d \
  --name "${KIBANA_CONTAINER_NAME}" \
  --restart unless-stopped \
  -e "ELASTICSEARCH_HOSTS=${ELASTICSEARCH_HOST}" \
  -e "KIBANA_INDEX=${KIBANA_INDEX}" \
  -p "${KIBANA_PORT}:5601" \
  "${REGISTRY}/kibana:${KIBANA_VERSION}"