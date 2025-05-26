#!/usr/bin/env bash
set -e

BANNER=$(
  cat <<EOF
████████╗███████╗███╗   ███╗██████╗  ██████╗ 
╚══██╔══╝██╔════╝████╗ ████║██╔══██╗██╔═══██╗
   ██║   █████╗  ██╔████╔██║██████╔╝██║   ██║
   ██║   ██╔══╝  ██║╚██╔╝██║██╔═══╝ ██║   ██║
   ██║   ███████╗██║ ╚═╝ ██║██║     ╚██████╔╝
   ╚═╝   ╚══════╝╚═╝     ╚═╝╚═╝      ╚═════╝ 
EOF
)

# shellcheck disable=SC1091
source "$(dirname "$0")/../utils.sh"
# shellcheck disable=SC1091
source "$(dirname "$0")/../config.sh"

show_banner_from_variable

log_step "Deploying grafana/tempo to the Kubernetes cluster..."
kubectl create namespace $TEMPO_NAMESPACE || true
helm install --values $PROJECT_ROOT/infrastructure/deployment/tempo.yaml $HELM_TEMPO_DEPLOYMENT_NAME --namespace=$TEMPO_NAMESPACE grafana/tempo-distributed
log_success "Grafana/tempo deployed successfully to the Kubernetes cluster in namespace $TEMPO_NAMESPACE."
