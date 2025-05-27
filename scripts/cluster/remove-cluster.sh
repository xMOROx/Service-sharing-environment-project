#!/usr/bin/env bash
set -e

# shellcheck disable=SC1091
source "$(dirname "$0")/../utils.sh"
# shellcheck disable=SC1091
source "$(dirname "$0")/../config.sh"

log_step "Removing cluster..."
kind delete cluster --name "$KIND_CLUSTER_NAME"
log_success "Cluster removed successfully."
