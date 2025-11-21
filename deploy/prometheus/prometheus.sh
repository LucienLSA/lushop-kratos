#!/usr/bin/env bash
set -euo pipefail

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
  --name "${PROMETHEUS_CONTAINER_NAME:-lushop-prometheus}" \
  --restart unless-stopped \
  --privileged=true \
  -u root \
  -p "${PROMETHEUS_PORT:-9090}:9090" \
  -v /etc/localtime:/etc/localtime:ro \
  -v "${PROM_STORAGE_DIR}":/prometheus/data \
  -v "${PROM_CONF_DIR}":/prometheus/conf \
  -v "${PROM_RULES_DIR}":/prometheus/rules \
  prom/prometheus:${PROMETHEUS_VERSION:-v2.52.0} \
  --config.file=/prometheus/conf/prometheus.yaml \
  --storage.tsdb.path=/prometheus/data \
  --web.enable-lifecycle