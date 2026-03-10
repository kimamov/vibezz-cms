.PHONY: dev dev-api dev-admin dev-web db db-up db-down migrate sqlc build test

dev: db-up dev-api dev-admin dev-web

dev-api:
	cd apps/api && go run ./cmd/server

dev-admin:
	cd apps/admin && npx nuxi dev --port 3001

dev-web:
	cd apps/web && npx nuxi dev --port 3000

db:
	docker compose up -d db

db-up: db
	@echo "Waiting for PostgreSQL..."
	@sleep 2
	@echo "PostgreSQL is ready."

db-down:
	docker compose down

migrate:
	cd apps/api && go run ./cmd/server migrate

sqlc:
	cd apps/api && sqlc generate

build:
	cd apps/api && go build -o ../../bin/api ./cmd/server

test:
	cd apps/api && go test ./...
