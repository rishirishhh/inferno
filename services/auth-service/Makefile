include ../../.env

DB_URL := postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/auth_db?sslmode=disable

auth-migrate-up:
	migrate -path ./internal/db/migrations \
	-database "$(DB_URL)" \
	up

auth-migrate-down:
	migrate -path ./internal/db/migrations \
	-database "$(DB_URL)" \
	down 1

auth-sqlc:
	sqlc generate