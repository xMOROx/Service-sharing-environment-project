#!/usr/bin/env bash
set -e

BANNER=$(
  cat <<'EOF'
 ██████╗██╗     ███████╗ █████╗ ███╗   ██╗    ██╗   ██╗██████╗ 
██╔════╝██║     ██╔════╝██╔══██╗████╗  ██║    ██║   ██║██╔══██╗
██║     ██║     █████╗  ███████║██╔██╗ ██║    ██║   ██║██████╔╝
██║     ██║     ██╔══╝  ██╔══██║██║╚██╗██║    ██║   ██║██╔═══╝ 
╚██████╗███████╗███████╗██║  ██║██║ ╚████║    ╚██████╔╝██║     
 ╚═════╝╚══════╝╚══════╝╚═╝  ╚═╝╚═╝  ╚═══╝     ╚═════╝ ╚═╝     
EOF
)

# shellcheck disable=SC1091
source "$(dirname "$0")/../utils.sh"
# shellcheck disable=SC1091
source "$(dirname "$0")/../config.sh"

show_banner_from_variable

log_step "Uninstalling application..."
helm uninstall $HELM_LOKI_DEPLOYMENT_NAME -n $LOKI_NAMESPACE || true
helm uninstall $HELM_TEMPO_DEPLOYMENT_NAME -n $TEMPO_NAMESPACE || true
helm uninstall $HELM_PROMTAIL_DEPLOYMENT_NAME -n $PROMTAIL_NAMESPACE || true
helm uninstall $HELM_EVENTEXPORTER_DEPLOYMENT_NAME -n $EVENTEXPORTER_NAMESPACE || true
helm uninstall $HELM_KUBE_PROMETHEUS_STACK_DEPLOYMENT_NAME -n $KUBE_PROMETHEUS_STACK_NAMESPACE || true
helm uninstall $HELM_OTEL_DEPLOYMENT_NAME -n $OTEL_COLLECTOR_NAMESPACE || true
helm uninstall $HELM_APP_DEPLOYMENT_NAME -n $APP_NAMESPACE || true
log_success "Application uninstalled."

log_info "Uninstallation complete."
