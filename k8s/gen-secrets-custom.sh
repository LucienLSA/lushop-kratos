#!/bin/bash

# Lushop K8s 自定义 Secrets 生成脚本
# 生成包含自定义密码的 K8s Secrets

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 配置
NAMESPACE="lushop"

# 默认密码 (可以修改)
MYSQL_ROOT_PASSWORD="${MYSQL_ROOT_PASSWORD:-MyRoot@123}"
MYSQL_DATABASE="${MYSQL_DATABASE:-lushop}"
MYSQL_USER="${MYSQL_USER:-lushop}"
MYSQL_PASSWORD="${MYSQL_PASSWORD:-lushopDb@123}"

REDIS_PASSWORD="${REDIS_PASSWORD:-Redis@123}"

NACOS_MYSQL_HOST="${NACOS_MYSQL_HOST:-mysql}"
NACOS_MYSQL_PORT="${NACOS_MYSQL_PORT:-3306}"
NACOS_MYSQL_DB="${NACOS_MYSQL_DB:-nacos}"
NACOS_MYSQL_USER="${NACOS_MYSQL_USER:-nacos}"
NACOS_MYSQL_PASSWORD="${NACOS_MYSQL_PASSWORD:-Nacos@123}"

ROCKETMQ_USERNAME="${ROCKETMQ_USERNAME:-rocketmq}"
ROCKETMQ_PASSWORD="${ROCKETMQ_PASSWORD:-Rmq@123}"

GRAFANA_ADMIN_USER="${GRAFANA_ADMIN_USER:-admin}"
GRAFANA_ADMIN_PASSWORD="${GRAFANA_ADMIN_PASSWORD:-Grafana@123}"

ELASTICSEARCH_PASSWORD="${ELASTICSEARCH_PASSWORD:-Elastic@123}"

log_info() {
    echo -e "${BLUE}[INFO]${NC} $(date +'%Y-%m-%d %H:%M:%S') - $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $(date +'%Y-%m-%d %H:%M:%S') - $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $(date +'%Y-%m-%d %H:%M:%S') - $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $(date +'%Y-%m-%d %H:%M:%S') - $1"
}

# 检查 kubectl
check_kubectl() {
    if ! command -v kubectl &> /dev/null; then
        log_error "kubectl 未安装"
        exit 1
    fi

    if ! kubectl cluster-info &> /dev/null; then
        log_error "无法连接到 K8s 集群"
        exit 1
    fi
}

# 创建命名空间
create_namespace() {
    if ! kubectl get namespace "$NAMESPACE" &> /dev/null; then
        kubectl create namespace "$NAMESPACE"
        log_success "命名空间 $NAMESPACE 已创建"
    else
        log_info "命名空间 $NAMESPACE 已存在"
    fi
}

# 生成 MySQL Secret
generate_mysql_secret() {
    log_info "生成 MySQL Secret..."

    kubectl create secret generic mysql-auth \
        --from-literal=mysql-root-password="$MYSQL_ROOT_PASSWORD" \
        --from-literal=mysql-database="$MYSQL_DATABASE" \
        --from-literal=mysql-user="$MYSQL_USER" \
        --from-literal=mysql-password="$MYSQL_PASSWORD" \
        -n "$NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -

    log_success "MySQL Secret 已创建"
    echo "  Root Password: $MYSQL_ROOT_PASSWORD"
    echo "  Database: $MYSQL_DATABASE"
    echo "  User: $MYSQL_USER"
    echo "  Password: $MYSQL_PASSWORD"
}

# 生成 Redis Secret
generate_redis_secret() {
    log_info "生成 Redis Secret..."

    kubectl create secret generic redis-auth \
        --from-literal=redis-password="$REDIS_PASSWORD" \
        -n "$NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -

    log_success "Redis Secret 已创建"
    echo "  Password: $REDIS_PASSWORD"
}

# 生成 Nacos Secret
generate_nacos_secret() {
    log_info "生成 Nacos Secret..."

    kubectl create secret generic nacos-auth \
        --from-literal=mysql-host="$NACOS_MYSQL_HOST" \
        --from-literal=mysql-port="$NACOS_MYSQL_PORT" \
        --from-literal=mysql-db="$NACOS_MYSQL_DB" \
        --from-literal=mysql-user="$NACOS_MYSQL_USER" \
        --from-literal=mysql-password="$NACOS_MYSQL_PASSWORD" \
        -n "$NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -

    log_success "Nacos Secret 已创建"
    echo "  MySQL Host: $NACOS_MYSQL_HOST:$NACOS_MYSQL_PORT"
    echo "  Database: $NACOS_MYSQL_DB"
    echo "  User: $NACOS_MYSQL_USER"
    echo "  Password: $NACOS_MYSQL_PASSWORD"
}

