#!/usr/bin/env bash
set -e

BANNER=$(
  cat <<EOF
 ██████╗ ████████╗██╗     ██████╗      ██████╗ ██████╗ ██╗     ██╗     ███████╗ ██████╗████████╗ ██████╗ ██████╗ 
██╔═══██╗╚══██╔══╝██║     ██╔══██╗    ██╔════╝██╔═══██╗██║     ██║     ██╔════╝██╔════╝╚══██╔══╝██╔═══██╗██╔══██╗
██║   ██║   ██║   ██║     ██████╔╝    ██║     ██║   ██║██║     ██║     █████╗  ██║        ██║   ██║   ██║██████╔╝
██║   ██║   ██║   ██║     ██╔═══╝     ██║     ██║   ██║██║     ██║     ██╔══╝  ██║        ██║   ██║   ██║██╔══██╗
╚██████╔╝   ██║   ███████╗██║         ╚██████╗╚██████╔╝███████╗███████╗███████╗╚██████╗   ██║   ╚██████╔╝██║  ██║
 ╚═════╝    ╚═╝   ╚══════╝╚═╝          ╚═════╝ ╚═════╝ ╚══════╝╚══════╝╚══════╝ ╚═════╝   ╚═╝    ╚═════╝ ╚═╝  ╚═╝
                                                                                                                 
EOF
)

NAMESPACE="otlp-collector"
PROJECT_ROOT=$(git rev-parse --show-toplevel)
DEPLOYMENT_NAME="otlp-collector"
cd "$PROJECT_ROOT" || exit 1

# shellcheck disable=SC1091
source "$PROJECT_ROOT/scripts/utils.sh"

show_banner_from_variable

log_step "Adding otlp-collector helm repository..."
helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts
helm repo update
log_success "Successfully added otlp-collector helm repository."

log_step "Deploying otlp-collector to the Kubernetes cluster..."
kubectl create namespace $NAMESPACE || true
helm install $DEPLOYMENT_NAME open-telemetry/opentelemetry-collector \
  --namespace $NAMESPACE \
  --values $PROJECT_ROOT/infrastructure/deployment/otlp-collector.yaml
log_success "Otlp-collector deployed successfully to the Kubernetes cluster in namespace $NAMESPACE."
