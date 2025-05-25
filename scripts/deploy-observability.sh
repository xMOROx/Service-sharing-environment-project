#!/usr/bin/env bash
set -e

PROJECT_ROOT=$(git rev-parse --show-toplevel)
cd "$PROJECT_ROOT" || exit 1

$PROJECT_ROOT/scripts/deploy-prometheus.sh
$PROJECT_ROOT/scripts/deploy-loki.sh
$PROJECT_ROOT/scripts/deploy-eventexporter.sh
$PROJECT_ROOT/scripts/deploy-promtail.sh
$PROJECT_ROOT/scripts/deploy-otlp-collector.sh
