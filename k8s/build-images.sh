#!/bin/bash

# Lushop Docker Image Build Script
# This script builds Docker images for Lushop services

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuration
REGISTRY="${REGISTRY:-lushop}"
TAG="${TAG:-latest}"
SERVICES=("gateway" "user" "goods" "order" "inventory" "userauth" "userop")

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

# Check Docker
check_docker() {
    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed"
        exit 1
    fi

    if ! docker info &> /dev/null; then
        log_error "Docker daemon is not running"
        exit 1
    fi
}

# Build service image
build_service() {
    local service="$1"
    local image_name="${REGISTRY}/${service}:${TAG}"

    log_info "Building image for ${service}..."

    # Determine Dockerfile location (relative to project root)
    if [ "${service}" = "gateway" ] && [ -f "../lushop/Dockerfile" ]; then
        DOCKERFILE="../lushop/Dockerfile"
        CONTEXT="../lushop"
    elif [ -f "../service/${service}/Dockerfile" ]; then
        DOCKERFILE="../service/${service}/Dockerfile"
        CONTEXT="../service/${service}"
    else
        log_warning "Dockerfile not found for ${service}, skipping..."
        return 1
    fi

    # Build image
    docker build -t "${image_name}" -f "${DOCKERFILE}" "${CONTEXT}"

    log_success "Built ${image_name}"
    return 0
}

# Push image
push_service() {
    local service="$1"
    local image_name="${REGISTRY}/${service}:${TAG}"

    log_info "Pushing ${image_name}..."
    docker push "${image_name}"
    log_success "Pushed ${image_name}"
}

# List images
list_images() {
    log_info "Local Lushop images:"
    echo ""
    docker images | grep lushop || echo "No lushop images found"
    echo ""
}

# Clean images
clean_images() {
    log_warning "This will remove all lushop images. Continue? (y/N)"
    read -r response
    if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
        log_info "Cleaning images..."
        docker images | grep lushop | awk '{print $3}' | xargs docker rmi -f 2>/dev/null || true
        log_success "Images cleaned"
    fi
}

# Build all services
build_all() {
    log_info "Building all services..."

    local built_count=0
    local total_count=${#SERVICES[@]}

    for service in "${SERVICES[@]}"; do
        if build_service "${service}"; then
            built_count=$((built_count + 1))
        fi
    done

    log_success "Built ${built_count}/${total_count} services"
}

# Push all services
push_all() {
    log_info "Pushing all services..."

    for service in "${SERVICES[@]}"; do
        push_service "${service}"
    done

    log_success "All services pushed"
}

# Main function
main() {
    check_docker

    case "$1" in
        all)
            build_all
            ;;
        push)
            if [ -z "$2" ]; then
                push_all
            else
                push_service "$2"
            fi
            ;;
        list)
            list_images
            ;;
        clean)
            clean_images
            ;;
        gateway)
            build_service "lushop"
            ;;
        "")
            echo "Usage: $0 {all|push [service]|list|clean|gateway|<service_name>}"
            echo ""
            echo "Commands:"
            echo "  all               Build all services"
            echo "  push [service]    Push all services or specific service"
            echo "  list              List all built images"
            echo "  clean             Remove all lushop images"
            echo "  gateway           Build only gateway service"
            echo "  <service_name>    Build specific service (goods, inventory, order, user, userauth, userop)"
            echo ""
            echo "Environment variables:"
            echo "  REGISTRY          Docker registry (default: localhost:5000)"
            echo "  TAG               Image tag (default: latest)"
            echo ""
            exit 1
            ;;
        *)
            # Build specific service
            service="$1"
            if [[ " ${SERVICES[*]} " =~ " ${service} " ]]; then
                build_service "${service}"
            else
                log_error "Unknown service: ${service}"
                echo "Available services: ${SERVICES[*]}"
                exit 1
            fi
            ;;
    esac
}

main "$@"
