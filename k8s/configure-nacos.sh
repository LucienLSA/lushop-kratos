#!/bin/bash

# Lushop Nacos 配置自动导入脚本
# 用于在 K8s 部署时自动将服务配置导入到 Nacos

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 配置
NAMESPACE="lushop"
NACOS_URL="http://localhost:8848"
USERNAME="nacos"
PASSWORD="nacos"
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

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

# 等待 Nacos 启动
wait_for_nacos() {
    log_info "等待 Nacos 启动..."

    local timeout=120
    while [ $timeout -gt 0 ]; do
        if curl -s "$NACOS_URL/nacos/v1/console/health/readiness" > /dev/null 2>&1; then
            log_success "Nacos 已启动"
            return 0
        fi
        sleep 2
        timeout=$((timeout - 2))
        log_info "等待 Nacos... 剩余 ${timeout} 秒"
    done

    log_error "Nacos 启动超时"
    exit 1
}

# 创建命名空间
create_namespace() {
    log_info "创建 Nacos 命名空间: lushop"

    curl -X POST "$NACOS_URL/nacos/v1/console/namespaces" \
        -u "$USERNAME:$PASSWORD" \
        -d "customNamespaceId=de9c6a0e-1fbc-425d-8d3b-09066fea6889" \
        -d "namespaceName=lushop" \
        -d "namespaceDesc=Lushop microservices namespace" \
        || log_warning "命名空间可能已存在"
}

# 导入配置到 Nacos
import_config() {
    local data_id=$1
    local config_file=$2
    local desc=$3

    if [ ! -f "$config_file" ]; then
        log_error "配置文件不存在: $config_file"
        return 1
    fi

    log_info "导入配置: $data_id"

    # 获取密码（如果配置中有占位符，需要替换）
    local mysql_pass=$(kubectl get secret mysql-auth -n "$NAMESPACE" -o jsonpath='{.data.mysql-password}' 2>/dev/null | base64 -d 2>/dev/null || echo "lushopDb@123")
    local redis_pass=$(kubectl get secret redis-auth -n "$NAMESPACE" -o jsonpath='{.data.redis-password}' 2>/dev/null | base64 -d 2>/dev/null || echo "Redis@123")

    # 读取配置文件内容并替换占位符
    local config_content
    config_content=$(sed "s/YourDBPasswordHere/$mysql_pass/g; s/YourRedisPasswordHere/$redis_pass/g" "$config_file")

    # 发布配置
    local response
    response=$(curl -s -X POST "$NACOS_URL/nacos/v1/cs/configs" \
        -u "$USERNAME:$PASSWORD" \
        -d "dataId=$data_id" \
        -d "group=lushop_grpc" \
        -d "content=$config_content" \
        -d "type=YAML" \
        -d "tenant=de9c6a0e-1fbc-425d-8d3b-09066fea6889" \
        -d "desc=$desc")

    if echo "$response" | grep -q "true"; then
        log_success "配置 $data_id 导入成功"
    else
        log_error "配置 $data_id 导入失败: $response"
        return 1
    fi
}

# 验证配置
verify_configs() {
    log_info "验证配置导入..."

    local configs=("user.yaml" "goods.yaml" "order.yaml" "inventory.yaml" "userop.yaml" "userauth.yaml" "gateway.yaml")
    local failed=0

    for config in "${configs[@]}"; do
        local response
        response=$(curl -s "$NACOS_URL/nacos/v1/cs/configs" \
            -u "$USERNAME:$PASSWORD" \
            -d "dataId=$config" \
            -d "group=lushop_grpc" \
            -d "tenant=de9c6a0e-1fbc-425d-8d3b-09066fea6889")

        if echo "$response" | grep -q "config"; then
            log_success "配置 $config 验证通过"
        else
            log_error "配置 $config 验证失败"
            ((failed++))
        fi
    done

    if [ $failed -gt 0 ]; then
        log_error "有 $failed 个配置验证失败"
        return 1
    else
        log_success "所有配置验证通过"
    fi
}

# 主函数
main() {
    echo "=========================================="
    echo "  Lushop Nacos 配置导入脚本"
    echo "=========================================="
    echo ""

    case "${1:-import}" in
        import)
            wait_for_nacos
            create_namespace

            log_info "开始导入服务配置..."

            # 导入各服务配置
            import_config "user.yaml" "$PROJECT_ROOT/service/user/configs/nacos-config-k8s.yaml" "User service configuration"
            import_config "goods.yaml" "$PROJECT_ROOT/service/goods/configs/nacos-config-k8s.yaml" "Goods service configuration"
            import_config "order.yaml" "$PROJECT_ROOT/service/order/configs/nacos-config-k8s.yaml" "Order service configuration"
            import_config "inventory.yaml" "$PROJECT_ROOT/service/inventory/configs/nacos-config-k8s.yaml" "Inventory service configuration"
            import_config "userop.yaml" "$PROJECT_ROOT/service/userop/configs/nacos-config-k8s.yaml" "UserOp service configuration"
            import_config "userauth.yaml" "$PROJECT_ROOT/service/userauth/configs/nacosRemote-k8s.yaml" "UserAuth service configuration"
            import_config "gateway.yaml" "$PROJECT_ROOT/lushop/configs/nacos-config-k8s.yaml" "Gateway service configuration"

            verify_configs
            ;;
        verify)
            verify_configs
            ;;
        *)
            echo "用法: $0 {import|verify}"
            echo ""
            echo "命令:"
            echo "  import  - 导入所有配置到 Nacos"
            echo "  verify  - 验证配置是否正确导入"
            exit 1
            ;;
    esac

    echo ""
    log_success "配置导入完成！"
}

main "$@"
