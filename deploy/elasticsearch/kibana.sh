#!/usr/bin/env bash
set -euo pipefail

KIBANA_VERSION="${KIBANA_VERSION:-7.10.1}"
KIBANA_CONTAINER_NAME="${KIBANA_CONTAINER_NAME:-lushop-kibana}"
ELASTICSEARCH_HOST="${ELASTICSEARCH_HOST:-http://localhost:9200}"

docker run -d \
  --name "${KIBANA_CONTAINER_NAME}" \
  --restart unless-stopped \
  -e "ELASTICSEARCH_HOSTS=${ELASTICSEARCH_HOST}" \
  -p "${KIBANA_PORT:-5601}:5601" \
  "docker.elastic.co/kibana/kibana:${KIBANA_VERSION}"