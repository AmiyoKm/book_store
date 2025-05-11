MIGRATION_PATH = ./cmd/migrate/migrations


.PHONY: migrate-create
migrate-create:
	@migrate create -seq -ext sql -dir ${MIGRATION_PATH} $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@echo "Running migrations..."
	@migrate -source=file://${MIGRATION_PATH}  -database="postgres://admin:adminpassword@localhost:5432/book_store?sslmode=disable" up
	@echo "Migrations completed."

.PHONY: migrate-down
migrate-down:
	@echo "Rolling back migrations..."
	@migrate -source=file://${MIGRATION_PATH}  -database=postgres://admin:adminpassword@localhost:5432/book_store?sslmode=disable down
	@echo "Rollback completed."