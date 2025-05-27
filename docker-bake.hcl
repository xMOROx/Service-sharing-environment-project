group "default" {
  targets = ["order-service", "inventory-service"]
}

target "order-service" {
  context = "."
  dockerfile = "./infrastructure/docker/order-service/Dockerfile"
  tags       = ["localhost:5001/order-service:0.1.0", "order-service:0.1.0"]
}

target "inventory-service" {
  context = "."
  dockerfile = "./infrastructure/docker/inventory-service/Dockerfile"
  tags       = ["localhost:5001/inventory-service:0.1.0", "inventory-service:0.1.0"]
}
