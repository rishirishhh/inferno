.PHONY: up down restart logs clean ps

up:
	docker compose up -d

build:
	docker compose up -d --build

down:
	docker compose down -v

restart:
	docker compose restart

logs:
	docker compose logs -f

ps:
	docker ps

clean:
	docker compose down -v

postgres:
	docker exec -it inferno-postgres psql -U admin -d inferno

redis:
	docker exec -it redis redis-cli

kafka-logs:
	docker compose logs -f kafka

prometheus-logs:
	docker compose logs -f prometheus

grafana-logs:
	docker compose logs -f grafana