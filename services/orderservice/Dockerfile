# Stage 1: build
FROM golang:1.24.6 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Copy all source
COPY . .

# Build the binary for orderservice
WORKDIR /app/services/orderservice
RUN CGO_ENABLED=0 GOOS=linux go build -o /orderservice ./cmd/main.go

# Stage 2: runtime
FROM alpine:3.18
RUN apk add --no-cache ca-certificates

# Copy binary from builder
COPY --from=builder /orderservice /usr/local/bin/orderservice

# Expose service port
EXPOSE 8083

# Run service
CMD ["/usr/local/bin/orderservice"]
