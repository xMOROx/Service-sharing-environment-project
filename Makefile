.PHONY: all proto clean order inventory deploy images uninstall clean-images

ORDER_DIR=./services/order-service
INVENTORY_DIR=./services/inventory-service
ORDER_EXEC=order-service
INVENTORY_EXEC=inventory-service

all: proto order inventory build-images deploy

build-images:
	docker buildx bake

proto: 
	mkdir -p proto/inventory
	mkdir -p proto/order
	protoc --go_out=.. --go-grpc_out=.. --proto_path=proto proto/inventory.proto
	protoc --go_out=.. --go-grpc_out=.. --proto_path=proto proto/order.proto

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
	helm install demo ./infrastructure/deployment/ -n default

uninstall:
	helm uninstall demo -n default
