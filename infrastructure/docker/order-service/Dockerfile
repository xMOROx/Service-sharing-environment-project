FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY ./services/order-service/go.mod ./services/order-service/go.sum ./services/order-service/
COPY ./proto ./proto

COPY . .

WORKDIR /app/services/order-service
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/order-service-app ./main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/order-service-app /app/order-service-app
EXPOSE 50052
CMD ["/app/order-service-app"]
