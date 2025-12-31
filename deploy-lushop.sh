#!/bin/bash
set -euo pipefail

# deploy-lushop.sh
# One-click deploy for lushop: create namespace, apply manifests, wait for core services.

NS=lushop
K8S_DIR="./k8s"
BUILD_IMAGES=false

usage() {
  cat <<EOF
Usage: $0 [--build]
  --build    Build service images before applying manifests (requires ./k8s/build-images.sh)
EOF
}

if [ "${1:-}" = "--help" ]; then
  usage
  exit 0
fi

if [ "${1:-}" = "--build" ]; then
  BUILD_IMAGES=true
fi

echo "One-click deploy: namespace=${NS}, k8s dir=${K8S_DIR}, build_images=${BUILD_IMAGES}"

if [ "$BUILD_IMAGES" = true ]; then
  if [ -x "${K8S_DIR}/build-images.sh" ]; then
    echo "Building images..."
    (cd "$K8S_DIR" && ./build-images.sh all)
  else
    echo "Build script not found or not executable: ${K8S_DIR}/build-images.sh"
    exit 1
  fi
fi

echo "Creating namespace ${NS} (if not exists)..."
kubectl create namespace "$NS" --dry-run=client -o yaml | kubectl apply -f -

echo "Applying kustomize manifests..."
kubectl apply -k "$K8S_DIR"

echo "Waiting for core deployments to become available (timeout 5m)..."
for deploy in lushop-gateway user-service goods-service order-service inventory-service userauth-service userop-service; do
  echo "Waiting for ${deploy}..."
  kubectl -n "$NS" wait --for=condition=available "deployment/${deploy}" --timeout=300s || echo "Warning: ${deploy} not available within timeout"
done

echo "Deployment finished. Run 'kubectl get pods -n ${NS}' to check status."


