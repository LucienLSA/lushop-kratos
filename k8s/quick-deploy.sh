#!/bin/bash

# Lushop Quick Deployment Script
# This script provides a quick way to deploy Lushop to Kubernetes

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
NAMESPACE="lushop"
TIMEOUT="600"

# Functions
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

# Check prerequisites
check_prerequisites() {
    log_info "Checking prerequisites..."

    # Check if kubectl is installed
    if ! command -v kubectl &> /dev/null; then
        log_error "kubectl is not installed. Please install it first."
        log_info "Visit: https://kubernetes.io/docs/tasks/tools/"
        exit 1
    fi

    # Check if connected to Kubernetes cluster
    if ! kubectl cluster-info &> /dev/null; then
        log_error "Cannot connect to Kubernetes cluster."
        log_info "Please check your kubeconfig and cluster connectivity."
        exit 1
    fi

    # Check if kustomize is available
    if ! command -v kustomize &> /dev/null && ! kubectl kustomize --help &> /dev/null; then
        log_warning "kustomize not found. Will use kubectl apply -k instead."
    fi

    log_success "Prerequisites check passed"
}

# Setup prerequisites
setup_prerequisites() {
    log_info "Setting up prerequisites..."

    # Create namespace if it doesn't exist
    if ! kubectl get namespace "$NAMESPACE" &> /dev/null; then
        log_info "Creating namespace: $NAMESPACE"
        kubectl create namespace "$NAMESPACE"
    else
        log_warning "Namespace $NAMESPACE already exists"
    fi

    # Check if NGINX Ingress is installed
    if ! kubectl get deployment ingress-nginx-controller -n ingress-nginx &> /dev/null 2>&1; then
        log_warning "NGINX Ingress Controller not found."
        log_info "Installing NGINX Ingress Controller..."

        # Install NGINX Ingress Controller
        kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.1/deploy/static/provider/cloud/deploy.yaml

        # Wait for ingress controller to be ready
        log_info "Waiting for NGINX Ingress Controller to be ready..."
        kubectl wait --for=condition=available --timeout=300s deployment/ingress-nginx-controller -n ingress-nginx
    fi

    log_success "Prerequisites setup completed"
}

# Deploy infrastructure
deploy_infrastructure() {
    local skip_infra="${SKIP_INFRASTRUCTURE:-false}"

    if [ "$skip_infra" = "true" ]; then
        log_info "Skipping infrastructure deployment (--skip-infrastructure)"
        return
    fi

    log_info "Deploying infrastructure services..."

    # Deploy infrastructure
    kubectl apply -f infrastructure.yaml

    # Wait for infrastructure services
    log_info "Waiting for infrastructure services to be ready..."

    # Wait for MySQL
    kubectl wait --for=condition=available --timeout=300s deployment/mysql -n "$NAMESPACE" || {
        log_warning "MySQL deployment not ready, continuing..."
    }

    # Wait for Redis
    kubectl wait --for=condition=available --timeout=300s deployment/redis -n "$NAMESPACE" || {
        log_warning "Redis deployment not ready, continuing..."
    }

    # Wait for Nacos
    kubectl wait --for=condition=available --timeout=300s deployment/nacos -n "$NAMESPACE" || {
        log_warning "Nacos deployment not ready, continuing..."
    }

    log_success "Infrastructure deployment completed"
}

# Configure Nacos
configure_nacos() {
    log_info "Configuring Nacos with service configurations..."

    # Port forward Nacos
    kubectl port-forward -n "$NAMESPACE" svc/nacos 8848:8848 &
    NACOS_PF_PID=$!

    # Wait for port forward
    sleep 3

    # Import configurations using the configure script
    if ! ./configure-nacos.sh import; then
        log_error "Failed to configure Nacos"
        kill $NACOS_PF_PID 2>/dev/null || true
        return 1
    fi

    # Stop port forward
    kill $NACOS_PF_PID 2>/dev/null || true

    log_success "Nacos configuration completed"
}

# Deploy application services
deploy_application() {
    log_info "Deploying Lushop application services..."

    # Apply kustomization
    if command -v kustomize &> /dev/null; then
        kustomize build . | kubectl apply -f -
    else
        kubectl apply -k .
    fi

    log_success "Application services deployed"
}

# Wait for all services to be ready
wait_for_services() {
    log_info "Waiting for all services to be ready (timeout: ${TIMEOUT}s)..."

    local services=("lushop-gateway" "goods-service" "inventory-service" "order-service" "user-service" "userauth-service" "userop-service")

    for service in "${services[@]}"; do
        log_info "Waiting for $service..."
        if ! kubectl wait --for=condition=available --timeout="${TIMEOUT}s" deployment/"$service" -n "$NAMESPACE"; then
            log_error "Service $service failed to become ready"
            return 1
        fi
    done

    log_success "All services are ready"
}

