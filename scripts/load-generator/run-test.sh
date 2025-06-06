#!/usr/bin/env bash

if ! kubectl get configmap order-proto --namespace demo >/dev/null 2>&1; then
	kubectl create configmap order-proto --from-file=./proto --namespace demo
	echo "ConfigMap 'order-proto' created."
fi

helm install grpc-checkavailability "$(dirname "$0")/../../infrastructure/grpc-load-generator" --namespace demo \
  --set grpc.call=order.OrderService.CheckItemAvailability \
  --set-literal grpc.payload='{"product_id": "P001"}'

helm install grpc-buildorder "$(dirname "$0")/../../infrastructure/grpc-load-generator" --namespace demo \
  --set grpc.call=order.OrderService.BuildOrder \
  --set-literal grpc.payload='[{"product_id": "P002", "requested_quantity": 2}, {"product_id": "P003", "requested_quantity": 3}]' \
  --set grpc.concurrency=5 \
  --set grpc.requests=50
