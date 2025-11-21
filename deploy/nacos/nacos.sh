#!/usr/bin/env bash
set -euo pipefail

NACOS_VERSION="${NACOS_VERSION:-2.3.2}"
NACOS_CONTAINER_NAME="${NACOS_CONTAINER_NAME:-lushop-nacos}"
DATA_ROOT="${DATA_DIR:-$HOME/lushop-data}"
NACOS_ROOT="${DATA_ROOT}/nacos"

mkdir -p "${NACOS_ROOT}/"{conf,data,logs}

docker run -d \
  --name "${NACOS_CONTAINER_NAME}" \
  --restart unless-stopped \
  -p "${NACOS_HTTP_PORT:-8848}:8848" \
  -p "${NACOS_GRPC_PORT:-9848}:9848" \
  -p "${NACOS_GRPCS_PORT:-9849}:9849" \
  -e MODE="${NACOS_MODE:-standalone}" \
  -e SPRING_DATASOURCE_PLATFORM="${SPRING_DATASOURCE_PLATFORM:-mysql}" \
  -e MYSQL_SERVICE_HOST="${NACOS_DB_HOST:-localhost}" \
  -e MYSQL_SERVICE_PORT="${NACOS_DB_PORT:-3306}" \
  -e MYSQL_SERVICE_USER="${NACOS_DB_USER:-root}" \
  -e MYSQL_SERVICE_PASSWORD="${NACOS_DB_PASSWORD:-root123456}" \
  -e MYSQL_SERVICE_DB_NAME="${NACOS_DB_NAME:-nacos}" \
  -e TIME_ZONE="${NACOS_TIMEZONE:-Asia/Shanghai}" \
  -e NACOS_AUTH_ENABLE="${NACOS_AUTH_ENABLE:-true}" \
  -e JVM_XMS="${NACOS_JVM_XMS:-256m}" \
  -e JVM_XMX="${NACOS_JVM_XMX:-256m}" \
  -v "${NACOS_ROOT}/logs":/home/nacos/logs \
  -v "${NACOS_ROOT}/data":/home/nacos/data \
  -v "${NACOS_ROOT}/conf":/home/nacos/conf \
  "nacos/nacos-server:${NACOS_VERSION}"