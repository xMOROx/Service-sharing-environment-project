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

# shellcheck disable=SC1091
source "$(dirname "$0")/../utils.sh"
# shellcheck disable=SC1091
source "$(dirname "$0")/../config.sh"

show_banner_from_variable

log_step "Deploying grafana/loki to the Kubernetes cluster..."
kubectl create namespace "$LOKI_NAMESPACE" || true

helm upgrade \
  --install "$HELM_LOKI_DEPLOYMENT_NAME"  \
  --namespace="$LOKI_NAMESPACE" \
  --values "$LOKI_VALUES_PATH" \
  --version "$HELM_LOKI_VERSION" \
  grafana/loki

log_success "Grafana/loki deployed successfully to the Kubernetes cluster in namespace $LOKI_NAMESPACE."
