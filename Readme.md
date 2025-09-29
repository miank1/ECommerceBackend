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
