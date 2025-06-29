# Load all variables from .env
ifneq (,$(wildcard .env))
	include .env
	export
endif

MIGRATION_PATH = ./cmd/migrate/migrations


.PHONY: migrate-create
migrate-create:
	@migrate create -seq -ext sql -dir ${MIGRATION_PATH} $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@echo "Running migrations..."
	@migrate -source=file://${MIGRATION_PATH}  -database="${DB_ADDR}" up
	@echo "Migrations completed."

.PHONY: migrate-down
migrate-down:
	@echo "Rolling back migrations..."
	@migrate -source=file://${MIGRATION_PATH}  -database="${DB_ADDR}" down
	@echo "Rollback completed."

.PHONY : gen-docs
gen-docs:
	@echo "Generating API documentation..."
	@swag init -g ./api/main.go -d cmd,internal && swag fmt
	@echo "API documentation generated successfully."