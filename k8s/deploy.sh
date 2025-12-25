#!/bin/bash

# Lushop Kubernetes Deployment Script
# This script deploys the Lushop microservices to Kubernetes

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
NAMESPACE="lushop"
KUBECONFIG_PATH="${KUBECONFIG:-$HOME/.kube/config}"
REGISTRY="crpi-t2nj22cqo7df6n3a.cn-shanghai.personal.cr.aliyuncs.com/aliyun1123466419"

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

check_dependencies() {
    log_info "Checking dependencies..."

    if ! command -v kubectl &> /dev/null; then
        log_error "kubectl is not installed. Please install it first."
        exit 1
    fi

    if ! command -v kustomize &> /dev/null; then
        log_warning "kustomize is not installed. Using kubectl kustomize instead."
    fi

    log_success "Dependencies check passed"
}

check_kube_connection() {
    log_info "Checking Kubernetes connection..."

    if ! kubectl cluster-info &> /dev/null; then
        log_error "Cannot connect to Kubernetes cluster. Please check your kubeconfig."
        exit 1
    fi

    log_success "Kubernetes connection established"
}

create_namespace() {
    log_info "Creating namespace: $NAMESPACE"

    if ! kubectl get namespace "$NAMESPACE" &> /dev/null; then
        kubectl apply -f namespace.yaml
        log_success "Namespace created"
    else
        log_warning "Namespace already exists"
    fi
}

deploy_infrastructure() {
    log_info "Deploying infrastructure components..."

    # This assumes you have infrastructure deployment scripts in the deploy directory
    # You may need to run these separately or integrate them here

    log_info "Please ensure the following infrastructure is running:"
    log_info "  - MySQL database"
    log_info "  - Redis cache"
    log_info "  - Nacos service registry"
    log_info "  - Consul (if used)"
    log_info "  - Elasticsearch (if used)"
    log_info "  - RocketMQ (if used)"

    # Example infrastructure deployment (uncomment and modify as needed)
    # ./deploy/mysql/mysql.sh
    # ./deploy/redis/redis.sh
    # ./deploy/nacos/nacos.sh
}

deploy_services() {
    log_info "Deploying Lushop services..."

    # Use kustomize to deploy all resources
    if command -v kustomize &> /dev/null; then
        kustomize build . | kubectl apply -f -
    else
        kubectl apply -k .
    fi

    log_success "Services deployed successfully"
}

wait_for_deployments() {
    log_info "Waiting for deployments to be ready..."

    local services=("lushop-gateway" "goods-service" "inventory-service" "order-service" "user-service" "userauth-service" "userop-service")

    for service in "${services[@]}"; do
        log_info "Waiting for $service..."
        kubectl wait --for=condition=available --timeout=300s deployment/"$service" -n "$NAMESPACE"
    done

    log_success "All deployments are ready"
}

check_services() {
    log_info "Checking service status..."

    kubectl get pods -n "$NAMESPACE"
    kubectl get services -n "$NAMESPACE"
    kubectl get ingress -n "$NAMESPACE"
}

update_images() {
    log_info "Updating images..."

    # Pull latest images
    local services=("lushop" "goods" "inventory" "order" "user" "userauth" "userop")

    for service in "${services[@]}"; do
        local image="${REGISTRY}/${service}:latest"
        log_info "Pulling image: $image"
        docker pull "$image" || log_warning "Failed to pull $image"
    done

    # Restart deployments to use new images
    for service in "${services[@]}"; do
        log_info "Restarting deployment: $service"
        kubectl rollout restart deployment/"$service" -n "$NAMESPACE"
    done

    log_success "Images updated and deployments restarted"
}

rollback() {
    local service="$1"
    local revision="$2"

    if [ -z "$service" ] || [ -z "$revision" ]; then
        log_error "Usage: rollback <service> <revision>"
        exit 1
    fi

    log_info "Rolling back $service to revision $revision..."
    kubectl rollout undo deployment/"$service" --to-revision="$revision" -n "$NAMESPACE"
    log_success "Rollback completed"
}

scale_service() {
    local service="$1"
    local replicas="$2"

    if [ -z "$service" ] || [ -z "$replicas" ]; then
        log_error "Usage: scale_service <service> <replicas>"
        exit 1
    fi

    log_info "Scaling $service to $replicas replicas..."
    kubectl scale deployment "$service" --replicas="$replicas" -n "$NAMESPACE"
    log_success "Service scaled successfully"
}

show_usage() {
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  deploy      Deploy all services to Kubernetes"
    echo "  update      Update images and restart deployments"
    echo "  status      Show status of all services"
    echo "  rollback    Rollback a service to previous revision"
    echo "  scale       Scale a service"
    echo "  cleanup     Remove all resources"
    echo ""
    echo "Examples:"
    echo "  $0 deploy"
    echo "  $0 update"
    echo "  $0 rollback user-service 1"
    echo "  $0 scale goods-service 3"
}

cleanup() {
    log_warning "This will delete all Lushop resources in namespace $NAMESPACE"
    read -p "Are you sure? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        log_info "Cleaning up resources..."
        kubectl delete namespace "$NAMESPACE" --ignore-not-found=true
        log_success "Cleanup completed"
    fi
}

# Main script
main() {
    local command="$1"

    case "$command" in
        deploy)
            check_dependencies
            check_kube_connection
            create_namespace
            deploy_infrastructure
            deploy_services
            wait_for_deployments
            check_services
            ;;
        update)
            check_dependencies
            check_kube_connection
            update_images
            wait_for_deployments
            ;;
        status)
            check_dependencies
            check_kube_connection
            check_services
            ;;
        rollback)
            check_dependencies
            check_kube_connection
            rollback "$2" "$3"
            ;;
        scale)
            check_dependencies
            check_kube_connection
            scale_service "$2" "$3"
            ;;
        cleanup)
            check_dependencies
            check_kube_connection
            cleanup
            ;;
        *)
            show_usage
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"
