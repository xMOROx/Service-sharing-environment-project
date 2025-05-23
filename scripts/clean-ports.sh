#!/bin/bash
set -e

PROJECT_ROOT=$(git rev-parse --show-toplevel)

# shellcheck disable=SC1091
source "$PROJECT_ROOT/scripts/utils.sh"

log_step "Finding prometheus process id"
# Allow pgrep to fail without exiting the script (due to set -e) by using '|| true'
prometheus_process_id=$(pgrep -f "kubectl port-forward .*prometheus.*$PROMETHEUS_PORT" || true)

if [ -n "$prometheus_process_id" ]; then
  log_success "Found prometheus process id(s): $prometheus_process_id"
  log_step "Killing prometheus process(es)"
  # Unquote variable to allow kill to handle multiple PIDs if pgrep returns them
  if kill $prometheus_process_id; then
    log_success "Killed prometheus process id(s): $prometheus_process_id"
  else
    log_error "Failed to kill prometheus process id(s): $prometheus_process_id. They might have already exited or you may lack permissions."
  fi
else
  log_info "Prometheus port-forward process not found."
fi

echo

log_step "Finding grafana process id"
grafana_process_id=$(pgrep -f "kubectl port-forward .*grafana.*$GRAFANA_PORT" || true)

if [ -n "$grafana_process_id" ]; then
  log_success "Found grafana process id(s): $grafana_process_id"
  log_step "Killing grafana process(es)"
  if kill $grafana_process_id; then
    log_success "Killed grafana process id(s): $grafana_process_id"
  else
    log_error "Failed to kill grafana process id(s): $grafana_process_id. They might have already exited or you may lack permissions."
  fi
else
  log_info "Grafana port-forward process not found."
fi
