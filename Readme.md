Userservice → port 8081

Productservice → port 8082

Orderservice → port 8083

Cartservice → port 8084

Each has its own Dockerfile (same pattern).

Then docker-compose.yml can build all services separately.

                           ┌──────────────────────────┐
                           │        Clients           │
                           │ (Web, Mobile, Postman)   │
                           └───────────┬─────────────┘
                                       │
                                HTTPS / API Gateway (Optional)
                                       │
                                ┌──────┴───────┐
                                │   ALB (Load  │
                                │ Balancer)    │
                                └──────┬───────┘
                                       │
                 ┌─────────────────────┼───────────────────────┐
                 │                     │                       │
        ┌────────▼─────────┐  ┌────────▼─────────┐   ┌────────▼─────────┐
        │   User Service   │  │ Product Service  │   │  Order Service   │
        │  (ECS Task/Farg.)│  │  (ECS Task)      │   │  (ECS Task)      │
        │ Port 8081        │  │ Port 8082        │   │ Port 8083        │
        └────────┬─────────┘  └────────┬─────────┘   └────────┬─────────┘
                 │                     │                       │
                 │                     │                       │
                 │                     │                       │
        ┌────────▼─────────┐  ┌────────▼─────────┐   ┌────────▼─────────┐
        │  Cart Service    │  │ Payment Service   │   │ Notification Svc │
        │  (ECS Task)      │  │  (ECS Task)       │   │  (ECS Task)      │
        │  Port 8084       │  │  Port 8085        │   │  Port 8086       │
        └────────┬─────────┘  └────────┬─────────┘   └────────┬─────────┘
                 │                     │                       │
                 └─────────────────────┼───────────────────────┘
                                       │
                             ┌─────────▼─────────┐
                             │   RDS Proxy (Opt) │
                             │ Connection Pooling│
                             └─────────┬─────────┘
                                       │
                             ┌─────────▼─────────┐
                             │  AWS RDS (Postgres)│
                             │  ecommerce DB      │
                             └────────────────────┘


🔹 Sequence Explanation

Client sends order request through API Gateway.

API Gateway checks with User Service that the user is valid.

Gateway forwards the request to Order Service.

Order Service calls Product Service to check stock.

If stock is available, it calls Payment Service.

After payment success, Order Service saves the order in Order DB.

Order Service publishes an OrderCreated event to the Broker.

Inventory Service consumes it and updates stock.

Notification Service consumes it and sends an email/SMS.

Client gets back the confirmation response.

-------------------------------------------------------------------------------------------

🔹 Core Services

User Service → authentication, registration, profile, addresses.

Product Service → product CRUD, details, metadata.

Catalog / Search Service → full-text search, filtering, recommendations (often backed by ElasticSearch).

Cart Service → add/remove/view cart, backed by Redis/DB.

Order Service → create/track orders, orchestrates checkout.

Payment Service → integrates with external payment gateways (Stripe, Razorpay, PayPal).

Inventory Service → stock management, reserve/release items.

Notification Service → emails, SMS, push notifications (confirmation, updates).

🔹 Optional but Common Services

Review & Rating Service → product reviews, ratings, moderation.

Recommendation Service → personalized product suggestions (could be ML based).

Shipping/Logistics Service → shipping options, tracking, integration with courier APIs.

Analytics Service → events, monitoring, sales reports.

Admin/Backoffice Service → for sellers/admins to manage products, orders, discounts.

Promotion/Coupon Service → handle promo codes, discounts, campaigns.

🔹 Infrastructure Components (not “services” but needed)

API Gateway → single entry point (NGINX, Kong, Traefik, AWS API Gateway).

Message Broker / Event Bus → Kafka, RabbitMQ, AWS SQS/SNS for async communication.

Databases → each service has its own DB (Postgres, MySQL, Mongo, Redis, Elastic).

Monitoring/Logging → Prometheus, Grafana, ELK, OpenTelemetry.

🔹 Total Count

Mandatory core services: 8

Optional services (common in real-world): +6
👉 So anywhere from 8 (MVP) to 14+ (full-blown system) depending on how deep you want to go.

✅ If you are following the roadmap.sh Ecommerce API project, you need at least these 8 core services:

User

Product

Catalog/Search

Cart

Order

Payment

Inventory

Notification


My Guidance

Since you’re building step by step, do this roadmap:

Now (Phase 2, MVP) → Build Catalog/Search Service as a separate service.

Keep it cleanly separated from Product Service (don’t merge).

But for speed, fetch data directly from Product DB.

Next (Phase 3/4) → Refactor Catalog/Search to consume Product Service APIs or Product events.

This keeps microservice boundaries clean.

Later (Production scale) → Add ElasticSearch + Event-driven sync for high-performance search.


UserService - API

GET - http://localhost:8081/health

✅ That’s the complete flow for userservice:
health → register → login → fetch profile


