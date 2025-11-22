#!/usr/bin/env bash
set -euo pipefail

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 脚本目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
NAMESPACE="lushop"
TIMEOUT=300

# 函数：打印信息
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

success() {
    echo -e "${GREEN}✓${NC} $1"
}

fail() {
    echo -e "${RED}✗${NC} $1"
}

# 函数：检查前置条件
check_prerequisites() {
    step "检查前置条件..."
    
    # 检查 kubectl
    if ! command -v kubectl &> /dev/null; then
        error "kubectl 未安装"
        exit 1
    fi
    
    # 检查集群连接
    if ! kubectl cluster-info &> /dev/null; then
        error "无法连接到 Kubernetes 集群"
        exit 1
    fi
    success "kubectl 和集群连接正常"
    
    # 检查存储类
    if ! kubectl get storageclass &> /dev/null; then
        warn "无法获取存储类，请确保集群已配置 StorageClass"
    else
        local default_sc=$(kubectl get storageclass -o jsonpath='{.items[?(@.metadata.annotations.storageclass\.kubernetes\.io/is-default-class=="true")].metadata.name}')
        if [ -z "$default_sc" ]; then
            warn "没有默认存储类，PVC 可能需要手动指定 storageClassName"
        else
            success "默认存储类: $default_sc"
        fi
    fi
    
    # 检查镜像是否存在
    local images=("lushop/user:latest" "lushop/goods:latest" "lushop/order:latest" 
                  "lushop/inventory:latest" "lushop/userop:latest" "lushop/userauth:latest" 
                  "lushop/gateway:latest")
    local missing_images=()
    
    for image in "${images[@]}"; do
        if ! docker images "$image" | grep -q "$(echo $image | cut -d: -f1)"; then
            missing_images+=("$image")
        fi
    done
    
    if [ ${#missing_images[@]} -gt 0 ]; then
        warn "以下镜像未找到:"
        for img in "${missing_images[@]}"; do
            echo "  - $img"
        done
        echo ""
        read -p "是否继续部署? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            info "请先运行 ./build-images.sh all 构建镜像"
            exit 1
        fi
    else
        success "所有镜像已构建"
    fi
    
    echo ""
}

# 函数：等待 Pod 就绪
wait_for_pod() {
    local label=$1
    local name=$2
    step "等待 $name 就绪..."
    
    if kubectl wait --for=condition=ready pod -l "$label" -n "$NAMESPACE" --timeout=${TIMEOUT}s 2>/dev/null; then
        success "$name 已就绪"
        return 0
    else
        error "$name 启动超时或失败"
        kubectl get pods -l "$label" -n "$NAMESPACE"
        kubectl describe pods -l "$label" -n "$NAMESPACE" | tail -20
        return 1
    fi
}

# 函数：部署命名空间
deploy_namespace() {
    step "创建命名空间..."
    kubectl apply -f "${SCRIPT_DIR}/base/namespace.yaml"
    success "命名空间已创建"
    echo ""
}

# 函数：部署基础设施层
deploy_infrastructure() {
    step "部署基础设施层..."
    
    # Redis
    info "部署 Redis..."
    kubectl apply -k "${SCRIPT_DIR}/base/redis"
    wait_for_pod "app.kubernetes.io/name=redis" "Redis"
    echo ""
    
    # MySQL
    info "部署 MySQL..."
    kubectl apply -k "${SCRIPT_DIR}/base/mysql"
    wait_for_pod "app.kubernetes.io/name=mysql" "MySQL"
    echo ""
    
    # RocketMQ
    info "部署 RocketMQ..."
    kubectl apply -k "${SCRIPT_DIR}/base/rocketmq"
    wait_for_pod "app.kubernetes.io/name=rocketmq" "RocketMQ"
    echo ""
    
    success "基础设施层部署完成"
    echo ""
}

# 函数：初始化数据库
init_database() {
    step "初始化数据库..."
    
    # 检查是否有初始化脚本
    local init_script="${SCRIPT_DIR}/../scripts/init_db.sql"
    if [ ! -f "$init_script" ]; then
        warn "数据库初始化脚本不存在: $init_script"
        warn "请手动初始化数据库或创建初始化脚本"
        return 0
    fi
    
    # 获取 MySQL root 密码
    local mysql_password
    if ! mysql_password=$(kubectl get secret mysql-auth -n "$NAMESPACE" -o jsonpath='{.data.mysql-root-password}' 2>/dev/null | base64 -d); then
        error "无法获取 MySQL root 密码"
        return 1
    fi
    
    # 使用 port-forward 连接 MySQL
    info "建立 MySQL 连接..."
    kubectl port-forward -n "$NAMESPACE" svc/mysql 3306:3306 > /dev/null 2>&1 &
    local pf_pid=$!
    sleep 5
    
    # 检查 MySQL 是否可连接
    if ! mysql -h 127.0.0.1 -uroot -p"$mysql_password" -e "SELECT 1" &> /dev/null; then
        warn "无法连接到 MySQL，请手动初始化数据库"
        kill $pf_pid 2>/dev/null || true
        return 0
    fi
    
    # 导入数据库脚本
    info "导入数据库脚本..."
    if mysql -h 127.0.0.1 -uroot -p"$mysql_password" < "$init_script" 2>/dev/null; then
        success "数据库初始化完成"
    else
        warn "数据库初始化失败，请检查脚本和连接"
    fi
    
    # 停止 port-forward
    kill $pf_pid 2>/dev/null || true
    echo ""
}

# 函数：部署服务治理组件
deploy_governance() {
    step "部署服务治理组件..."
    
    # Nacos
    info "部署 Nacos..."
    kubectl apply -k "${SCRIPT_DIR}/base/nacos"
    wait_for_pod "app.kubernetes.io/name=nacos" "Nacos"
    echo ""
    
    # Consul
    info "部署 Consul..."
    kubectl apply -k "${SCRIPT_DIR}/base/consul"
    wait_for_pod "app.kubernetes.io/name=consul" "Consul"
    echo ""
    
    # 监控组件
    info "部署监控组件..."
    kubectl apply -k "${SCRIPT_DIR}/base/prometheus"
    kubectl apply -k "${SCRIPT_DIR}/base/grafana"
    echo ""
    
    # 链路追踪
    info "部署链路追踪..."
    kubectl apply -k "${SCRIPT_DIR}/base/jaeger"
    echo ""
    
    # 日志组件
    info "部署日志组件..."
    kubectl apply -k "${SCRIPT_DIR}/base/elasticsearch"
    wait_for_pod "app.kubernetes.io/name=elasticsearch" "Elasticsearch"
    kubectl apply -k "${SCRIPT_DIR}/base/kibana"
    echo ""
    
    success "服务治理组件部署完成"
    echo ""
}

# 函数：配置 Nacos
configure_nacos() {
    step "配置 Nacos..."
    
    info "Nacos 配置需要手动完成："
    echo ""
    echo "1. 使用 port-forward 访问 Nacos 控制台:"
    echo "   kubectl port-forward -n $NAMESPACE svc/nacos 8848:8848"
    echo ""
    echo "2. 访问: http://localhost:8848/nacos"
    echo "   默认账号: nacos / nacos"
    echo ""
    echo "3. 创建命名空间: de9c6a0e-1fbc-425d-8d3b-09066fea6889"
    echo ""
    echo "4. 为每个服务导入配置（参考各服务的 configs/nacos-config.yaml）"
    echo ""
    
    read -p "是否已配置 Nacos? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        warn "请先配置 Nacos，然后重新运行部署脚本"
        exit 1
    fi
    
    echo ""
}

# 函数：部署业务服务
deploy_services() {
    step "部署业务服务..."
    
    # 基础服务
    info "部署基础服务..."
    kubectl apply -k "${SCRIPT_DIR}/base/services/user"
    kubectl apply -k "${SCRIPT_DIR}/base/services/goods"
    kubectl apply -k "${SCRIPT_DIR}/base/services/inventory"
    kubectl apply -k "${SCRIPT_DIR}/base/services/userop"
    kubectl apply -k "${SCRIPT_DIR}/base/services/userauth"
    echo ""
    
    # 等待基础服务就绪
    info "等待基础服务就绪..."
    wait_for_pod "app.kubernetes.io/name=user" "User Service"
    wait_for_pod "app.kubernetes.io/name=goods" "Goods Service"
    wait_for_pod "app.kubernetes.io/name=inventory" "Inventory Service"
    echo ""
    
    # 依赖服务
    info "部署依赖服务..."
    kubectl apply -k "${SCRIPT_DIR}/base/services/order"
    echo ""
    
    # API 网关
    info "部署 API 网关..."
    kubectl apply -k "${SCRIPT_DIR}/base/services/gateway"
    echo ""
    
    success "业务服务部署完成"
    echo ""
}

# 函数：验证部署
verify_deployment() {
    step "验证部署..."
    
    echo ""
    info "Pod 状态:"
    kubectl get pods -n "$NAMESPACE"
    echo ""
    
    info "服务状态:"
    kubectl get svc -n "$NAMESPACE"
    echo ""
    
    # 检查失败的 Pod
    local failed_pods=$(kubectl get pods -n "$NAMESPACE" -o jsonpath='{.items[?(@.status.phase!="Running" && @.status.phase!="Succeeded")].metadata.name}')
    if [ -n "$failed_pods" ]; then
        warn "以下 Pod 状态异常:"
        echo "$failed_pods"
        echo ""
        info "查看详细信息:"
        echo "kubectl describe pod <pod-name> -n $NAMESPACE"
        echo "kubectl logs <pod-name> -n $NAMESPACE"
    else
        success "所有 Pod 运行正常"
    fi
    
    echo ""
}

# 函数：显示访问信息
show_access_info() {
    step "访问信息..."
    
    echo ""
    info "服务访问方式:"
    echo ""
    echo "1. Gateway (NodePort):"
    echo "   NODE_IP=\$(kubectl get nodes -o jsonpath='{.items[0].status.addresses[?(@.type==\"InternalIP\")].address}')"
    echo "   curl http://\$NODE_IP:30080/api/goods/list"
    echo ""
    echo "2. Port Forward:"
    echo "   kubectl port-forward -n $NAMESPACE svc/gateway-service 8001:8001"
    echo "   curl http://localhost:8001/health"
    echo ""
    echo "3. 监控和日志:"
    echo "   - Prometheus: kubectl port-forward -n $NAMESPACE svc/prometheus 9090:9090"
    echo "   - Grafana: kubectl port-forward -n $NAMESPACE svc/grafana 3000:3000"
    echo "   - Jaeger: kubectl port-forward -n $NAMESPACE svc/jaeger 16686:16686"
    echo "   - Kibana: kubectl port-forward -n $NAMESPACE svc/kibana 5601:5601"
    echo ""
}

# 主函数
main() {
    case "${1:-deploy}" in
        deploy)
            echo -e "${CYAN}========================================${NC}"
            echo -e "${CYAN}  Lushop K8s 部署脚本${NC}"
            echo -e "${CYAN}========================================${NC}"
            echo ""
            
            check_prerequisites
            deploy_namespace
            deploy_infrastructure
            init_database
            deploy_governance
            configure_nacos
            deploy_services
            verify_deployment
            show_access_info
            
            success "部署完成！"
            ;;
        status)
            kubectl get all -n "$NAMESPACE"
            ;;
        logs)
            if [ -z "${2:-}" ]; then
                error "请指定服务名称"
                echo "用法: $0 logs <service-name>"
                exit 1
            fi
            kubectl logs -f deployment/${2}-service -n "$NAMESPACE" || kubectl logs -f -l app.kubernetes.io/name=${2} -n "$NAMESPACE"
            ;;
        delete)
            warn "这将删除所有部署的资源"
            read -p "确认删除? (y/N): " -n 1 -r
            echo
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                kubectl delete -k "${SCRIPT_DIR}/base/"
                kubectl delete namespace "$NAMESPACE"
                success "已删除所有资源"
            else
                info "已取消"
            fi
            ;;
        *)
            echo "用法: $0 {deploy|status|logs|delete}"
            echo ""
            echo "命令:"
            echo "  deploy  - 部署所有服务（默认）"
            echo "  status  - 查看部署状态"
            echo "  logs    - 查看服务日志"
            echo "  delete  - 删除所有资源"
            exit 1
            ;;
    esac
}

main "$@"

