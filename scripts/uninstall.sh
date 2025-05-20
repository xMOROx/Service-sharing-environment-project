#!/bin/bash
set -e

BANNER=$(cat <<'EOF'
 ██████╗██╗     ███████╗ █████╗ ███╗   ██╗    ██╗   ██╗██████╗ 
██╔════╝██║     ██╔════╝██╔══██╗████╗  ██║    ██║   ██║██╔══██╗
██║     ██║     █████╗  ███████║██╔██╗ ██║    ██║   ██║██████╔╝
██║     ██║     ██╔══╝  ██╔══██║██║╚██╗██║    ██║   ██║██╔═══╝ 
╚██████╗███████╗███████╗██║  ██║██║ ╚████║    ╚██████╔╝██║     
 ╚═════╝╚══════╝╚══════╝╚═╝  ╚═╝╚═╝  ╚═══╝     ╚═════╝ ╚═╝     
EOF
)

PROJECT_ROOT=$(git rev-parse --show-toplevel)
cd "$PROJECT_ROOT" || exit 1

# shellcheck disable=SC1091
source "$PROJECT_ROOT/scripts/utils.sh"

show_banner_from_variable

log_step "Uninstalling application..."
helm uninstall demo -n default
log_success "Application uninstalled."

log_info "Uninstallation complete."
