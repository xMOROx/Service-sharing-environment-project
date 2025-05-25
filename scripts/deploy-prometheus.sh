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

NAMESPACE="prometheus"
PROJECT_ROOT=$(git rev-parse --show-toplevel)
DEPLOYMENT_NAME="prometheus"
cd "$PROJECT_ROOT" || exit 1

# shellcheck disable=SC1091
source "$PROJECT_ROOT/scripts/utils.sh"

show_banner_from_variable

log_step "Adding prometheus helm repository..."
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
log_success "Successfully added prometheus helm repository."

log_step "Deploying prometheus to the Kubernetes cluster..."
kubectl create namespace $NAMESPACE || true
helm upgrade --install $DEPLOYMENT_NAME prometheus-community/kube-prometheus-stack -n $NAMESPACE --values $PROJECT_ROOT/infrastructure/deployment/kube-prometheus-stack.yaml
log_success "Prometheus deployed successfully to the Kubernetes cluster in namespace $NAMESPACE."
