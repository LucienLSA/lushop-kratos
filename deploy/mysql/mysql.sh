#!/usr/bin/env bash
set -euo pipefail

# 镜像仓库配置
REGISTRY="${REGISTRY:-crpi-t2nj22cqo7df6n3a.cn-shanghai.personal.cr.aliyuncs.com/k8s_study6}"

MYSQL_VERSION="${MYSQL_VERSION:-8.0}"
MYSQL_CONTAINER_NAME="${MYSQL_CONTAINER_NAME:-lushop-mysql}"
MYSQL_PORT="${MYSQL_PORT:-3306}"
MYSQL_ROOT_PASSWORD="${MYSQL_ROOT_PASSWORD:-root123456}"
MYSQL_DATABASE="${MYSQL_DATABASE:-lushop}"
MYSQL_USER="${MYSQL_USER:-lushop}"
MYSQL_PASSWORD="${MYSQL_PASSWORD:-lushop123456}"
MYSQL_TIMEZONE="${MYSQL_TIMEZONE:-Asia/Shanghai}"
MYSQL_CHARACTER_SET="${MYSQL_CHARACTER_SET:-utf8mb4}"
MYSQL_COLLATION="${MYSQL_COLLATION:-utf8mb4_unicode_ci}"
DATA_ROOT="${DATA_DIR:-$HOME/lushop-data}"
MYSQL_DATA_DIR="${DATA_ROOT}/mysql"

mkdir -p "${MYSQL_DATA_DIR}"

docker run -d \
  --name "${MYSQL_CONTAINER_NAME}" \
  --restart unless-stopped \
  -e MYSQL_ROOT_PASSWORD="${MYSQL_ROOT_PASSWORD}" \
  -e MYSQL_DATABASE="${MYSQL_DATABASE}" \
  -e MYSQL_USER="${MYSQL_USER}" \
  -e MYSQL_PASSWORD="${MYSQL_PASSWORD}" \
  -e TZ="${MYSQL_TIMEZONE}" \
  -p "${MYSQL_PORT}:3306" \
  -v "${MYSQL_DATA_DIR}":/var/lib/mysql \
  "${REGISTRY}/mysql:${MYSQL_VERSION}" \
  mysqld --character-set-server="${MYSQL_CHARACTER_SET}" --collation-server="${MYSQL_COLLATION}"