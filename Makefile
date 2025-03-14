
# Database Configuration (Default values, can be overridden)
POSTGRES_SQL_HOST ?= localhost
POSTGRES_SQL_PORT ?= 5432
POSTGRES_SQL_USER ?= postgres
POSTGRES_SQL_PASSWORD ?= postgres
POSTGRES_SQL_DATABASE_NAME ?= postgres

MIGRATIONS_DIR = src/internal/database/migrations
SCHEMA_FILE = src/internal/database/schema/schema.sql

all: help

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  generate_migration_files  Generate a new migration file"
	@echo "  migrate_up                Apply all pending migrations"
	@echo "  migrate_down              Rollback the last migration"
	@echo "  gen_sqlc                  Generate SQLC Go code"
	@echo "  reset_database            Reset the database"
	@echo "  build_release             Build the release binary"
	@echo "  prod                      Run the application in production mode"
	@echo "  dev                       Run the application in development mode"

# Check if go-migrate is installed, otherwise install it
install_go_migration_if_not_installed:
	@command -v migrate >/dev/null 2>&1 || go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Generate a new migration file (manually triggered)
generate_migration_files: install_go_migration_if_not_installed
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq init_schema

# Apply all pending migrations
migrate_up: install_go_migration_if_not_installed
	migrate -database "postgres://$(POSTGRES_SQL_USER):$(POSTGRES_SQL_PASSWORD)@$(POSTGRES_SQL_HOST):$(POSTGRES_SQL_PORT)/$(POSTGRES_SQL_DATABASE_NAME)?sslmode=disable" -path $(MIGRATIONS_DIR) up

# Rollback the last migration
migrate_down: install_go_migration_if_not_installed
	migrate -database "postgres://$(POSTGRES_SQL_USER):$(POSTGRES_SQL_PASSWORD)@$(POSTGRES_SQL_HOST):$(POSTGRES_SQL_PORT)/$(POSTGRES_SQL_DATABASE_NAME)?sslmode=disable" -path $(MIGRATIONS_DIR) down 1

# Check if sqlc is installed, otherwise install it
install_sqlc_if_not_installed:
	@command -v sqlc >/dev/null 2>&1 || go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Generate SQLC Go code (ensures schema is applied and migrations are up before running)
gen_sqlc: install_sqlc_if_not_installed migrate_up
	cd src && sqlc generate

reset_database:
	PGPASSWORD=$(POSTGRES_SQL_PASSWORD) psql -U $(POSTGRES_SQL_USER) -h $(POSTGRES_SQL_HOST) -d template1 -c "DROP DATABASE IF EXISTS $(POSTGRES_SQL_DATABASE_NAME);"
	PGPASSWORD=$(POSTGRES_SQL_PASSWORD) psql -U $(POSTGRES_SQL_USER) -h $(POSTGRES_SQL_HOST) -d template1 -c "CREATE DATABASE $(POSTGRES_SQL_DATABASE_NAME);"

build_release:
	cd src && go build -ldflags "-s -w" -o dist/release .


install_air_if_not_installed:
	@command -v air >/dev/null 2>&1 || go install github.com/air-verse/air@latest


prod:
	docker-compose -f docker-compose.yaml up --build -d
	@echo "Application is running in production mode"
	@echo "Before request need do migration: make migrate_up"
	@echo "check the environment variables in the docker-compose.yaml file"
	@echo "for the database connection"

	$(MAKE) migrate_up

dev: install_air_if_not_installed
	cd src && air