.PHONY: all proto clean order inventory deploy images uninstall uninstall-observability clean-images forward-ports clean-ports deploy-observability create-cluster remove-cluster help 

ORDER_DIR=./services/order-service
INVENTORY_DIR=./services/inventory-service
ORDER_EXEC=order-service
INVENTORY_EXEC=inventory-service

# ============== Main Targets ==============
all: deploy

# ============== Build Targets ==============
build-images:
	docker buildx bake

proto:
	./scripts/generate-proto.sh

order:
	cd ${ORDER_DIR} && go build .

inventory:
	cd ${INVENTORY_DIR} && go build .

# ============== Deployment Targets ==============
deploy:
	./scripts/service-management/deploy.sh

deploy-observability:
	./scripts/observability/deploy-observability.sh # Ensure this script is updated to call sub-scripts from scripts/observability/

# ============== Cluster Management ==============
create-cluster:
	./scripts/cluster/create-cluster.sh

remove-cluster:
	./scripts/cluster/remove-cluster.sh

# ============== Port Management ==============
forward-ports:
	./scripts/service-management/forward-ports.sh

clean-ports:
	./scripts/service-management/clean-ports.sh

# ============== Cleaning Targets ==============
uninstall: clean-ports
	./scripts/service-management/uninstall.sh

uninstall-observability:
	./scripts/observability/uninstall-observability.sh

clean: clean-images
	find . -name "*.pb.go" -delete
	rm -rf proto/inventory proto/order

clean-images:
	rm -f ${ORDER_DIR}/${ORDER_EXEC}
	rm -f ${INVENTORY_DIR}/${INVENTORY_EXEC}

# ============== Load generator ==============
test:
	./scripts/load-generator/run-test.sh

test-clean:
	./scripts/load-generator/clean.sh

# ============== Help Target ==============
help:
	@echo "Available targets:"
	@echo "  all                   - Deploy everything (default)"
	@echo "  build-images          - Build docker images using docker-bake"
	@echo "  proto                 - Generate protobuf files"
	@echo "  order                 - Build the order service"
	@echo "  inventory             - Build the inventory service"
	@echo "  deploy                - Deploy services"
	@echo "  deploy-observability  - Deploy observability stack"
	@echo "  uninstall             - Uninstall services and clean ports"
	@echo "  uninstall-observability - Uninstall observability stack"
	@echo "  forward-ports         - Forward ports for services"
	@echo "  clean-ports           - Clean forwarded ports"
	@echo "  create-cluster        - Create local Kind cluster"
	@echo "  remove-cluster        - Remove local Kind cluster"
	@echo "  clean                 - Clean generated files and images"
	@echo "  clean-images          - Clean built executables"
	@echo "  help                  - Show this help message"

