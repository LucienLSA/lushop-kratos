#!/usr/bin/env bash
set -euo pipefail

MYSQL_VERSION="${MYSQL_VERSION:-8.0}"
MYSQL_CONTAINER_NAME="${MYSQL_CONTAINER_NAME:-lushop-mysql}"
MYSQL_PORT="${MYSQL_PORT:-3306}"
MYSQL_ROOT_PASSWORD="${MYSQL_ROOT_PASSWORD:-root123456}"
MYSQL_DATABASE="${MYSQL_DATABASE:-lushop}"
MYSQL_USER="${MYSQL_USER:-lushop}"
MYSQL_PASSWORD="${MYSQL_PASSWORD:-lushop123456}"
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
  -p "${MYSQL_PORT}:3306" \
  -v "${MYSQL_DATA_DIR}":/var/lib/mysql \
  "mysql:${MYSQL_VERSION}" \
  --character-set-server=utf8mb4 \
  --collation-server=utf8mb4_unicode_ci