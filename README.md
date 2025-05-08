# Documentation
## 0. Division of labor
| Mateusz Bywalec | Szymon Głomski | Dawid Wołek | Patryk Zajdel |
|-|-|-|-|
| go related work | go related work | grafana | k8s + docker |


## 1. Introduction
The goal of this project is to gain practical experience with OpenTelemetry by building a simple client-server application in Go. A lightweight implementation of both the server and client will be created, followed by the integration of tracing, metrics, and logging capabilities using the OpenTelemetry SDK. The collected data will be sent to an OpenTelemetry Collector and then visualized in Grafana for monitoring and analysis. The aim is to achieve a minimal yet functional observability setup for the application.


## 2. Theoretical background/technology stack
| Deployment | Observability   | App                              | Testing            |
| ---------- | --------------- | -------------------------------- | ------------------ |
| k8s        |  OTel Collector | Go                               | \<load generator\> |
| Docker     |  Grafana        | gRPC + protobuf                  |                    |
| Helm       |  Prometheus     | otel-go + pprof                  |                    |

## 3. Case study concept description
### Application Idea: “Inventory & Order Management”
**Domain:**  
A simple system for managing orders and inventory levels.

### Services

1. **Inventory Service** – Manages inventory levels:
    *   `[Unary]` Retrieves detailed static information about a specific product.
    *   `[Unary]` Adds a new product definition to the system.
    *   `[Unary]` Modifies static details of an existing product.
    *   `[Unary]` Marks a product as discontinued or removes its definition.
    *   `[Unary]` Directly adjusts the stock quantity for a product (e.g., for returns, manual corrections, after confirmed order fulfillment).
    *   `[Client-Streaming]` Processes a sequence of incoming stock items representing a bulk shipment or transfer.
    *   `[Unary]` Provides the current stock level and availability status for a single product upon a direct request.
    *   `[Server-Streaming]` Streams a list of products, possibly filtered, along with their current stock levels.
    *   `[Server-Streaming]` Allows clients to subscribe to and receive ongoing notifications when product stock levels fall below specified thresholds.
    *   `[Bidirectional-Streaming]` Engages in a continuous, two-way communication session with a client to:
        *   Receive and process a series of requests to provisionally add, update, or remove items for a pending order being built by the client.
        *   Send back real-time availability feedback, status of any soft reservations for items in the client's session, and proactive alerts (e.g., if stock for an item in the client's session becomes critically low due to external factors).

2.  **Order Service** – Handles order processing:
    *   `[Unary]` Checks the availability of individual items with the `Inventory Service` via direct, one-off requests.
    *   `[Bidirectional-Streaming]` Initiates and manages an interactive, continuous communication session with the `Inventory Service` to dynamically build an order, check availability for multiple items sequentially, receive live feedback, and potentially request soft reservations.
    *   Internally determines whether to confirm or reject an order based on inventory availability feedback and other business logic.
    *   If an order is confirmed:
        *   `[Bidirectional-Streaming]` Communicates the finalization of a pending order to the `Inventory Service` (often as part of an ongoing order, leading to firming up soft reservations or triggering stock deduction).
        *   `[Unary]` Alternatively, sends explicit instructions to the `Inventory Service` to definitively reserve and/or decrement stock for all items in the confirmed order (if not handled via an interactive session).
    *   If an order is rejected or an interactive order-building session is cancelled:
        *   `[Bidirectional-Streaming]` Notifies the `Inventory Service` to release any soft reservations made for that session.

## 4. Solution architecture
mermaid
## 5. Environment configuration description
## 6. Installation method
## 7. How to reproduce - step by step
### 1. Infrastructure as Code approach
## 8. Demo deployment steps:
### 1. Configuration set-up
### 2. Data preparation
### 3. Execution procedure
### 4. Results presentation
## 9. Using AI in the project
## 10. Summary – conclusions
## 11. References
