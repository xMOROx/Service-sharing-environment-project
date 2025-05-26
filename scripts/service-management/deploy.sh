#!/usr/bin/env bash
set -e

BANNER=$(
  cat <<EOF
██████╗ ███████╗██████╗ ██╗      ██████╗ ██╗   ██╗
██╔══██╗██╔════╝██╔══██╗██║     ██╔═══██╗╚██╗ ██╔╝
██║  ██║█████╗  ██████╔╝██║     ██║   ██║ ╚████╔╝
██║  ██║██╔══╝  ██╔═══╝ ██║     ██║   ██║  ╚██╔╝
██████╔╝███████╗██║     ███████╗╚██████╔╝   ██║
╚═════╝ ╚══════╝╚═╝     ╚══════╝ ╚═════╝    ╚═╝
EOF
)

# shellcheck disable=SC1091
source "$(dirname "$0")/../utils.sh"
# shellcheck disable=SC1091
source "$(dirname "$0")/../config.sh"

show_banner_from_variable

log_step "Running protobuf generation script..."
bash "$PROJECT_ROOT/scripts/generate-proto.sh"

log_step "Building order service..."
(cd services/order-service && go build .)
log_success "Order service built."

log_step "Building inventory service..."
(cd services/inventory-service && go build .)
log_success "Inventory service built."

log_step "Building Docker images..."
docker buildx bake
log_success "Docker images built."

log_step "Helm dependencies updating..."
helm dependency update ./infrastructure/deployment/
log_success "Helm dependencies updated."

log_step "Deploying observability plane ..."
$PROJECT_ROOT/scripts/observability/deploy-observability.sh
log_success "Observability plane deployed."

log_step "Deploying application..."
helm install demo ./infrastructure/deployment/ -n $APP_NAMESPACE --create-namespace || exit 1
log_success "Application deployed."

log_info "Deployment complete. Make sure your Kubernetes cluster is configured and accessible."
