#!/usr/bin/env bash
set -e

BANNER=$(
  cat <<EOF
██████╗ ██████╗  ██████╗ ████████╗ ██████╗ 
██╔══██╗██╔══██╗██╔═══██╗╚══██╔══╝██╔═══██╗
██████╔╝██████╔╝██║   ██║   ██║   ██║   ██║
██╔═══╝ ██╔══██╗██║   ██║   ██║   ██║   ██║
██║     ██║  ██║╚██████╔╝   ██║   ╚██████╔╝
╚═╝     ╚═╝  ╚═╝ ╚═════╝    ╚═╝    ╚═════╝ 
EOF
)

# shellcheck disable=SC1091
source "$(dirname "$0")/utils.sh"

show_banner_from_variable

log_step "Ensuring proto directories exist..."
mkdir -p proto/inventory
mkdir -p proto/order
log_success "Proto directories ensured."

log_step "Generating protobuf files..."
protoc --go_out=.. --go-grpc_out=.. --proto_path=proto proto/inventory.proto
protoc --go_out=.. --go-grpc_out=.. --proto_path=proto proto/order.proto
log_success "Protobuf files generated."

log_info "Protobuf generation complete."
