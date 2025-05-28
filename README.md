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
    * `[Unary]` Retrieves detailed static information about a specific product.
    * `[Unary]` Adds a new product definition to the system.
    * `[Unary]` Modifies static details of an existing product.
    * `[Unary]` Marks a product as discontinued or removes its definition.
    * `[Unary]` Directly adjusts the stock quantity for a product (e.g., for returns, manual corrections, after confirmed order fulfillment).
    * `[Client-Streaming]` Processes a sequence of incoming stock items representing a bulk shipment or transfer.
    * `[Unary]` Provides the current stock level and availability status for a single product upon a direct request.
    * `[Server-Streaming]` Streams a list of products, possibly filtered, along with their current stock levels.
    * `[Server-Streaming]` Allows clients to subscribe to and receive ongoing notifications when product stock levels fall below specified thresholds.
    * `[Bidirectional-Streaming]` Engages in a continuous, two-way communication session with a client to:
        * Receive and process a series of requests to provisionally add, update, or remove items for a pending order being built by the client.
        * Send back real-time availability feedback, status of any soft reservations for items in the client's session, and proactive alerts (e.g., if stock for an item in the client's session becomes critically low due to external factors).

2. **Order Service** – Handles order processing:
    * `[Unary]` Checks the availability of individual items with the `Inventory Service` via direct, one-off requests.
    * `[Bidirectional-Streaming]` Initiates and manages an interactive, continuous communication session with the `Inventory Service` to dynamically build an order, check availability for multiple items sequentially, receive live feedback, and potentially request soft reservations.
    * Internally determines whether to confirm or reject an order based on inventory availability feedback and other business logic.
    * If an order is confirmed:
        * `[Bidirectional-Streaming]` Communicates the finalization of a pending order to the `Inventory Service` (often as part of an ongoing order, leading to firming up soft reservations or triggering stock deduction).
        * `[Unary]` Alternatively, sends explicit instructions to the `Inventory Service` to definitively reserve and/or decrement stock for all items in the confirmed order (if not handled via an interactive session).
    * If an order is rejected or an interactive order-building session is cancelled:
        * `[Bidirectional-Streaming]` Notifies the `Inventory Service` to release any soft reservations made for that session.

## 4. Solution architecture

![Solution architecture diagram showing the interaction between services, including the Inventory Service, Order Service, OpenTelemetry Collector, and Grafana.](./images/architecture.svg)

## 5. Environment configuration description

### Linux Setup

To set up the development environment on a Linux system, you will need the following tools:

* **Go**: The programming language used for the services. Installation instructions can be found at [https://golang.org/doc/install](https://golang.org/doc/install).
* **Docker**: For containerizing applications. Install Docker Engine from [https://docs.docker.com/engine/install/](https://docs.docker.com/engine/install/). Ensure Docker Buildx is available (usually included with recent Docker versions).
* **Kind**: For running local Kubernetes clusters. Installation guide: [https://kind.sigs.k8s.io/docs/user/quick-start/#installation](https://kind.sigs.k8s.io/docs/user/quick-start/#installation).
* **Helm**: The package manager for Kubernetes. Installation guide: [https://helm.sh/docs/intro/install/](https://helm.sh/docs/intro/install/).
* **Protobuf Compiler (`protoc`)**: For generating Go code from `.proto` files. Installation instructions can be found at [https://grpc.io/docs/protoc-installation/](https://grpc.io/docs/protoc-installation/).
  * You will also need the Go plugins for protoc: [https://grpc.io/docs/languages/go/quickstart/](https://grpc.io/docs/languages/go/quickstart/)

    Ensure that your `GOBIN` directory (usually `$GOPATH/bin` or `$HOME/go/bin`) is in your system's `PATH`.
* **make**: To use the Makefile for common operations. Usually available by default on most Linux distributions. If not, install it using your distribution's package manager (e.g., `sudo apt install make` on Debian/Ubuntu).
* **bash**: The scripts in this project are written in bash. Ensure it's available (default on most Linux systems).

Once these tools are installed, you can proceed with the installation and deployment steps outlined in this README.

### Windows Setup: WSL2 + Docker Desktop

Setting up a development environment on Windows can be done using **WSL2 (Windows Subsystem for Linux)** in combination with **Docker Desktop**.  
Most steps follow the **Linux Setup** - except for the **Kind** and **Docker** installation.  
Docker Desktop automatically:

- integrates with WSL2 distributions,  
- provides a built-in Kubernetes cluster.

#### Docker Desktop Configuration

1. **Enable WSL integration** for your WSL distribution:  
Open Docker Desktop and go to  `Settings -> Resources -> WSL Integration`, then make sure your desired distro is enabled.

2. **Enable Kubernetes support**:  
Navigate to `Settings -> Kubernetes -> Enable Kubernetes` and enable Kubernetes cluster.

#### Verifying Setup

Use the following commands **inside your WSL terminal (e.g. Ubuntu)**:

```bash
docker version
```

- The output should indicate: `Server: Docker Desktop`.

```bash
kubectl config get-contexts
```

- You should see a context named `docker-desktop`.

```bash
kubectl config current-context
```

- This should return `docker-desktop`.  
  If not, switch using:

```bash
kubectl config use-context docker-desktop
```


## 6. Installation method

## 7. How to reproduce - step by step

### 1. Infrastructure as Code approach

## 8. Demo deployment steps

### 1. Configuration set-up

### 2. Data preparation

### 3. Execution procedure

### 4. Results presentation

## 9. Using AI in the project

## 10. Summary – conclusions

## 11. References

* **OpenTelemetry:** [https://opentelemetry.io/](https://opentelemetry.io/)
* **Go:** [https://golang.org/](https://golang.org/)
* **Docker:** [https://www.docker.com/](https://www.docker.com/)
* **Kubernetes:** [https://kubernetes.io/](https://kubernetes.io/)
* **Kind (Kubernetes in Docker):** [https://kind.sigs.k8s.io/](https://kind.sigs.k8s.io/)
* **Helm:** [https://helm.sh/](https://helm.sh/)
* **gRPC:** [https://grpc.io/](https://grpc.io/)
* **Protocol Buffers:** [https://developers.google.com/protocol-buffers](https://developers.google.com/protocol-buffers)

### Helm Charts Used (from Artifact Hub or other sources)

* **Kube Prometheus Stack:** For Prometheus and Grafana deployment (Version: 72.6.2).
  * Typically found on [Artifact Hub](https://artifacthub.io/packages/helm/prometheus-community/kube-prometheus-stack)
* **Loki:** For log aggregation (Version: 6.30.0).
  * Typically found on [Artifact Hub](https://artifacthub.io/packages/helm/grafana/loki)
* **Tempo:** For distributed tracing (Version: 1.40.2).
  * Typically found on [Artifact Hub](https://artifacthub.io/packages/helm/grafana/tempo)
* **OpenTelemetry Collector:** For collecting telemetry data (Version: 0.125.0).
  * Typically found on [Artifact Hub](https://artifacthub.io/packages/helm/open-telemetry/opentelemetry-collector)
* **Promtail:** For shipping logs to Loki (Version: 6.16.6).
  * Typically found on [Artifact Hub](https://artifacthub.io/packages/helm/grafana/promtail)
* **Kubernetes Event Exporter:** For exporting Kubernetes events (Version: 3.5.3).
  * Typically found on [Artifact Hub](https://artifacthub.io/packages/helm/bitnami/kubernetes-event-exporter)
