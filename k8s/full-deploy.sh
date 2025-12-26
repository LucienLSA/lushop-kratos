#!/bin/bash

# Lushop 完整 K8s 部署脚本
# 部署完整的微服务架构：基础设施 + 业务服务

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 配置
NAMESPACE="lushop"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# 日志函数
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

# 检查依赖
check_dependencies() {
    log_info "检查依赖..."

    if ! command -v kubectl &> /dev/null; then
        log_error "kubectl 未安装"
        exit 1
    fi

    if ! command -v docker &> /dev/null; then
        log_error "docker 未安装"
        exit 1
    fi

    if ! kubectl cluster-info &> /dev/null; then
        log_error "无法连接到 K8s 集群"
        exit 1
    fi

    log_success "依赖检查通过"
}

# 环境准备
setup_environment() {
    log_info "设置环境..."

    # 创建命名空间
    kubectl create namespace "$NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -

    # 检查 StorageClass
    if ! kubectl get storageclass | grep -q "local-path"; then
        log_warning "未找到 local-path StorageClass，正在安装..."
        kubectl apply -f https://raw.githubusercontent.com/rancher/local-path-provisioner/v0.0.24/deploy/local-path-storage.yaml
        kubectl wait --for=condition=available --timeout=300s deployment/local-path-provisioner -n local-path-storage
    fi

    # 检查 NGINX Ingress
    if ! kubectl get deployment ingress-nginx-controller -n ingress-nginx &>/dev/null; then
        log_info "安装 NGINX Ingress Controller..."
        kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.1/deploy/static/provider/cloud/deploy.yaml
        kubectl wait --for=condition=available --timeout=300s deployment/ingress-nginx-controller -n ingress-nginx
    fi

    log_success "环境设置完成"
}

# 构建镜像
build_images() {
    log_info "构建服务镜像..."

    cd "$SCRIPT_DIR"
    if [ -f "build-images.sh" ]; then
        ./build-images.sh all
    else
        log_error "build-images.sh 不存在"
        exit 1
    fi

    log_success "镜像构建完成"
}

# 生成 Secrets
generate_secrets() {
    log_info "生成 K8s Secrets..."

    cd "$SCRIPT_DIR"
    if [ -f "gen-secrets-custom.sh" ]; then
        ./gen-secrets-custom.sh
    else
        log_error "gen-secrets-custom.sh 不存在"
        exit 1
    fi

    log_success "Secrets 生成完成"
}

# 部署基础设施
deploy_infrastructure() {
    log_info "部署基础设施服务..."

    local components=("mysql" "redis" "nacos" "consul" "rocketmq" "elasticsearch" "grafana" "jaeger" "kibana" "prometheus")

    for component in "${components[@]}"; do
        log_info "部署 $component..."
        if kubectl apply -k "base/$component" 2>/dev/null; then
            log_success "$component 部署成功"
        else
            log_warning "$component 部署失败或已存在"
        fi
    done

    # 等待基础设施服务就绪
    log_info "等待基础设施服务就绪..."
    kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=mysql -n "$NAMESPACE" --timeout=300s || log_warning "MySQL 未就绪"
    kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=redis -n "$NAMESPACE" --timeout=300s || log_warning "Redis 未就绪"
    kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=nacos -n "$NAMESPACE" --timeout=300s || log_warning "Nacos 未就绪"

    log_success "基础设施部署完成"
}

# 配置 Nacos
configure_nacos() {
    log_info "配置 Nacos..."

    cd "$SCRIPT_DIR"
    if [ -f "configure-nacos.sh" ]; then
        ./configure-nacos.sh import
    else
        log_warning "configure-nacos.sh 不存在，跳过 Nacos 配置"
    fi

    log_success "Nacos 配置完成"
}

