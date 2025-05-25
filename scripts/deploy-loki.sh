#!/usr/bin/env bash
set -e

BANNER=$(
  cat <<EOF
██╗      ██████╗ ██╗  ██╗██╗
██║     ██╔═══██╗██║ ██╔╝██║
██║     ██║   ██║█████╔╝ ██║
██║     ██║   ██║██╔═██╗ ██║
███████╗╚██████╔╝██║  ██╗██║
╚══════╝ ╚═════╝ ╚═╝  ╚═╝╚═╝
EOF
)

NAMESPACE="loki"
DEPLOYMENT_NAME="loki"

PROJECT_ROOT=$(git rev-parse --show-toplevel)
cd "$PROJECT_ROOT" || exit 1

# shellcheck disable=SC1091
source "$PROJECT_ROOT/scripts/utils.sh"

show_banner_from_variable

log_step "Adding grafana/loki helm repository..."
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update
log_success "Successfully added grafana/loki helm repository."

log_step "Deploying grafana/loki to the Kubernetes cluster..."
kubectl create namespace $NAMESPACE || true
helm install --values $PROJECT_ROOT/infrastructure/deployment/loki.yaml $DEPLOYMENT_NAME --namespace=$NAMESPACE grafana/loki
log_success "Grafana/loki deployed successfully to the Kubernetes cluster in namespace $NAMESPACE."
