group "default" {
  targets = ["order-service", "inventory-service"]
}

target "order-service" {
  context    = "./docker/order-service"
  dockerfile = "Dockerfile"
  tags       = ["order-service:0.1.0"]
}

target "inventory-service" {
  context    = "./docker/inventory-service"
  dockerfile = "Dockerfile"
  tags       = ["inventory-service:0.1.0"]
}
