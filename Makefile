.PHONY: run build up down migrate swag test help

# Variables
DB_CONTAINER=checkin_db
DB_USER=user
DB_NAME=checkin_db

run:
	go run cmd/server/main.go

build:
	go build -o bin/server cmd/server/main.go

up:
	docker-compose up -d

down:
	docker-compose down

migrate:
	docker exec -i $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME) < migrations/001_initial_schema.up.sql

swag:
	swag init -g cmd/server/main.go

test:
	go test ./...

help:
	@echo "Available commands:"
	@echo "  make run       - Run the server"
	@echo "  make build     - Build the server binary"
	@echo "  make up        - Start Docker containers"
	@echo "  make down      - Stop Docker containers"
	@echo "  make migrate   - Apply database migrations"
	@echo "  make swag      - Generate Swagger documentation"
	@echo "  make test      - Run tests"
