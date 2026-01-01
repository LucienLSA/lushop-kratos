#!/usr/bin/env bash
# Docker 部署环境变量配置文件模板
# 复制此文件的内容到 .env 文件并修改相应的值
# 注意：不要将包含真实密码的 .env 文件提交到版本控制系统

cat << 'EOF'
# =============================================================================
# 全局配置
# =============================================================================

# 镜像仓库配置
REGISTRY=crpi-t2nj22cqo7df6n3a.cn-shanghai.personal.cr.aliyuncs.com/k8s_study6

# 数据根目录
DATA_DIR=$HOME/lushop-data

# =============================================================================
# MySQL 配置
# =============================================================================

MYSQL_VERSION=8.0
MYSQL_CONTAINER_NAME=lushop-mysql
MYSQL_PORT=3306
MYSQL_ROOT_PASSWORD=your_secure_root_password
MYSQL_DATABASE=lushop
MYSQL_USER=lushop
MYSQL_PASSWORD=your_secure_mysql_password
MYSQL_TIMEZONE=Asia/Shanghai
MYSQL_CHARACTER_SET=utf8mb4
MYSQL_COLLATION=utf8mb4_unicode_ci

# =============================================================================
# Redis 配置
# =============================================================================

REDIS_VERSION=7-alpine
REDIS_PORT=6379
REDIS_NAME=lushop-redis
REDIS_PASSWORD=your_secure_redis_password
REDIS_MAXMEMORY=256mb
REDIS_MAXMEMORY_POLICY=allkeys-lru
REDIS_TCP_KEEPALIVE=300

# =============================================================================
# Nacos 配置
# =============================================================================

NACOS_VERSION=v2.3.2
NACOS_CONTAINER_NAME=lushop-nacos
NACOS_HTTP_PORT=8848
NACOS_GRPC_PORT=9848
NACOS_GRPCS_PORT=9849
NACOS_MODE=standalone
NACOS_DB_HOST=localhost
NACOS_DB_PORT=3306
NACOS_DB_USER=root
NACOS_DB_PASSWORD=your_secure_root_password
NACOS_DB_NAME=nacos
NACOS_TIMEZONE=Asia/Shanghai
NACOS_AUTH_ENABLE=true
NACOS_JVM_XMS=256m
NACOS_JVM_XMX=256m

# =============================================================================
# Consul 配置
# =============================================================================

CONSUL_VERSION=1.16
CONSUL_CONTAINER_NAME=lushop-consul
CONSUL_HTTP_PORT=8500
CONSUL_SERF_LAN_PORT=8301
CONSUL_SERF_WAN_PORT=8302
CONSUL_SERVER_PORT=8300
CONSUL_DNS_PORT=8600

# =============================================================================
# Elasticsearch 配置
# =============================================================================

ELASTICSEARCH_VERSION=7.10.1
ELASTICSEARCH_CONTAINER_NAME=lushop-elasticsearch
ELASTICSEARCH_HTTP_PORT=9200
ELASTICSEARCH_TRANSPORT_PORT=9300
ES_JAVA_OPTS="-Xms512m -Xmx512m"

# =============================================================================
# Kibana 配置
# =============================================================================

KIBANA_VERSION=7.10.1
KIBANA_CONTAINER_NAME=lushop-kibana
KIBANA_PORT=5601
ELASTICSEARCH_HOST=http://localhost:9200
KIBANA_INDEX=.kibana

# =============================================================================
# Grafana 配置
# =============================================================================

GRAFANA_VERSION=10.3.3
GRAFANA_CONTAINER_NAME=lushop-grafana
GRAFANA_PORT=3000
GRAFANA_ADMIN_USER=admin
GRAFANA_ADMIN_PASSWORD=your_secure_grafana_password

# =============================================================================
# Prometheus 配置
# =============================================================================

PROMETHEUS_VERSION=v2.52.0
PROMETHEUS_CONTAINER_NAME=lushop-prometheus
PROMETHEUS_PORT=9090

# =============================================================================
# Jaeger 配置
# =============================================================================

JAEGER_VERSION=1.52
JAEGER_CONTAINER_NAME=lushop-jaeger
JAEGER_UI_PORT=16686
JAEGER_COLLECTOR_PORT=14268
JAEGER_COLLECTOR_GRPC_PORT=14250

# =============================================================================
# RocketMQ 配置 (用于 docker-compose)
# =============================================================================

ROCKETMQ_VERSION=5.3.3
ROCKETMQ_DATA_DIR=./data
ROCKETMQ_ACCESS_KEY=lushop
ROCKETMQ_SECRET_KEY=your_secure_rocketmq_secret
ROCKETMQ_DASHBOARD_PASSWORD=your_secure_dashboard_password
ROCKETMQ_ACL_ENABLE=true
ROCKETMQ_AUTH_ENABLE=true
ROCKETMQ_AUTHZ_ENABLE=true

# =============================================================================
# 使用说明
# =============================================================================
#
# 1. 创建 .env 文件:
#    cp deploy/env-template.sh /tmp/env-template.sh
#    bash /tmp/env-template.sh > .env
#
# 2. 编辑 .env 文件，修改所有密码和敏感配置
#
# 3. 运行脚本时会自动加载 .env 文件中的变量:
#    ./deploy/mysql/mysql.sh
#
# 4. 或者手动设置环境变量:
#    MYSQL_ROOT_PASSWORD=xxx ./deploy/mysql/mysql.sh
#
# 安全提醒:
# - 所有默认密码都应该修改
# - 生产环境使用强密码
# - 定期轮换密码
# - 不要将 .env 文件提交到版本控制系统

EOF