# 生成 RocketMQ Secret
generate_rocketmq_secret() {
    log_info "生成 RocketMQ Secret..."

    kubectl create secret generic rocketmq-credentials \
        --from-literal=username="$ROCKETMQ_USERNAME" \
        --from-literal=password="$ROCKETMQ_PASSWORD" \
        -n "$NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -

    log_success "RocketMQ Secret 已创建"
    echo "  Username: $ROCKETMQ_USERNAME"
    echo "  Password: $ROCKETMQ_PASSWORD"
}

# 生成 Grafana Secret
generate_grafana_secret() {
    log_info "生成 Grafana Secret..."

    kubectl create secret generic grafana-admin \
        --from-literal=admin-user="$GRAFANA_ADMIN_USER" \
        --from-literal=admin-password="$GRAFANA_ADMIN_PASSWORD" \
        -n "$NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -

    log_success "Grafana Secret 已创建"
    echo "  Admin User: $GRAFANA_ADMIN_USER"
    echo "  Admin Password: $GRAFANA_ADMIN_PASSWORD"
}

# 生成 Elasticsearch Secret
generate_elasticsearch_secret() {
    log_info "生成 Elasticsearch Secret..."

    kubectl create secret generic elasticsearch-auth \
        --from-literal=password="$ELASTICSEARCH_PASSWORD" \
        -n "$NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -

    log_success "Elasticsearch Secret 已创建"
    echo "  Password: $ELASTICSEARCH_PASSWORD"
}

# 显示所有生成的密码
show_passwords() {
    echo ""
    log_info "=== 生成的密码摘要 ==="
    echo "MySQL:"
    echo "  Root Password: $MYSQL_ROOT_PASSWORD"
    echo "  User Password: $MYSQL_PASSWORD"
    echo "Redis:"
    echo "  Password: $REDIS_PASSWORD"
    echo "Nacos:"
    echo "  MySQL Password: $NACOS_MYSQL_PASSWORD"
    echo "RocketMQ:"
    echo "  Password: $ROCKETMQ_PASSWORD"
    echo "Grafana:"
    echo "  Admin Password: $GRAFANA_ADMIN_PASSWORD"
    echo "Elasticsearch:"
    echo "  Password: $ELASTICSEARCH_PASSWORD"
    echo ""
    log_warning "请保存这些密码，用于后续配置和访问"
}

# 验证 Secrets
verify_secrets() {
    log_info "验证 Secrets 创建..."

    local secrets=("mysql-auth" "redis-auth" "nacos-auth" "rocketmq-credentials" "grafana-admin" "elasticsearch-auth")
    local failed=0

    for secret in "${secrets[@]}"; do
        if kubectl get secret "$secret" -n "$NAMESPACE" &> /dev/null; then
            log_success "Secret $secret 已创建"
        else
            log_error "Secret $secret 创建失败"
            ((failed++))
        fi
    done

    if [ $failed -gt 0 ]; then
        log_error "有 $failed 个 Secret 创建失败"
        return 1
    else
        log_success "所有 Secrets 验证通过"
    fi
}

# 主函数
main() {
    echo "=========================================="
    echo "  Lushop K8s 自定义 Secrets 生成脚本"
    echo "=========================================="
    echo ""

    # 显示配置
    echo "将生成以下密码 (可通过环境变量覆盖):"
    echo "MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD}"
    echo "MYSQL_PASSWORD=${MYSQL_PASSWORD}"
    echo "REDIS_PASSWORD=${REDIS_PASSWORD}"
    echo "NACOS_MYSQL_PASSWORD=${NACOS_MYSQL_PASSWORD}"
    echo "ROCKETMQ_PASSWORD=${ROCKETMQ_PASSWORD}"
    echo "GRAFANA_ADMIN_PASSWORD=${GRAFANA_ADMIN_PASSWORD}"
    echo "ELASTICSEARCH_PASSWORD=${ELASTICSEARCH_PASSWORD}"
    echo ""

    read -p "确认生成 Secrets? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log_info "已取消"
        exit 0
    fi

    check_kubectl
    create_namespace
    generate_mysql_secret
    generate_redis_secret
    generate_nacos_secret
    generate_rocketmq_secret
    generate_grafana_secret
    generate_elasticsearch_secret
    verify_secrets
    show_passwords

    echo ""
    log_success "所有 Secrets 已生成完成！"
    log_info "查看所有 Secrets: kubectl get secrets -n $NAMESPACE"
}

main "$@"
