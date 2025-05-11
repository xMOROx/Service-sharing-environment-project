.PHONY: proto clean

proto: 
	mkdir -p proto/inventory
	mkdir -p proto/order
	protoc --go_out=. --go-grpc_out=. --proto_path=proto proto/inventory.proto
	protoc --go_out=. --go-grpc_out=. --proto_path=proto proto/order.proto

clean:
	find . -name "*.pb.go" -delete
	rm -rf proto/inventory proto/order
