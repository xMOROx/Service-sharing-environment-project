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

NAMESPACE="event-exporter"
PROJECT_ROOT=$(git rev-parse --show-toplevel)
DEPLOYMENT_NAME="event-exporter"
cd "$PROJECT_ROOT" || exit 1

# shellcheck disable=SC1091
source "$PROJECT_ROOT/scripts/utils.sh"

show_banner_from_variable

log_step "Adding bitnami helm repository..."
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
log_success "Successfully added bitnami helm repository."

log_step "Deploying event-exporter to the Kubernetes cluster..."
kubectl create namespace $NAMESPACE || true
helm install $DEPLOYMENT_NAME bitnami/kubernetes-event-exporter --values $PROJECT_ROOT/infrastructure/deployment/eventexporter.yaml -n $NAMESPACE
log_success "Event-exporter deployed successfully to the Kubernetes cluster in namespace $NAMESPACE."
