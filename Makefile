.PHONY: all proto clean order inventory deploy images uninstall clean-images forward-ports clean-ports deploy-observability create-cluster remove-cluster

ORDER_DIR=./services/order-service
INVENTORY_DIR=./services/inventory-service
ORDER_EXEC=order-service
INVENTORY_EXEC=inventory-service

all: deploy

build-images:
	docker buildx bake

proto:
	./scripts/generate-proto.sh

clean: clean-images
	find . -name "*.pb.go" -delete
	rm -rf proto/inventory proto/order

clean-images:
	rm -f ${ORDER_DIR}/${ORDER_EXEC}
	rm -f ${INVENTORY_DIR}/${INVENTORY_EXEC}

order:
	cd ${ORDER_DIR} && go build .

inventory:
	cd ${INVENTORY_DIR} && go build .

deploy:
	./scripts/deploy.sh

deploy-observability:
	./scripts/deploy-observability.sh

uninstall: clean-ports
	./scripts/uninstall.sh

forward-ports:
	./scripts/forward-ports.sh

clean-ports:
	./scripts/clean-ports.sh

create-cluster:
	./scripts/create-cluster.sh

remove-cluster:
	./scripts/remove-cluster.sh

