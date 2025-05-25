#!/usr/bin/env bash
set -e

PROJECT_ROOT=$(git rev-parse --show-toplevel)
cd "$PROJECT_ROOT" || exit 1

# shellcheck disable=SC1091
source "$PROJECT_ROOT/scripts/utils.sh"

log_step "Creating Kubernetes cluster with Kind..."
kind create cluster --config $PROJECT_ROOT/infrastructure/kind/kind.yaml --name k8s-lab
log_success "Kubernetes cluster created with given name k8s-lab."
