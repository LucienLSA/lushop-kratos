#!/usr/bin/env bash
set -euo pipefail

# 镜像仓库配置
REGISTRY="${REGISTRY:-crpi-t2nj22cqo7df6n3a.cn-shanghai.personal.cr.aliyuncs.com/k8s_study6}"

REDIS_VERSION="${REDIS_VERSION:-7-alpine}"
REDIS_PORT="${REDIS_PORT:-6379}"
REDIS_NAME="${REDIS_NAME:-lushop-redis}"
REDIS_PASSWORD="${REDIS_PASSWORD:-123456}"
REDIS_MAXMEMORY="${REDIS_MAXMEMORY:-256mb}"
REDIS_MAXMEMORY_POLICY="${REDIS_MAXMEMORY_POLICY:-allkeys-lru}"
REDIS_TCP_KEEPALIVE="${REDIS_TCP_KEEPALIVE:-300}"
DATA_ROOT="${DATA_DIR:-$HOME/lushop-data}"
REDIS_DIR="${DATA_ROOT}/redis"

mkdir -p "${REDIS_DIR}"

CONF_FILE="${REDIS_DIR}/redis.conf"
if [ ! -f "${CONF_FILE}" ]; then
  cat > "${CONF_FILE}" <<EOF
bind 0.0.0.0
protected-mode no
port 6379
daemonize no
tcp-backlog 511
timeout 0
tcp-keepalive ${REDIS_TCP_KEEPALIVE}
loglevel notice
logfile ""
databases 16
save 900 1
save 300 10
save 60 10000
stop-writes-on-bgsave-error yes
rdbcompression yes
rdbchecksum yes
dbfilename dump.rdb
dir /data
maxmemory ${REDIS_MAXMEMORY}
maxmemory-policy ${REDIS_MAXMEMORY_POLICY}
requirepass ${REDIS_PASSWORD}
appendonly yes
appendfsync everysec
EOF
fi

docker run -d \
  --name "${REDIS_NAME}" \
  --restart unless-stopped \
  -p "${REDIS_PORT}:6379" \
  -e REDIS_PASSWORD="${REDIS_PASSWORD}" \
  -v "${CONF_FILE}":/etc/redis/redis.conf \
  -v "${REDIS_DIR}":/data \
  "${REGISTRY}/redis:${REDIS_VERSION}" \
  redis-server /etc/redis/redis.conf --appendonly yes --requirepass "${REDIS_PASSWORD}"