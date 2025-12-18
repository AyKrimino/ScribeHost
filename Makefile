BINARY_NAME=scribehost
BINARY_DIR=bin
BINARY_PATH=$(BINARY_DIR)/$(BINARY_NAME)
MIGRATIONS_DIR=migrations
LOCAL_ENV=.env.local
PROD_ENV=.env.production

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BINARY_DIR)
	go build -o $(BINARY_PATH) server.go
	@echo "Build complete."

run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_PATH)

dev:
	@echo "Starting development server with Air..."
	@which air > /dev/null || (echo "Error: 'air' not found. Install it with 'go install github.com/air-verse/air@latest'"; exit 1)
	air

# Usage: make migrate-create name=your_migration_name
migrate-create:
	@echo "Creating migration: $(name)"
	@goose -dir $(MIGRATIONS_DIR) create $(name) sql

migrate-up:
	@echo "Applying migrations up..."
	@goose -env $(LOCAL_ENV) -dir $(MIGRATIONS_DIR) up

migrate-down:
	@echo "Rolling back last migration..."
	@goose -env $(LOCAL_ENV) -dir $(MIGRATIONS_DIR) down

clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_PATH)
	@echo "Clean complete."

# Docker Compose Commands
up:
	@echo "Starting services with Docker Compose..."
	docker compose up -d
	@echo "Services started. Use 'make logs' to view logs."

down:
	@echo "Stopping and removing services..."
	docker compose down
	@echo "Services stopped and removed."

logs:
	@echo "Viewing Docker Compose logs (Ctrl+C to stop)..."
	docker compose logs -f

help:
	@echo "Available commands:"
	@echo "  make build            - Build the application"
	@echo "  make run              - Build and run the application"
	@echo "  make dev              - Run with Air (live reload)"
	@echo "  make migrate-create name=<name> - Create new migration"
	@echo "  make migrate-up       - Apply up migrations"
	@echo "  make migrate-down     - Rollback last migration"
	@echo "  make clean            - Remove binary"
	@echo "  make help             - Show this help"

.PHONY: build run dev migrate-create migrate-up migrate-down clean help up down logs
