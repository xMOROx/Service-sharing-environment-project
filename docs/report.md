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
   - Adds/removes products  
   - Provides inventory status upon request

2. **Order Service** - Handles order processing:
   - Accepts orders  
   - Checks availability with the `Inventory Service`
   - Confirms or rejects the order

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
