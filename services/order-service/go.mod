module Service-sharing-environment-project/services/order-service

go 1.24.2

require (
	google.golang.org/grpc v1.72.0
	google.golang.org/protobuf v1.36.6
	example.com/inventory-service v0.0.0
)

require (
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250218202821-56aae31c358a // indirect
)
replace example.com/inventory-service => ../inventory-service