.PHONY: proto clean

proto: inventory-service order-service

inventory-service:
	cd services/inventory-service && protoc --go_out=. --go-grpc_out=. --proto_path=../../proto ../../proto/inventory.proto

order-service:
	cd services/order-service && protoc --go_out=. --go-grpc_out=. --proto_path=../../proto ../../proto/order.proto

clean:
	find . -name "*.pb.go" -delete
