Order Service — Responsibilities

Create new orders (user buys products).

Fetch order details.

List all orders for a user.

(Later) Handle status updates (pending, paid, shipped, delivered, cancelled).

(Later) Integrate with Payment Service.

id          UUID (primary key)
user_id     UUID (references users.id)
status      TEXT (e.g., "pending", "completed")
total_price NUMERIC(10,2)
created_at  TIMESTAMP
updated_at  TIMESTAMP


id          UUID (primary key)
order_id    UUID (references orders.id)
product_id  UUID (references products.id)
quantity    INT
price       NUMERIC(10,2) -- price at purchase time
