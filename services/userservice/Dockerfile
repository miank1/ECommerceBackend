# Stage 1: build
FROM golang:1.24.6 AS builder

WORKDIR /app

# Copy dependency files first (better cache)
COPY go.mod go.sum ./
RUN go mod download

# Copy everything
COPY . . 

# Build only userservice
WORKDIR /app/services/userservice
RUN CGO_ENABLED=0 GOOS=linux go build -o /userservice ./cmd/main.go

# Stage 2: runtime
FROM alpine:3.18
RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /userservice /usr/local/bin/userservice

EXPOSE 8081
CMD ["/usr/local/bin/userservice"]
