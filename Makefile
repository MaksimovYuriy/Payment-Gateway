DB_URL=postgres://payment_gateway:payment_gateway@localhost:5434/payment_gateway?sslmode=disable
MIGRATIONS_DIR=./migrations

build-up:
	docker compose up -d --build 

migrate-up:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DB_URL) GOOSE_MIGRATION_DIR=$(MIGRATIONS_DIR) goose up

migrate-down:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DB_URL) GOOSE_MIGRATION_DIR=$(MIGRATIONS_DIR) goose down

migrate-status:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DB_URL) GOOSE_MIGRATION_DIR=$(MIGRATIONS_DIR) goose status
