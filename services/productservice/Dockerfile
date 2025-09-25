# Stage 1: build
FROM golang:1.24.6-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy go.mod and go.sum from repo root
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code
COPY . .

# Build the binary for productservice
WORKDIR /app/services/productservice
RUN CGO_ENABLED=0 GOOS=linux go build -o /productservice ./cmd/main.go

# Stage 2: runtime
FROM alpine:3.18
RUN apk add --no-cache ca-certificates

COPY --from=builder /productservice /usr/local/bin/productservice

EXPOSE 8082
CMD ["/usr/local/bin/productservice"]
