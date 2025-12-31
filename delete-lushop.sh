set -euo pipefail

NS=lushop
TS=$(date +%Y%m%dT%H%M%S)
OUTDIR="./${NS}-backup-${TS}"
mkdir -p "$OUTDIR"

echo "1/3: Exporting all resources in namespace '$NS' to $OUTDIR"
kubectl get all,cm,secret,ing,sts,daemonset,cronjob,pvc -n "$NS" -o yaml > "$OUTDIR/${NS}-all-resources.yaml" || true
kubectl get roles,rolebindings,serviceaccounts -n "$NS" -o yaml > "$OUTDIR/${NS}-rbac.yaml" || true
echo "Backup saved."

echo "2/3: Deleting Lushop-related resources in namespace '$NS' (deployments, services, ingresses, configmaps, cronjobs, statefulsets)"

# Delete controllers (deployments/statefulsets/daemonsets/replicasets/jobs/cronjobs) - delete all in namespace
kubectl -n "$NS" delete deployment,sts,daemonset,replicaset,job,cronjob --all --ignore-not-found || true

# Delete pods (force remove running pods)
kubectl -n "$NS" delete pod --all --ignore-not-found || true

# Delete services and ingresses
kubectl -n "$NS" delete svc --all --ignore-not-found || true
kubectl -n "$NS" delete ingress --all --ignore-not-found || true

# Delete configmaps that follow the '*-config' naming convention (common in this repo)
kubectl -n "$NS" get configmap -o name | grep -E '\-config$' | xargs -r -n1 kubectl -n "$NS" delete || true

# Delete most secrets with a safe filter to avoid service account tokens
# This deletes secrets ending with '-secret' or containing 'auth' or 'credentials'
kubectl -n "$NS" get secret -o name | grep -E '(-secret$|auth|credentials)' | xargs -r -n1 kubectl -n "$NS" delete || true

# Optionally delete PVCs if DELETE_PVCS env var is set to "true"
if [ "${DELETE_PVCS:-false}" = "true" ]; then
  echo "DELETE_PVCS=true: Deleting PVCs in namespace $NS"
  kubectl -n "$NS" delete pvc --all --ignore-not-found || true
fi

echo "Requested resources deletion initiated. Waiting a short while for resources to terminate..."
sleep 5

echo "3/3: Remaining resources in namespace (if any):"
kubectl get all -n "$NS" || echo "No remaining standard resources or namespace does not exist."
kubectl get configmap -n "$NS" || true
kubectl get secret -n "$NS" || true

echo ""
echo "If you want to completely remove the namespace and all leftovers (including PVCs/PVs), run:"
echo "  kubectl delete namespace $NS"
echo "Or to force-clean finalizers on stuck resources, inspect 'kubectl get <kind> -n $NS' and remove finalizers as needed."