┌────────────────────────────────────────────────────────────────────────────┐
│                           API Gateway / ALB (Optional)                     │
│                   ──────────────────────────────────────                   │
│     Handles routing, load balancing, SSL termination, auth (future)        │
│     🟢 Pattern: API Gateway / Reverse Proxy                                 │
└────────────────────────────────────────────────────────────────────────────┘
                                       │
                                       ▼
        ┌────────────────────────────────────────────────────────┐
        │                 Application Services                   │
        └────────────────────────────────────────────────────────┘
                 │              │              │              │
                 │              │              │              │
 ┌───────────────▼────────────┐ ┌──────────────▼────────────┐ ┌──────────────▼────────────┐ ┌──────────────▼────────────┐
 │     User Service           │ │    Product Service         │ │    Order Service          │ │    Cart Service            │
 │  (port 8081 / ECS Task)    │ │ (port 8082 / ECS Task)     │ │ (port 8083 / ECS Task)    │ │ (port 8084 / ECS Task)     │
 │────────────────────────────│ │────────────────────────────│ │────────────────────────────│ │────────────────────────────│
 │ 🧩 Layers inside each svc: │ │                            │ │                            │ │                            │
 │                            │ │                            │ │                            │ │                            │
 │ Handler Layer (Gin)        │ │ Handler Layer (Gin)        │ │ Handler Layer (Gin)        │ │ Handler Layer (Gin)        │
 │  - REST endpoints          │ │  - REST endpoints          │ │  - REST endpoints          │ │  - REST endpoints          │
 │  🟢 Pattern: Controller     │ │  🟢 Pattern: Controller     │ │  🟢 Pattern: Controller     │ │  🟢 Pattern: Controller     │
 │                            │ │                            │ │                            │ │                            │
 │ Service Layer              │ │ Service Layer              │ │ Service Layer              │ │ Service Layer              │
 │  - Business logic          │ │  - Business logic          │ │  - Business logic          │ │  - Aggregates Products +   │
 │  🟢 Pattern: Service/Use-Case│ │  🟢 Pattern: Service/Use-Case│ │  🟢 Pattern: Service/Use-Case│ │    Orders (Aggregator)     │
 │                            │ │                            │ │                            │ │  🟢 Pattern: Aggregator     │
 │ Repository Layer           │ │ Repository Layer           │ │ Repository Layer           │ │ Repository Layer           │
 │  - CRUD / DB ops           │ │  - CRUD / DB ops           │ │  - CRUD / DB ops           │ │  - CRUD / DB ops           │
 │  🟢 Pattern: Repository     │ │  🟢 Pattern: Repository     │ │  🟢 Pattern: Repository     │ │  🟢 Pattern: Repository     │
 │                            │ │                            │ │                            │ │                            │
 │ DB Connection              │ │ DB Connection              │ │ DB Connection              │ │ DB Connection              │
 │  🟢 Pattern: Singleton      │ │  🟢 Pattern: Singleton      │ │  🟢 Pattern: Singleton      │ │  🟢 Pattern: Singleton      │
 │  🟢 Retry & Backoff logic   │ │  🟢 Retry & Backoff logic   │ │  🟢 Retry & Backoff logic   │ │  🟢 Retry & Backoff logic   │
 │                            │ │                            │ │                            │ │                            │
 │ DI (Repo → Service → Hdlr) │ │ DI (Repo → Service → Hdlr) │ │ DI (Repo → Service → Hdlr) │ │ DI (Repo → Service → Hdlr) │
 │  🟢 Pattern: Dependency Inj │ │  🟢 Pattern: Dependency Inj │ │  🟢 Pattern: Dependency Inj │ │  🟢 Pattern: Dependency Inj │
 │                            │ │                            │ │                            │ │                            │
 └───────────────┬────────────┘ └──────────────┬────────────┘ └──────────────┬────────────┘ └──────────────┬────────────┘
                 │                             │                             │                             │
                 │  Internal REST Calls (via Docker DNS / ECS Service DNS)    │
                 │────────────────────────────────────────────────────────────│
                 │
                 ▼
       ┌──────────────────────────────────────────────────────┐
       │                   Shared Packages                    │
       ├──────────────────────────────────────────────────────┤
       │ pkg/config  → 🟢 Config Pattern                      │
       │ pkg/logger  → 🟢 Singleton + Adapter Pattern          │
       │ pkg/jwt     → 🟢 Utility/Strategy (token signing)     │
       │ pkg/db      → 🟢 Retry + Singleton + Factory          │
       │ pkg/middleware → 🟢 Cross-Cutting Concerns Pattern    │
       └──────────────────────────────────────────────────────┘
                 │
                 ▼
       ┌──────────────────────────────────────────────────────┐
       │           AWS RDS (PostgreSQL Database)              │
       │  🟢 Pattern: Shared Persistent Store / Singleton DB   │
       │  Managed via RDS Proxy + Secrets Manager              │
       └──────────────────────────────────────────────────────┘

