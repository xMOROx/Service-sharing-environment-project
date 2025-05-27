#!/usr/bin/env bash
set -e

# shellcheck disable=SC1091
source "$(dirname "$0")/../utils.sh"
# shellcheck disable=SC1091
source "$(dirname "$0")/../config.sh"

log_step "Setting up local registry..."
reg_name='kind-registry'
reg_port='5001'
if [ "$(docker inspect -f '{{.State.Running}}' "${reg_name}" 2>/dev/null || true)" != 'true' ]; then
  docker run \
    -d --restart=always -p "127.0.0.1:${reg_port}:5000" --network bridge --name "${reg_name}" \
    registry:2
fi
log_success "Local registry container is running."

log_step "Creating Kubernetes cluster with Kind and local registry configuration..."
kind create cluster --config "$KIND_CONFIG_PATH" --name "$KIND_CLUSTER_NAME"
log_success "Kubernetes cluster created with given name $KIND_CLUSTER_NAME."

log_step "Configuring nodes to use the local registry..."
REGISTRY_DIR="/etc/containerd/certs.d/localhost:${reg_port}"
for node in $(kind get nodes --name "$KIND_CLUSTER_NAME"); do
  docker exec "${node}" mkdir -p "${REGISTRY_DIR}"
  cat <<EOF | docker exec -i "${node}" cp /dev/stdin "${REGISTRY_DIR}/hosts.toml"
[host."http://${reg_name}:5000"]
EOF
done
log_success "Nodes configured to use the local registry."

log_step "Connecting the registry to the cluster network..."
if [ "$(docker inspect -f='{{json .NetworkSettings.Networks.kind}}' "${reg_name}")" = 'null' ]; then
  docker network connect "kind" "${reg_name}"
fi
log_success "Registry connected to the cluster network."

log_step "Applying local registry hosting ConfigMap..."
kubectl apply -f "${PROJECT_ROOT}/infrastructure/kind/local-registry-configmap.yaml"
log_success "Local registry hosting ConfigMap applied."
