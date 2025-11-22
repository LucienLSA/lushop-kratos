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
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

# 镜像标签（可通过环境变量覆盖）
IMAGE_TAG="${IMAGE_TAG:-latest}"
IMAGE_REGISTRY="${IMAGE_REGISTRY:-lushop}"

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

# 函数：检查 iptables
check_iptables() {
    if [ ! -f /usr/sbin/iptables ] && [ ! -L /usr/sbin/iptables ]; then
        error "iptables 命令不存在，这会导致 Docker 网络功能失败"
        error ""
        error "修复方法（需要 root 权限）："
        error "  1. 创建 iptables 符号链接："
        error "     sudo update-alternatives --install /usr/sbin/iptables iptables /usr/sbin/iptables-legacy 10"
        error "     sudo update-alternatives --set iptables /usr/sbin/iptables-legacy"
        error ""
        error "  2. 或者重启 Docker 服务："
        error "     sudo systemctl restart docker"
        error ""
        error "  3. 验证修复："
        error "     /usr/sbin/iptables --version"
        error ""
        warn "继续构建可能会失败..."
        return 1
    fi
    
    if ! /usr/sbin/iptables --version &> /dev/null; then
        error "iptables 命令无法执行"
        return 1
    fi
    
    info "iptables 检查通过"
    return 0
}

# 函数：检查 Docker
check_docker() {
    if ! command -v docker &> /dev/null; then
        error "Docker 未安装，请先安装 Docker"
        exit 1
    fi
    
    if ! docker info &> /dev/null; then
        error "无法连接到 Docker daemon，请确保 Docker 正在运行"
        error "提示: sudo systemctl start docker 或重启 Docker Desktop"
        exit 1
    fi
    
    info "Docker 检查通过"
}

# 函数：构建服务镜像
build_service_image() {
    local service=$1
    local service_dir="${PROJECT_ROOT}/service/${service}"
    local dockerfile="${service_dir}/Dockerfile"
    local image_name="${IMAGE_REGISTRY}/${service}:${IMAGE_TAG}"
    
    if [ ! -d "$service_dir" ]; then
        error "服务目录不存在: $service_dir"
        return 1
    fi
    
    if [ ! -f "$dockerfile" ]; then
        warn "服务 $service 没有 Dockerfile，跳过"
        return 0
    fi
    
    # 特殊处理：order 服务需要访问 goods 目录，使用 service 目录作为构建上下文
    local build_context
    if [ "$service" = "order" ]; then
        build_context="${PROJECT_ROOT}/service"
        step "构建服务镜像: ${image_name} (使用特殊构建上下文)"
        info "  构建上下文: ${build_context} (order 服务需要访问 goods 目录)"
    else
        build_context="${service_dir}"
        step "构建服务镜像: ${image_name}"
        info "  构建上下文: ${build_context}"
    fi
    
    info "  Dockerfile: ${dockerfile}"
    
    cd "$build_context"
    
    if docker build \
        -f "$dockerfile" \
        -t "$image_name" \
        --build-arg BUILDKIT_INLINE_CACHE=1 \
        .; then
        success "服务 $service 镜像构建成功: ${image_name}"
        return 0
    else
        fail "服务 $service 镜像构建失败"
        return 1
    fi
}

# 函数：构建网关镜像
build_gateway_image() {
    local gateway_dir="${PROJECT_ROOT}/lushop"
    local dockerfile="${gateway_dir}/Dockerfile"
    local image_name="${IMAGE_REGISTRY}/gateway:${IMAGE_TAG}"
    
    if [ ! -d "$gateway_dir" ]; then
        error "网关目录不存在: $gateway_dir"
        return 1
    fi
    
    if [ ! -f "$dockerfile" ]; then
        error "网关没有 Dockerfile: $dockerfile"
        return 1
    fi
    
    step "构建网关镜像: ${image_name}"
    info "  构建上下文: ${gateway_dir}"
    info "  Dockerfile: ${dockerfile}"
    
    cd "$gateway_dir"
    
    if docker build \
        -f "$dockerfile" \
        -t "$image_name" \
        --build-arg BUILDKIT_INLINE_CACHE=1 \
        .; then
        success "网关镜像构建成功: ${image_name}"
        return 0
    else
        fail "网关镜像构建失败"
        return 1
    fi
}

