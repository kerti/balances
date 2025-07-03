include .env
export

SEEDER_DIR := ./migrations/seeder
SEED_FILES := $(shell bash -c "find $(SEEDER_DIR) -type f -name '*.sql' | sort")
MYSQL_CMD := MYSQL_PWD=$(DB_PASS) mysql -h $(DB_HOST) -P $(DB_PORT) -u $(DB_USER) $(DB_NAME)

.DEFAULT_GOAL := help

.PHONY: seed local genmock test test-report check-env help

local:
	@export $$(grep -vE '^\s*#' .env | grep -vE '^\s*$$' | xargs) && go run .

genmock:
	@mockgen -destination=mock/repository/repository.go -package=mock_repository -source=repository/repository.go
	@mockgen -destination=mock/service/service.go -package=mock_service -source=service/service.go

test:
	./coverage.sh

test-report:
	./coverage.sh --html

check-env:
	@[ -n "$(DB_HOST)" ] || (echo "❌ DB_HOST not set" && exit 1)
	@[ -n "$(DB_PORT)" ] || (echo "❌ DB_PORT not set" && exit 1)
	@[ -n "$(DB_USER)" ] || (echo "❌ DB_USER not set" && exit 1)
	@[ -n "$(DB_PASS)" ] || (echo "❌ DB_PASS not set" && exit 1)
	@[ -n "$(DB_NAME)" ] || (echo "❌ DB_NAME not set" && exit 1)

seed: check-env
	@echo "Running seed files in $(SEEDER_DIR)..."
	@for file in $(SEED_FILES); do \
		echo "Executing $$file..."; \
		$(MYSQL_CMD) < $$file || { echo "❌ Failed on $$file"; exit 1; }; \
	done
	@echo "✅ All seeders executed successfully."

help:
	@echo "Available targets:"
	@echo "  local        Run app with .env variables"
	@echo "  genmock      Generate mock files"
	@echo "  test         Run tests"
	@echo "  test-report  Generate HTML test coverage report"
	@echo "  seed         Run DB seeders from $(SEEDER_DIR)"
