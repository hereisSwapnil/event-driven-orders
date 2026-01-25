# Event Driven Orders - Low Level Design (LLD)

## Project Overview
Event Driven Orders is a scalable, asynchronous order processing system implemented in Go. It demonstrates core distributed system concepts including Event-Driven Architecture, delayed task scheduling using Redis, and reliable messaging with Kafka.

## Architecture
The project follows **Clean Architecture** principles to ensure separation of concerns:

[![mermaid-diagram-2026-01-25-160716.png](https://i.postimg.cc/8C40VGtP/mermaid-diagram-2026-01-25-160716.png)](https://postimg.cc/hJJr9HkH)

- **Domain Layer** (`internal/order/domain`): Contains core entities (`Order`, `OrderStatus`) and business rules. It has no external dependencies.
- **Service Layer** (`internal/order/service`): Orchestrates business use cases. Handles the flow between repositories and event producers.
- **Repository Layer** (`internal/order/repository`): Implements data access interfaces for PostgreSQL and Redis.
- **Transport Layer**:
  - **HTTP**: Handles REST requests (`cmd/order-service`).
  - **Kafka Consumers**: Handles asynchronous event processing (`cmd/scheduler-worker`, `cmd/processing-service`).

## Low Level Design (LLD)

### 1. Entities & Data Structures
**Order Entity**
- **ID**: UUID, unique identifier.
- **Status**: State tracking (`CREATED`, `PROCESSING`, `COMPLETED`, `FAILED`).
- **ScheduledAt**: Optional timestamp for delayed execution.
- **CreatedAt**: Timestamp of ingestion.

**Scheduler Store (Redis)**
- Uses **Sorted Sets (ZSET)** to manage time-based priority.
- **Score**: Unix timestamp of execution time.
- **Member**: Order ID.

### 2. Async & Event Model
The system is designed to handle orders asynchronously and decouple ingestion from processing:

- **Ingestion**: The API saves the order to Postgres and keeps it in `CREATED` state. It produces an `order.created` event.
- **Scheduling**: The **Scheduler Worker** consumes `order.created`.
  - If `ScheduledAt` is present, it adds the order to Redis with that timestamp.
  - If `ScheduledAt` is null (immediate), it uses `Now()` as the timestamp.
- **Dispatch Loop**: A background ticker (`1s`) polls Redis for items where `score <= Now`. Due items are moved to the `order.ready` topic.
- **Processing**: The **Processing Service** consumes `order.ready` and updates the DB state to `COMPLETED`.

### 3. API Specification
| Method | Endpoint | Body | Description |
|--------|----------|------|-------------|
| `POST` | `/orders` | `{"customer_name": "...", "total_price": 100, "scheduled_time": "..."}` | Creates a new order. |
| `GET` | `/orders` | `?id={uuid}` | Returns order details and status. |

### 4. State Machine (Order Status)
A state machine governs the lifecycle of an order:

- **CREATED**: Initial state after POST request.
- **PROCESSING**: Order picked up by the processing service (from `order.ready`).
- **COMPLETED**: Order logic successfully executed.

Flow: `CREATED` $\rightarrow$ [Scheduler] $\rightarrow$ [Processor] $\rightarrow$ `PROCESSING` $\rightarrow$ `COMPLETED`

## Setup & Running

### Prerequisites
- Docker & Docker Compose
- Go 1.23+ (optional, if running locally without Docker)

### Installation
1. **Clone the repository**:
   ```bash
   git clone https://github.com/hereisSwapnil/event-driven-orders.git
   cd event-driven-orders
   ```

2. **Run the system**:
   ```bash
   docker-compose up -d --build
   ```
   *This starts Postgres, Redis, Kafka (KRaft), and all 3 microservices.*

3. **Verify Status**:
   ```bash
   docker-compose ps
   ```

### Usage
**Immediate Order**:
```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{"id": "order-1", "customer_name": "Swapnil", "total_price": 500}'
```

**Scheduled Order**:
```bash
# Schedule for 10 seconds later
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "id": "order-future", 
    "customer_name": "Future Swapnil", 
    "total_price": 1000,
    "scheduled_time": "'$(date -v+10S -u +"%Y-%m-%dT%H:%M:%SZ")'"
  }'
```

## Technologies
- **Language**: Go (Golang) 1.23
- **Messaging**: Apache Kafka (KRaft)
- **Database**: PostgreSQL 15
- **Cache**: Redis 7
- **Containerization**: Docker
