#!/usr/bin/env bash
set -e

BANNER=$(
  cat <<EOF
██████╗ ██╗   ██╗██████╗  ██████╗ ███████╗ ██████╗ ██████╗ ██████╗ ███████╗
██╔══██╗╚██╗ ██╔╝██╔══██╗██╔═══██╗██╔════╝██╔════╝██╔═══██╗██╔══██╗██╔════╝
██████╔╝ ╚████╔╝ ██████╔╝██║   ██║███████╗██║     ██║   ██║██████╔╝█████╗  
██╔═══╝   ╚██╔╝  ██╔══██╗██║   ██║╚════██║██║     ██║   ██║██╔═══╝ ██╔══╝  
██║        ██║   ██║  ██║╚██████╔╝███████║╚██████╗╚██████╔╝██║     ███████╗
╚═╝        ╚═╝   ╚═╝  ╚═╝ ╚═════╝ ╚══════╝ ╚═════╝ ╚═════╝ ╚═╝     ╚══════╝
EOF
)

# shellcheck disable=SC1091
source "$(dirname "$0")/../utils.sh"
# shellcheck disable=SC1091
source "$(dirname "$0")/../config.sh"

show_banner_from_variable

log_step "Deploying grafana/pyroscope to the Kubernetes cluster..."
kubectl create namespace "$PYROSCOPE_NAMESPACE" || true

helm upgrade \
  --install "$HELM_PYROSCOPE_DEPLOYMENT_NAME" \
  --namespace="$PYROSCOPE_NAMESPACE" \
  --values="$PYROSCOPE_VALUES_PATH" \
  --version="$HELM_PYROSCOPE_CHART_VERSION" \
  grafana/pyroscope

log_success "Pyroscope deployed successfully to the Kubernetes cluster in namespace $PYROSCOPE_NAMESPACE."
