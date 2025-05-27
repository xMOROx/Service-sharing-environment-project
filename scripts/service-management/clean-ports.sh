#!/usr/bin/env bash
set -e

# shellcheck disable=SC1091
source "$(dirname "$0")/../utils.sh"
# shellcheck disable=SC1091
source "$(dirname "$0")/../config.sh"

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
