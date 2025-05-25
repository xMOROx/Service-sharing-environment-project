#!/usr/bin/env bash
set -e

PROJECT_ROOT=$(git rev-parse --show-toplevel)
KIND_CLUSTER_NAME="k8s-lab"
cd "$PROJECT_ROOT" || exit 1

# shellcheck disable=SC1091
source "$PROJECT_ROOT/scripts/utils.sh"

log_step "Removing cluster..."
kind delete cluster --name "$KIND_CLUSTER_NAME"
log_success "Cluster removed successfully."