# Show deployment status
show_status() {
    log_info "Deployment status:"

    echo ""
    echo "=== Pods ==="
    kubectl get pods -n "$NAMESPACE" -o wide

    echo ""
    echo "=== Services ==="
    kubectl get services -n "$NAMESPACE"

    echo ""
    echo "=== Deployments ==="
    kubectl get deployments -n "$NAMESPACE"

    echo ""
    echo "=== Ingress ==="
    kubectl get ingress -n "$NAMESPACE"

    echo ""
    echo "=== ConfigMaps ==="
    kubectl get configmap -n "$NAMESPACE"

    echo ""
    echo "=== Secrets ==="
    kubectl get secrets -n "$NAMESPACE"
}

# Get service endpoints
get_endpoints() {
    log_info "Service endpoints:"

    # Get ingress hosts
    local ingress_hosts=$(kubectl get ingress -n "$NAMESPACE" -o jsonpath='{.items[*].spec.rules[*].host}')

    if [ -n "$ingress_hosts" ]; then
        echo ""
        echo "=== External Access ==="
        for host in $ingress_hosts; do
            echo "https://$host"
        done
    fi

    # Get cluster internal services
    echo ""
    echo "=== Internal Services ==="
    kubectl get services -n "$NAMESPACE" -o custom-columns="NAME:.metadata.name,TYPE:.spec.type,CLUSTER-IP:.spec.clusterIP,PORTS:.spec.ports[*].port"
}

# Test services
test_services() {
    local ingress_host=$(kubectl get ingress -n "$NAMESPACE" -o jsonpath='{.items[0].spec.rules[0].host}' 2>/dev/null)

    if [ -n "$ingress_host" ]; then
        log_info "Testing service endpoints..."

        # Test health endpoint
        if curl -k -s "https://$ingress_host/health" > /dev/null; then
            log_success "Health check passed: https://$ingress_host/health"
        else
            log_warning "Health check failed. Service may still be starting up."
        fi

        # Test API endpoint
        if curl -k -s "https://$ingress_host/api/v1" > /dev/null; then
            log_success "API endpoint accessible: https://$ingress_host/api/v1"
        else
            log_warning "API endpoint not accessible yet."
        fi
    else
        log_warning "No ingress found. Skipping endpoint tests."
    fi
}

# Cleanup function
cleanup() {
    log_warning "Deployment failed. Cleaning up..."
    kubectl delete namespace "$NAMESPACE" --ignore-not-found=true
}

# Show usage
show_usage() {
    cat << EOF
Lushop Kubernetes Quick Deployment Script

USAGE:
    $0 [OPTIONS]

OPTIONS:
    -h, --help                    Show this help message
    -n, --namespace NAMESPACE     Kubernetes namespace (default: lushop)
    -t, --timeout SECONDS         Timeout for waiting operations (default: 600)
    --skip-infrastructure         Skip infrastructure deployment
    --skip-wait                   Skip waiting for services to be ready
    --skip-test                   Skip service testing
    --cleanup-on-failure          Cleanup resources on deployment failure

EXAMPLES:
    $0                            # Deploy everything
    $0 --skip-infrastructure      # Skip infrastructure, deploy only app services
    $0 --namespace production     # Deploy to production namespace
    $0 --cleanup-on-failure       # Cleanup on failure

ENVIRONMENT VARIABLES:
    SKIP_INFRASTRUCTURE           Set to 'true' to skip infrastructure
    SKIP_WAIT                     Set to 'true' to skip waiting
    SKIP_TEST                     Set to 'true' to skip testing

EOF
}

# Parse arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_usage
                exit 0
                ;;
            -n|--namespace)
                NAMESPACE="$2"
                shift 2
                ;;
            -t|--timeout)
                TIMEOUT="$2"
                shift 2
                ;;
            --skip-infrastructure)
                SKIP_INFRASTRUCTURE=true
                shift
                ;;
            --skip-wait)
                SKIP_WAIT=true
                shift
                ;;
            --skip-test)
                SKIP_TEST=true
                shift
                ;;
            --cleanup-on-failure)
                CLEANUP_ON_FAILURE=true
                shift
                ;;
            *)
                log_error "Unknown option: $1"
                show_usage
                exit 1
                ;;
        esac
    done
}

# Main deployment function
main() {
    log_info "Starting Lushop Kubernetes deployment"
    log_info "Namespace: $NAMESPACE"
    log_info "Timeout: ${TIMEOUT}s"

    # Set trap for cleanup on failure
    if [ "${CLEANUP_ON_FAILURE:-false}" = "true" ]; then
        trap cleanup ERR
    fi

    # Change to script directory
    cd "$SCRIPT_DIR"

    # Run deployment steps
    check_prerequisites
    setup_prerequisites
    deploy_infrastructure
    configure_nacos
    deploy_application

    if [ "${SKIP_WAIT:-false}" != "true" ]; then
        wait_for_services
    fi

    show_status
    get_endpoints

    if [ "${SKIP_TEST:-false}" != "true" ]; then
        test_services
    fi

    log_success "Lushop deployment completed successfully!"
    log_info "You can now access your services via the endpoints shown above."
}

# Run main function with all arguments
parse_args "$@"
main
