# Makefile for E-Commerce Microservices

# Variables
REGISTRY = 221095782431.dkr.ecr.ap-south-1.amazonaws.com
SERVICES = userservice productservice orderservice
DOCKER_COMPOSE = docker-compose.yaml

# Default Go build settings
GOOS = linux
GOARCH = amd64

.PHONY: all build docker-build docker-push run clean

# Build all services
all: build

build:
	@for service in $(SERVICES); do \
		echo "🚀 Building $$service..."; \
		cd services/$$service && \
		GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o bin/$$service ./cmd/main.go; \
		cd - >/dev/null; \
	done

# Docker build all services
docker-build:
	@for service in $(SERVICES); do \
		echo "🐳 Building Docker image for $$service..."; \
		docker build -t $$service -f services/$$service/Dockerfile .; \
	done

# Tag & Push to ECR
docker-push:
	@for service in $(SERVICES); do \
		echo "📤 Pushing $$service to ECR..."; \
		docker tag $$service:latest $(REGISTRY)/$$service:latest; \
		docker push $(REGISTRY)/$$service:latest; \
	done

# Run locally with Docker Compose
run:
	docker compose -f $(DOCKER_COMPOSE) up --build

# Clean build artifacts
clean:
	@for service in $(SERVICES); do \
		echo "🧹 Cleaning $$service..."; \
		rm -rf services/$$service/bin; \
	done
	docker system prune -f
