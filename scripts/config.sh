#!/bin/bash

PROJECT_ROOT=$(git rev-parse --show-toplevel)

APP_NAMESPACE="demo"
LOKI_NAMESPACE="loki"
PROMTAIL_NAMESPACE="promtail"
TEMPO_NAMESPACE="tempo"
OTEL_COLLECTOR_NAMESPACE="otel-collector"
KUBE_PROMETHEUS_STACK_NAMESPACE="prometheus"
EVENTEXPORTER_NAMESPACE="event-exporter"
PYROSCOPE_NAMESPACE="pyroscope"

HELM_APP_DEPLOYMENT_NAME="demo"
HELM_LOKI_DEPLOYMENT_NAME="loki"
HELM_OTEL_DEPLOYMENT_NAME="otel-collector"
HELM_PROMTAIL_DEPLOYMENT_NAME="promtail"
HELM_TEMPO_DEPLOYMENT_NAME="tempo"
HELM_KUBE_PROMETHEUS_STACK_DEPLOYMENT_NAME="prometheus"
HELM_EVENTEXPORTER_DEPLOYMENT_NAME="event-exporter"
HELM_PYROSCOPE_DEPLOYMENT_NAME="pyroscope"

HELM_LOKI_VERSION="6.30.0"
HELM_OTEL_COLLECTOR_VERSION="0.125.0"
HELM_PROMTAIL_VERSION="6.16.6"
HELM_TEMPO_VERSION="1.40.2"
HELM_KUBE_PROMETHEUS_STACK_CHART_VERSION="72.6.2"
HELM_EVENTEXPORTER_VERSION="3.5.3"
HELM_PYROSCOPE_CHART_VERSION="1.13.4"

KIND_CLUSTER_NAME="demo-k8s-lab"
KIND_CONFIG_PATH="${PROJECT_ROOT}/infrastructure/kind/kind.yaml"

ORDER_SERVICE_NAME="order-service"
INVENTORY_SERVICE_NAME="inventory-service"

# Paths
ORDER_SERVICE_DIR="${PROJECT_ROOT}/services/${ORDER_SERVICE_NAME}"
INVENTORY_SERVICE_DIR="${PROJECT_ROOT}/services/${INVENTORY_SERVICE_NAME}"
ORDER_SERVICE_EXEC="${ORDER_SERVICE_NAME}"
INVENTORY_SERVICE_EXEC="${INVENTORY_SERVICE_NAME}"

APP_CHART_PATH="${PROJECT_ROOT}/infrastructure/deployment"
LOKI_VALUES_PATH="${PROJECT_ROOT}/infrastructure/observability/loki-values.yaml"
PROMTAIL_VALUES_PATH="${PROJECT_ROOT}/infrastructure/observability/promtail-values.yaml"
TEMPO_VALUES_PATH="${PROJECT_ROOT}/infrastructure/observability/tempo-values.yaml"
OTEL_COLLECTOR_VALUES_PATH="${PROJECT_ROOT}/infrastructure/observability/otel-collector-values.yaml"
KUBE_PROMETHEUS_STACK_VALUES_PATH="${PROJECT_ROOT}/infrastructure/observability/kube-prometheus-stack-values.yaml"
EVENTEXPORTER_VALUES_PATH="${PROJECT_ROOT}/infrastructure/observability/event-exporter-values.yaml"
PYROSCOPE_VALUES_PATH="${PROJECT_ROOT}/infrastructure/observability/pyroscope-values.yaml"

GRAFANA_PORT=3000
INVENTORY_DASHBOARD="${PROJECT_ROOT}/infrastructure/dashboards/inventory.json"
ORDER_DASHBOARD="${PROJECT_ROOT}/infrastructure/dashboards/order.json"
