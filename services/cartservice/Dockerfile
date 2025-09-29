FROM golang:1.24.6 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
WORKDIR /app/services/cartservice
RUN CGO_ENABLED=0 GOOS=linux go build -o /cartservice ./cmd/main.go

FROM alpine:3.18
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /cartservice /usr/local/bin/cartservice

EXPOSE 8084
CMD ["/usr/local/bin/cartservice"]
