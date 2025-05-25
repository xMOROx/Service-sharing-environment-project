#!/usr/bin/env bash
set -e

./deploy-prometheus.sh
./deploy-loki.sh
./deploy-eventexporter.sh
./deploy-promtail.sh
