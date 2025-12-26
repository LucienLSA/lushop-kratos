#!/bin/bash

# Lushop Image Update Script
# This script updates container images and restarts deployments

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
NAMESPACE="lushop"
REGISTRY="crpi-t2nj22cqo7df6n3a.cn-shanghai.personal.cr.aliyuncs.com/aliyun1123466419"
SERVICES=("lushop" "goods" "inventory" "order" "user" "userauth" "userop")

# Functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if kubectl is available
check_kubectl() {
    if ! command -v kubectl &> /dev/null; then
        log_error "kubectl is not installed or not in PATH"
        exit 1
    fi
}

# Get the latest commit SHA from GitHub
get_latest_sha() {
    local repo="${GITHUB_REPOSITORY:-lushop/lushop-kratos}"
    local branch="${GITHUB_HEAD_REF:-main}"

    # If running in GitHub Actions, use the provided SHA
    if [ -n "$GITHUB_SHA" ]; then
        echo "$GITHUB_SHA"
        return
    fi

    # Otherwise, try to get it from git
    if command -v git &> /dev/null && git rev-parse --short HEAD &> /dev/null; then
        git rev-parse --short HEAD
    else
        log_warning "Could not determine latest commit SHA, using 'latest' tag"
        echo "latest"
    fi
}

# Update kustomization.yaml with new image tags
update_kustomization() {
    local new_tag="$1"
    local kustomization_file="kustomization.yaml"

    log_info "Updating kustomization.yaml with new tag: $new_tag"

    # Update image tags in kustomization.yaml
    for service in "${SERVICES[@]}"; do
        sed -i.bak "s|${REGISTRY}/${service}:.*|${REGISTRY}/${service}:${new_tag}|g" "$kustomization_file"
    done

    log_success "Kustomization file updated"
}

# Deploy updated configuration
deploy_updates() {
    log_info "Deploying updated configuration..."

    if command -v kustomize &> /dev/null; then
        kustomize build . | kubectl apply -f -
    else
        kubectl apply -k .
    fi

    log_success "Configuration deployed"
}

# Wait for deployments to be ready
wait_for_rollout() {
    local timeout="${ROLLOUT_TIMEOUT:-600}"

    log_info "Waiting for deployments to be ready (timeout: ${timeout}s)..."

    for service in "${SERVICES[@]}"; do
        # Convert service name to deployment name
        local deployment_name="$service"
        if [ "$service" = "lushop" ]; then
            deployment_name="lushop-gateway"
        fi

        log_info "Waiting for $deployment_name..."
        if ! kubectl rollout status deployment/"$deployment_name" -n "$NAMESPACE" --timeout="${timeout}s"; then
            log_error "Deployment $deployment_name failed to rollout"
            return 1
        fi
    done

    log_success "All deployments are ready"
}

# Check deployment status
check_status() {
    log_info "Checking deployment status..."

    echo "Pods:"
    kubectl get pods -n "$NAMESPACE" -o wide

    echo ""
    echo "Deployments:"
    kubectl get deployments -n "$NAMESPACE"

    echo ""
    echo "Services:"
    kubectl get services -n "$NAMESPACE"

    echo ""
    echo "ReplicaSets:"
    kubectl get rs -n "$NAMESPACE"
}

# Rollback deployment
rollback_deployment() {
    local service="$1"
    local revision="${2:-1}"

    local deployment_name="$service"
    if [ "$service" = "lushop" ]; then
        deployment_name="lushop-gateway"
    fi

    log_info "Rolling back $deployment_name to revision $revision..."
    kubectl rollout undo deployment/"$deployment_name" --to-revision="$revision" -n "$NAMESPACE"
    log_success "Rollback completed"
}

# Show usage
show_usage() {
    echo "Usage: $0 [OPTIONS] [SERVICE]"
    echo ""
    echo "Update container images for Lushop services"
    echo ""
    echo "Options:"
    echo "  -t, --tag TAG     Use specific image tag (default: latest commit SHA)"
    echo "  -n, --namespace NS Kubernetes namespace (default: lushop)"
    echo "  -r, --registry REG Container registry URL"
    echo "  --rollback REV    Rollback to specific revision"
    echo "  --status          Show current status"
    echo "  -h, --help        Show this help"
    echo ""
    echo "Services:"
    echo "  all               Update all services (default)"
    echo "  lushop            Update only lushop gateway"
    echo "  goods             Update only goods service"
    echo "  inventory         Update only inventory service"
    echo "  order             Update only order service"
    echo "  user              Update only user service"
    echo "  userauth          Update only userauth service"
    echo "  userop            Update only userop service"
    echo ""
    echo "Examples:"
    echo "  $0                    # Update all services with latest commit"
    echo "  $0 -t v1.2.3         # Update all services with specific tag"
    echo "  $0 goods             # Update only goods service"
    echo "  $0 --rollback 1 goods # Rollback goods service"
    echo "  $0 --status          # Show current status"
}

# Main function
main() {
    local tag=""
    local rollback_revision=""
    local show_status=false
    local target_service="all"

    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -t|--tag)
                tag="$2"
                shift 2
                ;;
            -n|--namespace)
                NAMESPACE="$2"
                shift 2
                ;;
            -r|--registry)
                REGISTRY="$2"
                shift 2
                ;;
            --rollback)
                rollback_revision="$2"
                shift 2
                ;;
            --status)
                show_status=true
                shift
                ;;
            -h|--help)
                show_usage
                exit 0
                ;;
            -*)
                log_error "Unknown option: $1"
                show_usage
                exit 1
                ;;
            *)
                target_service="$1"
                shift
                ;;
        esac
    done

    check_kubectl

    # Show status if requested
    if [ "$show_status" = true ]; then
        check_status
        exit 0
    fi

    # Handle rollback
    if [ -n "$rollback_revision" ]; then
        if [ "$target_service" = "all" ]; then
            log_error "Cannot rollback all services at once. Specify a service name."
            exit 1
        fi
        rollback_deployment "$target_service" "$rollback_revision"
        exit 0
    fi

    # Determine image tag
    if [ -z "$tag" ]; then
        tag=$(get_latest_sha)
    fi

    log_info "Using image tag: $tag"
    log_info "Target service: $target_service"

    # Update kustomization file
    update_kustomization "$tag"

    # Deploy updates
    deploy_updates

    # Wait for rollout
    if ! wait_for_rollout; then
        log_error "Deployment failed. You may want to rollback."
        exit 1
    fi

    # Show final status
    check_status

    log_success "Image update completed successfully!"
}

# Run main function
main "$@"
