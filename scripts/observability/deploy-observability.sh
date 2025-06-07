#!/usr/bin/env bash
set -e

# shellcheck disable=SC1091
source "$(dirname "$0")/../utils.sh"
# shellcheck disable=SC1091
source "$(dirname "$0")/../config.sh"

log_step "Adding required Helm repositories..."
helm repo add grafana https://grafana.github.io/helm-charts
helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add bitnami https://charts.bitnami.com/bitnami
log_step "Updating Helm repositories..."
helm repo update
log_success "Helm repositories updated successfully."

log_step "Deploying Prometheus..."
"$PROJECT_ROOT/scripts/observability/deploy-prometheus.sh"
log_step "Deploying Loki..."
"$PROJECT_ROOT/scripts/observability/deploy-loki.sh"
log_step "Deploying Tempo..."
"$PROJECT_ROOT/scripts/observability/deploy-tempo.sh"
# log_step "Deploying Eventexporter..."
# "$PROJECT_ROOT/scripts/observability/deploy-eventexporter.sh"
log_step "Deploying Promtail..."
"$PROJECT_ROOT/scripts/observability/deploy-promtail.sh"
log_step "Deploying Pyroscope..."
"$PROJECT_ROOT/scripts/observability/deploy-pyroscope.sh"
log_step "Deploying OTLP Collector..."
"$PROJECT_ROOT/scripts/observability/deploy-otlp-collector.sh"

log_success "Observability stack deployment complete."
