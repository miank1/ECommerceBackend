# E-Commerce Backend (Microservices Architecture)

A scalable E-Commerce backend built using Go (Golang), PostgreSQL, JWT Authentication, Docker, RabbitMQ, and Microservices Architecture.

## Services

Each service is maintained in its own repository and can be developed, deployed, and scaled independently.

| Service              | Port    | Repository                                | Responsibility                             |
| -------------------- | ------- | ----------------------------------------- | ------------------------------------------ |
| User Service         | 8081    | https://github.com/miank1/user_service    | Authentication, Registration, User Profile |
| Product Service      | 8082    | https://github.com/miank1/product_service | Product Management, Inventory Metadata     |
| Order Service        | 8083    | https://github.com/miank1/order_service   | Order Lifecycle Management                 |
| Cart Service         | 8084    | https://github.com/miank1/cart_service    | Shopping Cart Operations                   |
| Payment Service      | 8085    | https://github.com/miank1/payment_service | Payment Processing                         |
| Search Service       | Planned | Coming Soon                               | Catalog Search & Recommendations           |
| Notification Service | Planned | Coming Soon                               | Email, SMS & Push Notifications            |

---

## Architecture

```text
Clients
   │
   ▼
API Gateway / Load Balancer
   │
   ├───────────────┬───────────────┬───────────────┐
   ▼               ▼               ▼               ▼

User Service   Product Service  Cart Service  Order Service
                                        │
                                        ▼
                                  RabbitMQ
                                        │
                                        ▼
                                 Payment Service
                                        │
                                        ▼
                              Notification Service
```

---

## Technology Stack

### Backend

* Golang
* Gin Framework
* PostgreSQL
* GORM

### Authentication

* JWT Authentication
* Middleware-based Authorization

### Messaging

* RabbitMQ
* Event Driven Communication

### Infrastructure

* Docker
* Docker Compose
* AWS ECS / Fargate (Planned)
* AWS RDS PostgreSQL (Planned)

### Monitoring

* Prometheus (Planned)
* Grafana (Planned)

---

## Design Patterns Used

### Service Layer Pattern

Business logic separated from transport layer.

### Repository Pattern

Database access abstracted behind repositories.

### Dependency Injection

Repository → Service → Handler dependency flow.

### Middleware Pattern

Authentication, Logging, and Cross-Cutting Concerns.

### Event Driven Architecture

RabbitMQ based asynchronous service communication.

---

## Current Features

### User Service

* User Registration
* User Login
* JWT Generation
* Profile APIs

### Product Service

* Product CRUD
* Inventory Management

### Cart Service

* Add to Cart
* Update Cart
* Remove Items
* Checkout

### Order Service

* Create Order
* Track Order Status
* Update Order Status

### Payment Service

* Create Payment
* Update Payment Status
* Payment Tracking

---

## Roadmap

### Phase 1 (Completed)

* User Service
* Product Service
* Cart Service
* Order Service
* Payment Service
* JWT Authentication

### Phase 2 (In Progress)

* RabbitMQ Integration
* Event Driven Checkout Flow

### Phase 3

* Search Service
* Notification Service

### Phase 4

* Docker Compose
* CI/CD Pipeline
* AWS Deployment

### Phase 5

* Kafka
* Distributed Tracing
* Observability Stack
* Production Scale Architecture

```
```

██████████████████░░░░░░░░░░ 75%

✅ Core Microservices
✅ Authentication
✅ REST APIs
✅ PostgreSQL
✅ Docker
✅ Service-to-Service Communication
✅ Performance Basics
🚧 gRPC
✅ Kafka
❌ Kubernetes
❌ Observability
❌ Production Hardening

Create Product
      ↓
Save product_id
      ↓
Register User
      ↓
Save user_id
      ↓
Login
      ↓
Save JWT
      ↓
Add Item
      ↓
Checkout
      ↓
Order
      ↓
Kafka
      ↓
Payment
      ↓
Reduce Stock