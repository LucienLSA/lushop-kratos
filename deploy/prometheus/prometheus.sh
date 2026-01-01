#!/usr/bin/env bash
set -euo pipefail

# 镜像仓库配置
REGISTRY="${REGISTRY:-crpi-t2nj22cqo7df6n3a.cn-shanghai.personal.cr.aliyuncs.com/k8s_study6}"

# 版本和配置
PROMETHEUS_VERSION="${PROMETHEUS_VERSION:-v2.52.0}"
PROMETHEUS_CONTAINER_NAME="${PROMETHEUS_CONTAINER_NAME:-lushop-prometheus}"
PROMETHEUS_PORT="${PROMETHEUS_PORT:-9090}"

SCRIPT_DIR="$(cd -- "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DATA_ROOT="${DATA_DIR:-$HOME/lushop-data}"
PROM_DATA="${DATA_ROOT}/prometheus"
PROM_CONF_DIR="${PROM_DATA}/conf"
PROM_RULES_DIR="${PROM_DATA}/rules"
PROM_STORAGE_DIR="${PROM_DATA}/data"
mkdir -p "${PROM_CONF_DIR}" "${PROM_RULES_DIR}" "${PROM_STORAGE_DIR}"

DEFAULT_CONF_SOURCE="${PROMETHEUS_CONFIG:-${SCRIPT_DIR}/conf/prometheus.yaml}"
if [ ! -f "${DEFAULT_CONF_SOURCE}" ]; then
  echo "Prometheus config not found: ${DEFAULT_CONF_SOURCE}" >&2
  exit 1
fi
cp "${DEFAULT_CONF_SOURCE}" "${PROM_CONF_DIR}/prometheus.yaml"

docker run -d \
  --name "${PROMETHEUS_CONTAINER_NAME}" \
  --restart unless-stopped \
  --privileged=true \
  -u root \
  -p "${PROMETHEUS_PORT}:9090" \
  -v /etc/localtime:/etc/localtime:ro \
  -v "${PROM_STORAGE_DIR}":/prometheus/data \
  -v "${PROM_CONF_DIR}":/prometheus/conf \
  -v "${PROM_RULES_DIR}":/prometheus/rules \
  "${REGISTRY}/prometheus:${PROMETHEUS_VERSION}" \
  --config.file=/prometheus/conf/prometheus.yaml \
  --storage.tsdb.path=/prometheus/data \
  --web.enable-lifecycle \
  --web.listen-address=0.0.0.0:9090