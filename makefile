.PHONY: clean tidy run all migrate-up migrate-down compose-up compose-down

clean:
	go clean -modcache

tidy:
	go mod tidy

run:
	go run ./cmd/main.go

# Migrate up (apply all .up.sql)
run-up:
	go run ./cmd/main.go up

# Rollback last N migrations (default 1 if not specified)
# Usage: make run-down steps=2
run-down:
	go run ./cmd/main.go down $(steps)


all: clean tidy run

# Migration commands
migrate-up: 
	@echo "Running database migration UP..."
	@chmod +x scripts/run_migration.sh
	@./scripts/run_migration.sh up

migrate-down: 
	@echo "Running database migration DOWN..."
	@chmod +x scripts/run_migration.sh
	@./scripts/run_migration.sh down

# Setup database and run migrations
setup-db: migrate-up
	@echo "Database setup completed!"

# Reset database (drop and recreate tables)
reset-db: migrate-down migrate-up
	@echo "Database reset completed!"

compose-up:
	@docker-compose up -d
	@echo "Database setup completed!"

compose-down:
	@docker-compose down -d
	@echo "Database reset completed!"