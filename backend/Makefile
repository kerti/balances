include .env
export

BOLD := \033[1m
GREEN := \033[0;32m
RED := \033[0;31m
NC := \033[0m

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
	@[ -n "$(DB_HOST)" ] || (echo "❌ $(RED)DB_HOST$(NC) not set" && exit 1)
	@[ -n "$(DB_PORT)" ] || (echo "❌ $(RED)DB_PORT$(NC) not set" && exit 1)
	@[ -n "$(DB_USER)" ] || (echo "❌ $(RED)DB_USER$(NC) not set" && exit 1)
	@[ -n "$(DB_PASS)" ] || (echo "❌ $(RED)DB_PASS$(NC) not set" && exit 1)
	@[ -n "$(DB_NAME)" ] || (echo "❌ $(RED)DB_NAME$(NC) not set" && exit 1)

seed: check-env
	@echo "$(BOLD)Running seed files in $(GREEN)$(SEEDER_DIR)$(NC)..."
	@for file in $(SEED_FILES); do \
		printf "==> Executing $(GREEN)$$file$(NC)... "; \
		if $(MYSQL_CMD) < $$file; then \
			printf "$(GREEN)✅ OK$(NC)\n"; \
		else \
			echo "$(RED)❌ Failed on $$file$(NC)"; \
			exit 1; \
		fi; \
	done
	@echo "$(BOLD)✅ All seeders executed successfully.$(NC)"

help:
	@echo "Available targets:"
	@echo "  local        Run app with .env variables"
	@echo "  genmock      Generate mock files"
	@echo "  test         Run tests"
	@echo "  test-report  Generate HTML test coverage report"
	@echo "  seed         Run DB seeders from $(GREEN)$(BOLD)$(SEEDER_DIR)$(NC)"
