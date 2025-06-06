#!/usr/bin/env bash
uninstall_if_exists() {
  local release_name="$1"
  if helm ls -n "demo" | grep -q "^$release_name"; then
    echo "Uninstalling existing Helm release '$release_name'..."
    helm uninstall "$release_name" -n "demo"
  fi
}

if ! kubectl get configmap order-proto --namespace demo >/dev/null 2>&1; then
	kubectl create configmap order-proto --from-file=./proto --namespace demo
	echo "ConfigMap 'order-proto' created."
fi

uninstall_if_exists grpc-checkavailability
helm install grpc-checkavailability "$(dirname "$0")/../../infrastructure/grpc-load-generator" --namespace demo \
  --set grpc.call=order.OrderService.CheckItemAvailability \
  --set-literal grpc.payload='{"product_id": "P001"}'

uninstall_if_exists grpc-buildorder
helm install grpc-buildorder "$(dirname "$0")/../../infrastructure/grpc-load-generator" --namespace demo \
  --set grpc.call=order.OrderService.BuildOrder \
  --set-literal grpc.payload='[{"product_id": "P002", "requested_quantity": 2}, {"product_id": "P003", "requested_quantity": 3}]' \
  --set grpc.concurrency=5 \
  --set grpc.requests=50
