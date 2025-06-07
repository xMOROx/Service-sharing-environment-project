#!/usr/bin/env bash
set -e

BANNER=$(
  cat <<EOF
██████╗ ██████╗  ██████╗ ███╗   ███╗███████╗████████╗██╗  ██╗███████╗██╗   ██╗███████╗
██╔══██╗██╔══██╗██╔═══██╗████╗ ████║██╔════╝╚══██╔══╝██║  ██║██╔════╝██║   ██║██╔════╝
██████╔╝██████╔╝██║   ██║██╔████╔██║█████╗     ██║   ███████║█████╗  ██║   ██║███████╗
██╔═══╝ ██╔══██╗██║   ██║██║╚██╔╝██║██╔══╝     ██║   ██╔══██║██╔══╝  ██║   ██║╚════██║
██║     ██║  ██║╚██████╔╝██║ ╚═╝ ██║███████╗   ██║   ██║  ██║███████╗╚██████╔╝███████║
╚═╝     ╚═╝  ╚═╝ ╚═════╝ ╚═╝     ╚═╝╚══════╝   ╚═╝   ╚═╝  ╚═╝╚══════╝ ╚═════╝ ╚══════╝
EOF
)

# shellcheck disable=SC1091
source "$(dirname "$0")/../utils.sh"
# shellcheck disable=SC1091
source "$(dirname "$0")/../config.sh"

show_banner_from_variable

log_step "Deploying prometheus to the Kubernetes cluster..."
kubectl create namespace "$KUBE_PROMETHEUS_STACK_NAMESPACE" || true
if ! kubectl get configmap grafana-dashboards -n "$KUBE_PROMETHEUS_STACK_NAMESPACE" >/dev/null 2>&1; then
  kubectl create configmap grafana-dashboards --from-file="$ORDER_DASHBOARD" --from-file="$INVENTORY_DASHBOARD" -n "$KUBE_PROMETHEUS_STACK_NAMESPACE"
	echo "ConfigMap 'grafana-dashboards' created."
fi

helm upgrade \
  --install "$HELM_KUBE_PROMETHEUS_STACK_DEPLOYMENT_NAME" \
  --namespace "$KUBE_PROMETHEUS_STACK_NAMESPACE" \
  --values "$KUBE_PROMETHEUS_STACK_VALUES_PATH" \
  --version "$HELM_KUBE_PROMETHEUS_STACK_CHART_VERSION" \
  prometheus-community/kube-prometheus-stack

log_success "Prometheus deployed successfully to the Kubernetes cluster in namespace $KUBE_PROMETHEUS_STACK_NAMESPACE."
