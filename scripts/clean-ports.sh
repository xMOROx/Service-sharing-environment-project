#!/bin/bash
set -e

PROJECT_ROOT=$(git rev-parse --show-toplevel)

# shellcheck disable=SC1091
source "$PROJECT_ROOT/scripts/utils.sh"

log_step "Finding prometheus process id"
prometheus_process_id=$(pgrep -f "kubectl port-forward .*prometheus.*$PROMETHEUS_PORT" || exit 0)
if [ -z "$prometheus_process_id" ]; then
  log_info "Prometheus process not found. Exiting."
  exit 0
fi
log_success "Found prometheus process id: $prometheus_process_id"

log_step "Killing prometheus process"
kill "$prometheus_process_id"
log_success "Killed prometheus process id: $prometheus_process_id"
