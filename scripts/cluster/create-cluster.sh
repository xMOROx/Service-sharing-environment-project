#!/usr/bin/env bash
set -e

# shellcheck disable=SC1091
source "$(dirname "$0")/../utils.sh"
# shellcheck disable=SC1091
source "$(dirname "$0")/../config.sh"

log_step "Creating Kubernetes cluster with Kind..."
kind create cluster --config "$KIND_CONFIG_PATH" --name "$KIND_CLUSTER_NAME"
log_success "Kubernetes cluster created with given name $KIND_CLUSTER_NAME."
