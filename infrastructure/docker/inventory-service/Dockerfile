FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY ./services/inventory-service/go.mod ./services/inventory-service/go.sum ./services/inventory-service/
COPY ./proto ./proto

COPY . .

WORKDIR /app/services/inventory-service
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/inventory-service-app ./main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/inventory-service-app /app/inventory-service-app
EXPOSE 50051
CMD ["/app/inventory-service-app"]
