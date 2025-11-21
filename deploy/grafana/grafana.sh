#!/usr/bin/env bash
set -euo pipefail

DATA_ROOT="${DATA_DIR:-$HOME/lushop-data}"
GRAFANA_DATA_DIR="${DATA_ROOT}/grafana"
mkdir -p "${GRAFANA_DATA_DIR}"

docker run -d \
  --name "${GRAFANA_CONTAINER_NAME:-lushop-grafana}" \
  --restart unless-stopped \
  -p "${GRAFANA_PORT:-3000}:3000" \
  -v "${GRAFANA_DATA_DIR}":/var/lib/grafana \
  -e "GF_SECURITY_ADMIN_PASSWORD=${GRAFANA_ADMIN_PASSWORD:-admin}" \
  grafana/grafana:${GRAFANA_VERSION:-10.3.3}