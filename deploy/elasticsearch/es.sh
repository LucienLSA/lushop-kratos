#!/usr/bin/env bash
set -euo pipefail

# 镜像仓库配置
REGISTRY="${REGISTRY:-crpi-t2nj22cqo7df6n3a.cn-shanghai.personal.cr.aliyuncs.com/k8s_study6}"

ES_VERSION="${ELASTICSEARCH_VERSION:-7.10.1}"
ES_CONTAINER_NAME="${ELASTICSEARCH_CONTAINER_NAME:-lushop-elasticsearch}"
ELASTICSEARCH_HTTP_PORT="${ELASTICSEARCH_HTTP_PORT:-9200}"
ELASTICSEARCH_TRANSPORT_PORT="${ELASTICSEARCH_TRANSPORT_PORT:-9300}"
ES_JAVA_OPTS="${ES_JAVA_OPTS:--Xms512m -Xmx512m}"
DATA_ROOT="${DATA_DIR:-$HOME/lushop-data}"
ES_BASE_DIR="${DATA_ROOT}/elasticsearch"
ES_CONFIG_DIR="${ES_BASE_DIR}/config"
ES_DATA_DIR="${ES_BASE_DIR}/data"
ES_PLUGINS_DIR="${ES_BASE_DIR}/plugins"

mkdir -p "${ES_CONFIG_DIR}" "${ES_DATA_DIR}" "${ES_PLUGINS_DIR}"

CONFIG_FILE="${ES_CONFIG_DIR}/elasticsearch.yml"
if [ ! -f "${CONFIG_FILE}" ]; then
  cat > "${CONFIG_FILE}" <<'EOF'
cluster.name: lushop-es
node.name: lushop-es-node
path.data: /usr/share/elasticsearch/data
path.logs: /usr/share/elasticsearch/logs
network.host: 0.0.0.0
http.port: 9200
discovery.type: single-node
bootstrap.memory_lock: true
EOF
fi

docker run -d \
  --name "${ES_CONTAINER_NAME}" \
  --restart unless-stopped \
  -p "${ELASTICSEARCH_HTTP_PORT:-9200}:9200" \
  -p "${ELASTICSEARCH_TRANSPORT_PORT:-9300}:9300" \
  -e ES_JAVA_OPTS="${ES_JAVA_OPTS:--Xms512m -Xmx512m}" \
  -e "discovery.type=single-node" \
  -v "${ES_DATA_DIR}":/usr/share/elasticsearch/data \
  -v "${ES_PLUGINS_DIR}":/usr/share/elasticsearch/plugins \
  -v "${CONFIG_FILE}":/usr/share/elasticsearch/config/elasticsearch.yml \
  "${REGISTRY}/elasticsearch:${ES_VERSION}"