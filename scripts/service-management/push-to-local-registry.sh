#!/usr/bin/env bash
set -e

# shellcheck disable=SC1091
source "$(dirname "$0")/../utils.sh"
# shellcheck disable=SC1091
source "$(dirname "$0")/../config.sh"

log_step "Checking for local registry..."
if ! docker inspect -f '{{.State.Running}}' "kind-registry" 2>/dev/null | grep -q "true"; then
  log_warning "Local registry 'kind-registry' is not running. Skipping image push to local registry."
  exit 0
fi
log_success "Local registry 'kind-registry' found."

log_step "Pushing Docker images to local registry..."
docker push localhost:5001/order-service:0.1.0
docker push localhost:5001/inventory-service:0.1.0
log_success "Docker images pushed to local registry."
