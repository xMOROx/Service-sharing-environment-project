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

# shellcheck disable=SC1091
source "$(dirname "$0")/../utils.sh"
# shellcheck disable=SC1091
source "$(dirname "$0")/../config.sh"

show_banner_from_variable

log_step "Deploying grafana/promtail to the Kubernetes cluster..."
kubectl create namespace $PROMTAIL_NAMESPACE || true
helm upgrade --values $PROMTAIL_VALUES_PATH --install $HELM_PROMTAIL_DEPLOYMENT_NAME grafana/promtail -n $PROMTAIL_NAMESPACE
log_success "Grafana/promtail deployed successfully to the Kubernetes cluster in namespace $PROMTAIL_NAMESPACE."
