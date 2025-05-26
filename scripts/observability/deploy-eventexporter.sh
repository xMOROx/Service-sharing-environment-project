#!/usr/bin/env bash
set -e

BANNER=$(
  cat <<EOF
███████╗██╗   ██╗███████╗███╗   ██╗████████╗    ███████╗██╗  ██╗██████╗  ██████╗ ██████╗ ████████╗███████╗██████╗ 
██╔════╝██║   ██║██╔════╝████╗  ██║╚══██╔══╝    ██╔════╝╚██╗██╔╝██╔══██╗██╔═══██╗██╔══██╗╚══██╔══╝██╔════╝██╔══██╗
█████╗  ██║   ██║█████╗  ██╔██╗ ██║   ██║       █████╗   ╚███╔╝ ██████╔╝██║   ██║██████╔╝   ██║   █████╗  ██████╔╝
██╔══╝  ╚██╗ ██╔╝██╔══╝  ██║╚██╗██║   ██║       ██╔══╝   ██╔██╗ ██╔═══╝ ██║   ██║██╔══██╗   ██║   ██╔══╝  ██╔══██╗
███████╗ ╚████╔╝ ███████╗██║ ╚████║   ██║       ███████╗██╔╝ ██╗██║     ╚██████╔╝██║  ██║   ██║   ███████╗██║  ██║
╚══════╝  ╚═══╝  ╚══════╝╚═╝  ╚═══╝   ╚═╝       ╚══════╝╚═╝  ╚═╝╚═╝      ╚═════╝ ╚═╝  ╚═╝   ╚═╝   ╚══════╝╚═╝  ╚═╝
EOF
)

# shellcheck disable=SC1091
source "$(dirname "$0")/../utils.sh"
# shellcheck disable=SC1091
source "$(dirname "$0")/../config.sh"

show_banner_from_variable

log_step "Deploying event-exporter to the Kubernetes cluster..."
kubectl create namespace $EVENTEXPORTER_NAMESPACE || true
helm install $HELM_EVENTEXPORTER_DEPLOYMENT_NAME bitnami/kubernetes-event-exporter --values $EVENTEXPORTER_VALUES_PATH -n $EVENTEXPORTER_NAMESPACE
log_success "Event-exporter deployed successfully to the Kubernetes cluster in namespace $EVENTEXPORTER_NAMESPACE."
