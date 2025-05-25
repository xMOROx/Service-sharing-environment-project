#!/usr/bin/env bash
set -e

PROJECT_ROOT=$(git rev-parse --show-toplevel)

# shellcheck disable=SC1091
source "$PROJECT_ROOT/scripts/utils.sh"

GRAFANA_PORT=3000

log_step "Forwarding grafana port"
kubectl port-forward deployment/prometheus-grafana -n prometheus $GRAFANA_PORT
log_success "Forwarding grafana port to localhost:$GRAFANA_PORT"
