#!/usr/bin/env bash
set -e

# shellcheck disable=SC1091
source "$(dirname "$0")/../utils.sh"
# shellcheck disable=SC1091
source "$(dirname "$0")/../config.sh"

log_step "Forwarding grafana port"
kubectl port-forward deployment/prometheus-grafana -n $KUBE_PROMETHEUS_STACK_NAMESPACE $GRAFANA_PORT
log_success "Forwarding grafana port to localhost:$GRAFANA_PORT"
