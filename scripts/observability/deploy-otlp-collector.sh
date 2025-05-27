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

# shellcheck disable=SC1091
source "$(dirname "$0")/../utils.sh"
# shellcheck disable=SC1091
source "$(dirname "$0")/../config.sh"

show_banner_from_variable

log_step "Deploying otel-collector to the Kubernetes cluster..."
kubectl create namespace "$OTEL_COLLECTOR_NAMESPACE" || true

helm upgrade \
  --install "$HELM_OTEL_DEPLOYMENT_NAME" \
  --namespace "$OTEL_COLLECTOR_NAMESPACE" \
  --values "$OTEL_COLLECTOR_VALUES_PATH" \
  --version "$HELM_OTEL_COLLECTOR_VERSION" \
  open-telemetry/opentelemetry-collector

log_success "Otel-collector deployed successfully to the Kubernetes cluster in namespace $OTEL_COLLECTOR_NAMESPACE."
