#!/bin/bash
set -e

PROJECT_ROOT=$(git rev-parse --show-toplevel)

# shellcheck disable=SC1091
source "$PROJECT_ROOT/scripts/utils.sh"

PROMETHEUS_PORT=9090
GRAFANA_PORT=3000

log_step "Finding prometheus server name"
prometheus_pod=$(kubectl get pods -n default -l app.kubernetes.io/part-of=prometheus -l app.kubernetes.io/component=server -o jsonpath='{.items[*].metadata.name}')
log_success "Found prometheus pod: $prometheus_pod"

log_step "Forwarding prometheus port"
kubectl port-forward "${prometheus_pod}" -n default $PROMETHEUS_PORT:$PROMETHEUS_PORT &
log_success "Forwarding prometheus port to localhost:$PROMETHEUS_PORT"

log_step "Finding grafana pod name"
grafana_pod=$(kubectl get pods -n default -l app.kubernetes.io/name=grafana -o jsonpath='{.items[*].metadata.name}')
log_success "Found grafana pod: $grafana_pod"

log_step "Forwarding grafana port"
kubectl port-forward "${grafana_pod}" -n default $GRAFANA_PORT:$GRAFANA_PORT &
log_success "Forwarding grafana port to localhost:$GRAFANA_PORT"
