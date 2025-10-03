# Stage 1: build
FROM golang:1.24.6 AS builder

WORKDIR /app
COPY go.mod go.sum ./ 
RUN go mod download

# Copy all source
COPY . .

# Build the binary for catalogservice
WORKDIR /app/services/searchservice
RUN CGO_ENABLED=0 GOOS=linux go build -o /searchservice ./cmd/main.go

# Stage 2: runtime
FROM alpine:3.18
RUN apk add --no-cache ca-certificates

# Copy binary from builder
COPY --from=builder /searchservice /usr/local/bin/searchservice

# Expose service port
EXPOSE 8084

# Run service
CMD ["/usr/local/bin/searchservice"]