# 函数：构建所有服务镜像
build_all_services() {
    info "开始构建所有服务镜像..."
    
    local services=("user" "goods" "order" "inventory" "userop" "userauth")
    local failed_services=()
    local success_count=0
    
    for service in "${services[@]}"; do
        if build_service_image "$service"; then
            ((success_count++))
        else
            failed_services+=("$service")
        fi
        echo ""
    done
    
    if [ ${#failed_services[@]} -gt 0 ]; then
        error "以下服务构建失败: ${failed_services[*]}"
        return 1
    fi
    
    success "所有服务镜像构建完成 (${success_count}/${#services[@]})"
    return 0
}

# 函数：构建所有镜像
build_all() {
    info "开始构建所有镜像..."
    info "镜像标签: ${IMAGE_TAG}"
    info "镜像仓库: ${IMAGE_REGISTRY}"
    echo ""
    
    check_docker
    check_iptables || warn "iptables 检查失败，构建可能会遇到网络问题"
    echo ""
    
    local total_failed=0
    
    # 构建服务镜像
    if ! build_all_services; then
        ((total_failed++))
    fi
    echo ""
    
    # 构建网关镜像
    if ! build_gateway_image; then
        ((total_failed++))
    fi
    
    echo ""
    if [ $total_failed -eq 0 ]; then
        success "所有镜像构建完成！"
        info "使用 'docker images | grep ${IMAGE_REGISTRY}' 查看构建的镜像"
    else
        error "部分镜像构建失败，请检查错误信息"
        return 1
    fi
}

# 函数：仅构建服务镜像
build_services_only() {
    info "开始构建服务镜像..."
    info "镜像标签: ${IMAGE_TAG}"
    echo ""
    
    check_docker
    check_iptables || warn "iptables 检查失败，构建可能会遇到网络问题"
    echo ""
    build_all_services
    
    echo ""
    info "服务镜像构建完成！"
}

# 函数：仅构建网关镜像
build_gateway_only() {
    info "开始构建网关镜像..."
    info "镜像标签: ${IMAGE_TAG}"
    echo ""
    
    check_docker
    check_iptables || warn "iptables 检查失败，构建可能会遇到网络问题"
    echo ""
    build_gateway_image
    
    echo ""
    info "网关镜像构建完成！"
}

# 函数：构建指定服务
build_specific() {
    local service=$1
    
    check_docker
    check_iptables || warn "iptables 检查失败，构建可能会遇到网络问题"
    echo ""
    
    if [ "$service" = "gateway" ]; then
        build_gateway_image
    else
        build_service_image "$service"
    fi
}

# 函数：列出所有镜像
list_images() {
    info "Lushop 相关镜像:"
    echo ""
    docker images | grep -E "${IMAGE_REGISTRY}|REPOSITORY" || warn "未找到 ${IMAGE_REGISTRY} 镜像"
    echo ""
    
    # 显示镜像大小统计
    local total_size=$(docker images --format "{{.Size}}" | grep -E "^[0-9.]+[KMGT]i?B$" | awk '{sum+=$1} END {print sum}' 2>/dev/null || echo "0")
    if [ -n "$total_size" ] && [ "$total_size" != "0" ]; then
        info "总镜像大小: ${total_size}"
    fi
}

# 函数：清理未使用的镜像
clean_images() {
    warn "这将删除所有未使用的 ${IMAGE_REGISTRY} 镜像"
    read -p "确认删除? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        info "清理未使用的镜像..."
        docker images "${IMAGE_REGISTRY}/*" --format "{{.ID}}" | xargs -r docker rmi -f 2>/dev/null || true
        success "清理完成"
    else
        info "已取消"
    fi
}

# 函数：推送镜像到仓库
push_images() {
    local registry=${1:-""}
    
    if [ -z "$registry" ]; then
        error "请指定镜像仓库地址"
        echo ""
        echo "用法: $0 push <registry>"
        echo "示例: $0 push docker.io/your-username"
        echo "      $0 push registry.example.com:5000"
        exit 1
    fi
    
    info "推送镜像到 $registry ..."
    echo ""
    
    local services=("user" "goods" "order" "inventory" "userop" "userauth" "gateway")
    local failed=0
    
    for service in "${services[@]}"; do
        local image="${IMAGE_REGISTRY}/${service}:${IMAGE_TAG}"
        local target="${registry}/${IMAGE_REGISTRY}/${service}:${IMAGE_TAG}"
        
        if docker images "$image" | grep -q "$service"; then
            step "标记并推送 $image -> $target"
            docker tag "$image" "$target" || { fail "标记失败"; ((failed++)); continue; }
            docker push "$target" || { fail "推送失败"; ((failed++)); continue; }
            success "推送成功: $target"
        else
            warn "镜像 $image 不存在，跳过"
        fi
        echo ""
    done
    
    if [ $failed -eq 0 ]; then
        success "所有镜像推送完成"
    else
        error "$failed 个镜像推送失败"
        return 1
    fi
}

# 函数：显示帮助信息
show_help() {
    echo -e "${CYAN}Lushop 镜像构建脚本${NC}"
    echo ""
    echo "用法: $0 {命令} [选项]"
    echo ""
    echo -e "${YELLOW}命令:${NC}"
    echo "  all                 - 构建所有镜像（默认）"
    echo "  services, svc       - 仅构建服务镜像"
    echo "  gateway, gw         - 仅构建网关镜像"
    echo "  list, ls            - 列出所有 lushop 镜像"
    echo "  clean               - 清理未使用的镜像"
    echo "  push <registry>     - 推送镜像到仓库"
    echo "  <service-name>      - 构建指定服务"
    echo ""
    echo -e "${YELLOW}服务名称:${NC}"
    echo "  user, goods, order, inventory, userop, userauth, gateway"
    echo ""
    echo -e "${YELLOW}环境变量:${NC}"
    echo "  IMAGE_TAG           - 镜像标签（默认: latest）"
    echo "  IMAGE_REGISTRY      - 镜像仓库前缀（默认: lushop）"
    echo ""
    echo -e "${YELLOW}示例:${NC}"
    echo "  $0                          # 构建所有镜像"
    echo "  $0 user                     # 仅构建 user 服务"
    echo "  $0 gateway                  # 仅构建网关"
    echo "  $0 services                 # 构建所有服务（不含网关）"
    echo "  $0 list                     # 列出所有镜像"
    echo "  IMAGE_TAG=v1.0.0 $0 all     # 使用指定标签构建"
    echo "  $0 push docker.io/username  # 推送所有镜像到 Docker Hub"
    echo ""
}

# 主函数
main() {
    case "${1:-all}" in
        all|"")
            build_all
            ;;
        services|svc)
            build_services_only
            ;;
        gateway|gw)
            build_gateway_only
            ;;
        list|ls)
            list_images
            ;;
        clean)
            clean_images
            ;;
        push)
            push_images "${2:-}"
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            # 尝试构建指定服务
            if [ -n "$1" ]; then
                build_specific "$1"
            else
                show_help
                exit 1
            fi
            ;;
    esac
}

main "$@"

