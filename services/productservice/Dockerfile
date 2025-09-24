# Stage 1: build
FROM golang:1.24.6-alpine AS builder

# Install build dependencies
WORKDIR /app
RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files from root
COPY ../../go.mod ../../go.sum ./
RUN go mod download

# Copy all code
COPY . .

# Build the binary for productservice
WORKDIR /app/services/productservice
RUN CGO_ENABLED=1 GOOS=linux go build -o productservice ./cmd/main.go

# Stage 2: runtime
FROM alpine:3.18
RUN apk add --no-cache ca-certificates

# Copy built binary from builder
COPY --from=builder /app/services/productservice/productservice /usr/local/bin/productservice

# Expose service port
EXPOSE 8082

# Run service
CMD ["/usr/local/bin/productservice"]