# 部署业务服务
deploy_services() {
    log_info "部署业务服务..."

    # 导入镜像到 containerd
    log_info "导入服务镜像..."
    docker save \
        lushop/gateway:latest \
        lushop/user:latest \
        lushop/goods:latest \
        lushop/order:latest \
        lushop/inventory:latest \
        lushop/userauth:latest \
        lushop/userop:latest \
        -o /tmp/lushop-services.tar
    sudo ctr -n k8s.io images import /tmp/lushop-services.tar

    # 部署业务服务
    if kubectl apply -k . 2>/dev/null; then
        log_success "业务服务部署成功"
    else
        log_error "业务服务部署失败"
        exit 1
    fi

    # 等待服务就绪
    log_info "等待业务服务就绪..."
    local services=("lushop-gateway" "user-service" "goods-service" "order-service" "inventory-service" "userauth-service" "userop-service")

    for service in "${services[@]}"; do
        kubectl wait --for=condition=available --timeout=300s deployment/"$service" -n "$NAMESPACE" || log_warning "$service 未就绪"
    done

    log_success "业务服务部署完成"
}

# 设置 Ingress
setup_ingress() {
    log_info "设置 Ingress..."

    # 应用 Ingress 配置
    if kubectl apply -f ingress.yaml 2>/dev/null; then
        log_success "Ingress 设置完成"
    else
        log_warning "Ingress 配置不存在或已存在"
    fi
}

# 验证部署
verify_deployment() {
    log_info "验证部署..."

    echo ""
    echo "=== Pods 状态 ==="
    kubectl get pods -n "$NAMESPACE" -o wide

    echo ""
    echo "=== Services 状态 ==="
    kubectl get services -n "$NAMESPACE"

    echo ""
    echo "=== Ingress 状态 ==="
    kubectl get ingress -n "$NAMESPACE"

    echo ""
    echo "=== 访问地址 ==="
    echo "API Gateway: http://$(kubectl get svc lushop-gateway -n "$NAMESPACE" -o jsonpath='{.spec.clusterIP}'):8001"
    echo "Nacos: http://$(kubectl get svc nacos -n "$NAMESPACE" -o jsonpath='{.spec.clusterIP}'):8848"
    echo "Grafana: http://$(kubectl get svc grafana -n "$NAMESPACE" -o jsonpath='{.spec.clusterIP}'):3000"

    log_success "部署验证完成"
}

# 显示帮助
show_help() {
    echo "Lushop K8s 完整部署脚本"
    echo ""
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  full        完整部署 (默认)"
    echo "  infra       只部署基础设施"
    echo "  services    只部署业务服务"
    echo "  clean       清理所有资源"
    echo "  status      查看部署状态"
    echo "  help        显示此帮助"
    echo ""
    echo "环境变量:"
    echo "  MYSQL_ROOT_PASSWORD    MySQL root 密码 (默认: MyRoot@123)"
    echo "  MYSQL_PASSWORD         MySQL 用户密码 (默认: lushopDb@123)"
    echo "  REDIS_PASSWORD         Redis 密码 (默认: Redis@123)"
    echo "  NACOS_MYSQL_PASSWORD   Nacos MySQL 密码 (默认: Nacos@123)"
}

# 清理函数
cleanup() {
    log_info "清理资源..."

    # 删除命名空间 (会删除所有资源)
    kubectl delete namespace "$NAMESPACE" --ignore-not-found=true

    # 删除 Ingress
    kubectl delete namespace ingress-nginx --ignore-not-found=true

    log_success "清理完成"
}

# 主函数
main() {
    local action="${1:-full}"

    case "$action" in
        "full")
            check_dependencies
            setup_environment
            build_images
            generate_secrets
            deploy_infrastructure
            configure_nacos
            deploy_services
            setup_ingress
            verify_deployment
            ;;
        "infra")
            check_dependencies
            setup_environment
            generate_secrets
            deploy_infrastructure
            ;;
        "services")
            check_dependencies
            deploy_services
            setup_ingress
            verify_deployment
            ;;
        "clean")
            cleanup
            ;;
        "status")
            verify_deployment
            ;;
        "help"|*)
            show_help
            ;;
    esac
}

# 执行主函数
main "$@"
