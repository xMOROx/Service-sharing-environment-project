#!/bin/bash
set -e

BANNER="""
██████╗ ███████╗██████╗ ██╗      ██████╗ ██╗   ██╗
██╔══██╗██╔════╝██╔══██╗██║     ██╔═══██╗╚██╗ ██╔╝
██║  ██║█████╗  ██████╔╝██║     ██║   ██║ ╚████╔╝
██║  ██║██╔══╝  ██╔═══╝ ██║     ██║   ██║  ╚██╔╝
██████╔╝███████╗██║     ███████╗╚██████╔╝   ██║
╚═════╝ ╚══════╝╚═╝     ╚══════╝ ╚═════╝    ╚═╝
"""

PROJECT_ROOT=$(git rev-parse --show-toplevel)
cd "$PROJECT_ROOT" || exit 1

# shellcheck disable=SC1091
source "$PROJECT_ROOT/scripts/utils.sh"

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

log_step "Deploying application..."
helm install demo ./infrastructure/deployment/ -n default
log_success "Application deployed."

log_info "Deployment complete. Make sure your Kubernetes cluster is configured and accessible."
