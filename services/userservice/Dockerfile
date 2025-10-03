# Stage 1: build userservice
FROM golang:1.24.6-alpine AS builder
RUN apk add --no-cache gcc musl-dev

WORKDIR /src

# Copy module files from repo root
COPY go.mod go.sum ./
RUN go mod download

# Copy the userservice code + shared pkg folder
COPY services/userservice ./services/userservice
COPY pkg ./pkg

# Build the userservice binary
WORKDIR /src/services/userservice/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o /userservice main.go

# Stage 2: runtime
FROM alpine:3.18
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /userservice /usr/local/bin/userservice

EXPOSE 8081
CMD ["/usr/local/bin/userservice"]
