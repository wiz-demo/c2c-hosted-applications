# Order Processor & Inventory Sync Services

Backend services for order processing and inventory synchronization.

## Components

1. **Inventory Sync Agent (`main.go` / `go.mod`)**
   - Go HTTP/2 and SSH inventory synchronization agent.
   - Build:
     ```bash
     go build -buildvcs=true -o bin/c2c-hosted-applications .
     ```

2. **Order Processor Service (`pom.xml` / `src/main/java/io/wizdemo/c2c/OrderProcessorApplication.java`)**
   - Java 17 / Spring Boot order processing service.
   - Build:
     ```bash
     mvn -B -DskipTests package
     ```

## Health Endpoints

- Inventory Sync Agent: `GET /healthz` (default `127.0.0.1:8081`)
- Order Processor Service: `GET /healthz` (default `127.0.0.1:8082`)
