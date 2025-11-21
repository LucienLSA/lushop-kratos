#!/usr/bin/env bash
set -euo pipefail

ROCKETMQ_VERSION="${ROCKETMQ_VERSION:-5.3.3}"
DATA_ROOT="${DATA_DIR:-$HOME/lushop-data}"
ROCKETMQ_BASE="${DATA_ROOT}/rocketmq/${ROCKETMQ_VERSION}"
BROKER_DIR="${ROCKETMQ_BASE}/broker"
PROXY_DIR="${ROCKETMQ_BASE}/proxy"
DASHBOARD_DIR="${ROCKETMQ_BASE}/dashboard"
BROKER_CONF_DIR="${BROKER_DIR}/conf"
PROXY_CONF_DIR="${PROXY_DIR}/conf"
DASHBOARD_CONF_DIR="${DASHBOARD_DIR}/conf"

mkdir -p \
  "${BROKER_CONF_DIR}" "${BROKER_DIR}/logs" "${BROKER_DIR}/store" \
  "${PROXY_CONF_DIR}" "${PROXY_DIR}/logs" \
  "${DASHBOARD_CONF_DIR}"

ACCESS_KEY="${ROCKETMQ_ACCESS_KEY:-lushop}"
SECRET_KEY="${ROCKETMQ_SECRET_KEY:-lushop123456}"
BROKER_IP="${ROCKETMQ_BROKER_IP:-127.0.0.1}"
BROKER_CLUSTER="${ROCKETMQ_CLUSTER_NAME:-DefaultCluster}"
BROKER_NAME="${ROCKETMQ_BROKER_NAME:-broker-a}"

cat > "${BROKER_CONF_DIR}/broker.conf" <<EOF
brokerClusterName=${BROKER_CLUSTER}
brokerName=${BROKER_NAME}
brokerId=0
deleteWhen=04
fileReservedTime=48
brokerRole=ASYNC_MASTER
flushDiskType=ASYNC_FLUSH
brokerIP1=${BROKER_IP}
listenPort=10911
autoCreateTopicEnable=true
aclEnable=${ROCKETMQ_ACL_ENABLE:-true}
EOF

cat > "${BROKER_CONF_DIR}/tools.yml" <<EOF
accessKey: ${ACCESS_KEY}
secretKey: ${SECRET_KEY}
EOF

cat > "${PROXY_CONF_DIR}/rmq-proxy.json" <<EOF
{
  "rocketMQClusterName": "${BROKER_CLUSTER}",
  "remotingListenPort": 18680,
  "grpcServerPort": 18681,
  "enableACL": ${ROCKETMQ_ACL_ENABLE:-true},
  "authenticationEnabled": ${ROCKETMQ_AUTH_ENABLE:-true},
  "innerClientAuthenticationCredentials": "{\"accessKey\":\"${ACCESS_KEY}\", \"secretKey\":\"${SECRET_KEY}\"}",
  "enableAclRpcHookForClusterMode": true,
  "authorizationEnabled": ${ROCKETMQ_AUTHZ_ENABLE:-true}
}
EOF

cat > "${DASHBOARD_CONF_DIR}/users.properties" <<EOF
admin=${ROCKETMQ_DASHBOARD_PASSWORD:-admin123}
EOF

chown -R 3000:3000 "${ROCKETMQ_BASE}" || true

cat <<EOF
RocketMQ data prepared under: ${ROCKETMQ_BASE}

Update the following defaults as needed before starting docker-compose:
  ACCESS_KEY=${ACCESS_KEY}
  SECRET_KEY=${SECRET_KEY}
  DASHBOARD_PASSWORD=$(cat "${DASHBOARD_CONF_DIR}/users.properties" | cut -d'=' -f2)
EOF