#!/bin/bash

# Lushop Local Development Deployment Script
# This script helps deploy Lushop services for local development

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

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

    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed. Please install Docker first."
        exit 1
    fi

    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose is not installed. Please install Docker Compose first."
        exit 1
    fi

    log_success "Prerequisites check passed"
}

# Deploy infrastructure
deploy_infrastructure() {
    log_info "Deploying infrastructure services..."

    # Deploy infrastructure services
    docker-compose up -d mysql redis consul nacos jaeger rocketmq-namesrv rocketmq-broker

    # Wait for services to be ready
    log_info "Waiting for MySQL..."
    docker-compose exec -T mysql mysqladmin ping -h localhost --silent --wait=30

    log_info "Waiting for Redis..."
    docker-compose exec -T redis redis-cli --raw incr ping > /dev/null

    log_info "Waiting for Consul..."
    timeout=60
    while [ $timeout -gt 0 ]; do
        if curl -s http://localhost:8500/v1/status/leader > /dev/null 2>&1; then
            break
        fi
        sleep 2
        timeout=$((timeout - 2))
    done

    if [ $timeout -le 0 ]; then
        log_error "Consul failed to start"
        exit 1
    fi

    log_info "Waiting for Nacos..."
    timeout=120
    while [ $timeout -gt 0 ]; do
        if curl -s http://localhost:8848/nacos/v1/console/health/readiness > /dev/null 2>&1; then
            break
        fi
        sleep 2
        timeout=$((timeout - 2))
    done

    if [ $timeout -le 0 ]; then
        log_error "Nacos failed to start"
        exit 1
    fi

    log_success "Infrastructure deployment completed"
}

# Deploy application services
deploy_services() {
    log_info "Deploying application services..."

    # Deploy all application services
    docker-compose up -d lushop-gateway goods-service inventory-service order-service user-service userauth-service userop-service

    log_success "Application services deployed"
}

# Stop all services
stop_services() {
    log_info "Stopping all services..."
    docker-compose down
    log_success "All services stopped"
}

# Show status
show_status() {
    log_info "Service status:"
    echo ""
    docker-compose ps
    echo ""
    echo "=== Service URLs ==="
    echo "Consul UI:      http://localhost:8500"
    echo "Nacos UI:       http://localhost:8848/nacos (nacos/nacos)"
    echo "Jaeger UI:      http://localhost:16686"
    echo "API Gateway:    http://localhost:8001"
    echo ""
}

# Show logs
show_logs() {
    service="$1"
    if [ -z "$service" ]; then
        log_info "Showing logs for all services (Ctrl+C to stop)..."
        docker-compose logs -f
    else
        log_info "Showing logs for $service (Ctrl+C to stop)..."
        docker-compose logs -f "$service"
    fi
}

# Clean up
cleanup() {
    log_warning "This will remove all containers and volumes. Continue? (y/N)"
    read -r response
    if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
        log_info "Cleaning up..."
        docker-compose down -v --remove-orphans
        docker system prune -f
        log_success "Cleanup completed"
    fi
}

# Main function
main() {
    case "$1" in
        start|deploy)
            check_prerequisites
            deploy_infrastructure
            deploy_services
            sleep 5
            show_status
            ;;
        stop)
            stop_services
            ;;
        restart)
            stop_services
            sleep 2
            deploy_infrastructure
            deploy_services
            sleep 5
            show_status
            ;;
        status)
            show_status
            ;;
        logs)
            show_logs "$2"
            ;;
        cleanup|clean)
            cleanup
            ;;
        infrastructure|infra)
            check_prerequisites
            deploy_infrastructure
            ;;
        services|apps)
            deploy_services
            ;;
        *)
            echo "Usage: $0 {deploy|start|stop|restart|status|logs [service]|cleanup|infrastructure|services}"
            echo ""
            echo "Commands:"
            echo "  deploy/start    Deploy all services (infrastructure + applications)"
            echo "  stop            Stop all services"
            echo "  restart         Restart all services"
            echo "  status          Show service status"
            echo "  logs [service]  Show logs (all services if no service specified)"
            echo "  cleanup/clean   Remove all containers and volumes"
            echo "  infrastructure  Deploy only infrastructure services"
            echo "  services        Deploy only application services"
            echo ""
            exit 1
            ;;
    esac
}

main "$@"
