.PHONY: up down restart build logs ps clean

up:
	docker compose -f deployments/docker-compose.yml up --build

down:
	docker compose -f deployments/docker-compose.yml down

restart:
	docker compose -f deployments/docker-compose.yml down
	docker compose -f deployments/docker-compose.yml up --build

build:
	docker compose -f deployments/docker-compose.yml build

logs:
	docker compose -f deployments/docker-compose.yml logs -f

ps:
	docker compose -f deployments/docker-compose.yml ps

clean:
	docker compose -f deployments/docker-compose.yml down -v
