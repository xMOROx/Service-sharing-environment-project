#!/usr/bin/env bash
set -e

BANNER=$(
  cat <<EOF
██████╗ ██████╗  ██████╗ ███╗   ███╗████████╗ █████╗ ██╗██╗     
██╔══██╗██╔══██╗██╔═══██╗████╗ ████║╚══██╔══╝██╔══██╗██║██║     
██████╔╝██████╔╝██║   ██║██╔████╔██║   ██║   ███████║██║██║     
██╔═══╝ ██╔══██╗██║   ██║██║╚██╔╝██║   ██║   ██╔══██║██║██║     
██║     ██║  ██║╚██████╔╝██║ ╚═╝ ██║   ██║   ██║  ██║██║███████╗
╚═╝     ╚═╝  ╚═╝ ╚═════╝ ╚═╝     ╚═╝   ╚═╝   ╚═╝  ╚═╝╚═╝╚══════╝
                                                                
EOF
)

NAMESPACE="promtail"
PROJECT_ROOT=$(git rev-parse --show-toplevel)
DEPLOYMENT_NAME="promtail"
cd "$PROJECT_ROOT" || exit 1

# shellcheck disable=SC1091
source "$PROJECT_ROOT/scripts/utils.sh"

show_banner_from_variable

log_step "Adding grafana/promtail helm repository..."
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update
log_success "Successfully added grafana/promtail helm repository."

log_step "Deploying grafana/promtail to the Kubernetes cluster..."
kubectl create namespace $NAMESPACE || true
helm upgrade --values $PROJECT_ROOT/infrastructure/deployment/promtail.yaml --install $DEPLOYMENT_NAME grafana/promtail -n $NAMESPACE
log_success "Grafana/promtail deployed successfully to the Kubernetes cluster in namespace $NAMESPACE."
