#!/usr/bin/env bash
set -euo pipefail

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

# 脚本目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
NAMESPACE="lushop"

info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# 函数：生成随机密码
generate_password() {
    openssl rand -base64 32 | tr -d "=+/" | cut -c1-25
}

# 函数：创建命名空间（如果不存在）
create_namespace() {
    if ! kubectl get namespace "$NAMESPACE" &> /dev/null; then
        kubectl create namespace "$NAMESPACE"
        info "命名空间 $NAMESPACE 已创建"
    fi
}

# 函数：生成 MySQL Secret
generate_mysql_secret() {
    step "生成 MySQL Secret..."
    
    local root_password="${MYSQL_ROOT_PASSWORD:-$(generate_password)}"
    local db_name="${MYSQL_DATABASE:-lushop}"
    local db_user="${MYSQL_USER:-lushop}"
    local db_password="${MYSQL_PASSWORD:-$(generate_password)}"
    
    if kubectl get secret mysql-auth -n "$NAMESPACE" &> /dev/null; then
        warn "Secret mysql-auth 已存在，跳过创建"
        return 0
    fi
    
    kubectl create secret generic mysql-auth \
        --from-literal=mysql-root-password="$root_password" \
        --from-literal=mysql-database="$db_name" \
        --from-literal=mysql-user="$db_user" \
        --from-literal=mysql-password="$db_password" \
        -n "$NAMESPACE" \
        --dry-run=client -o yaml | kubectl apply -f -
    
    info "MySQL Secret 已创建"
    echo "  Root Password: $root_password"
    echo "  Database: $db_name"
    echo "  User: $db_user"
    echo "  Password: $db_password"
    echo ""
}

# 函数：生成 Redis Secret
generate_redis_secret() {
    step "生成 Redis Secret..."
    
    local password="${REDIS_PASSWORD:-$(generate_password)}"
    
    if kubectl get secret redis-auth -n "$NAMESPACE" &> /dev/null; then
        warn "Secret redis-auth 已存在，跳过创建"
        return 0
    fi
    
    kubectl create secret generic redis-auth \
        --from-literal=redis-password="$password" \
        -n "$NAMESPACE" \
        --dry-run=client -o yaml | kubectl apply -f -
    
    info "Redis Secret 已创建"
    echo "  Password: $password"
    echo ""
}

# 函数：生成 Nacos Secret
generate_nacos_secret() {
    step "生成 Nacos Secret..."
    
    local mysql_host="${NACOS_MYSQL_HOST:-mysql}"
    local mysql_port="${NACOS_MYSQL_PORT:-3306}"
    local mysql_db="${NACOS_MYSQL_DB:-nacos}"
    local mysql_user="${NACOS_MYSQL_USER:-nacos}"
    local mysql_password="${NACOS_MYSQL_PASSWORD:-$(generate_password)}"
    
    if kubectl get secret nacos-auth -n "$NAMESPACE" &> /dev/null; then
        warn "Secret nacos-auth 已存在，跳过创建"
        return 0
    fi
    
    kubectl create secret generic nacos-auth \
        --from-literal=mysql-host="$mysql_host" \
        --from-literal=mysql-port="$mysql_port" \
        --from-literal=mysql-db="$mysql_db" \
        --from-literal=mysql-user="$mysql_user" \
        --from-literal=mysql-password="$mysql_password" \
        -n "$NAMESPACE" \
        --dry-run=client -o yaml | kubectl apply -f -
    
    info "Nacos Secret 已创建"
    echo "  MySQL Host: $mysql_host"
    echo "  MySQL Port: $mysql_port"
    echo "  MySQL DB: $mysql_db"
    echo "  MySQL User: $mysql_user"
    echo "  MySQL Password: $mysql_password"
    echo ""
}

# 函数：生成 RocketMQ Secret
generate_rocketmq_secret() {
    step "生成 RocketMQ Secret..."
    
    local username="${ROCKETMQ_USERNAME:-rocketmq}"
    local password="${ROCKETMQ_PASSWORD:-$(generate_password)}"
    
    if kubectl get secret rocketmq-credentials -n "$NAMESPACE" &> /dev/null; then
        warn "Secret rocketmq-credentials 已存在，跳过创建"
        return 0
    fi
    
    kubectl create secret generic rocketmq-credentials \
        --from-literal=username="$username" \
        --from-literal=password="$password" \
        -n "$NAMESPACE" \
        --dry-run=client -o yaml | kubectl apply -f -
    
    info "RocketMQ Secret 已创建"
    echo "  Username: $username"
    echo "  Password: $password"
    echo ""
}

# 函数：生成 Grafana Secret
generate_grafana_secret() {
    step "生成 Grafana Secret..."
    
    local admin_user="${GRAFANA_ADMIN_USER:-admin}"
    local admin_password="${GRAFANA_ADMIN_PASSWORD:-$(generate_password)}"
    
    if kubectl get secret grafana-admin -n "$NAMESPACE" &> /dev/null; then
        warn "Secret grafana-admin 已存在，跳过创建"
        return 0
    fi
    
    kubectl create secret generic grafana-admin \
        --from-literal=admin-user="$admin_user" \
        --from-literal=admin-password="$admin_password" \
        -n "$NAMESPACE" \
        --dry-run=client -o yaml | kubectl apply -f -
    
    info "Grafana Secret 已创建"
    echo "  Admin User: $admin_user"
    echo "  Admin Password: $admin_password"
    echo ""
}

# 函数：生成 Elasticsearch Secret
generate_elasticsearch_secret() {
    step "生成 Elasticsearch Secret..."
    
    local password="${ELASTICSEARCH_PASSWORD:-$(generate_password)}"
    
    if kubectl get secret elasticsearch-auth -n "$NAMESPACE" &> /dev/null; then
        warn "Secret elasticsearch-auth 已存在，跳过创建"
        return 0
    fi
    
    kubectl create secret generic elasticsearch-auth \
        --from-literal=password="$password" \
        -n "$NAMESPACE" \
        --dry-run=client -o yaml | kubectl apply -f -
    
    info "Elasticsearch Secret 已创建"
    echo "  Password: $password"
    echo ""
}

# 主函数
main() {
    echo "=========================================="
    echo "  Lushop Secret 生成脚本"
    echo "=========================================="
    echo ""
    warn "注意: 此脚本会生成随机密码"
    warn "如需使用自定义密码，请设置环境变量:"
    echo "  MYSQL_ROOT_PASSWORD=your_password"
    echo "  REDIS_PASSWORD=your_password"
    echo "  等等..."
    echo ""
    read -p "继续生成 Secret? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        info "已取消"
        exit 0
    fi
    echo ""
    
    create_namespace
    generate_mysql_secret
    generate_redis_secret
    generate_nacos_secret
    generate_rocketmq_secret
    generate_grafana_secret
    generate_elasticsearch_secret
    
    echo ""
    info "所有 Secret 已生成！"
    warn "请保存这些密码，后续可能需要使用"
    echo ""
    info "查看所有 Secret:"
    echo "  kubectl get secrets -n $NAMESPACE"
}

main "$@"

