# Makefile для Zloy Billing Service

.PHONY: help build test clean dev logs status migrate-up migrate-down migrate-status migrate-reset

# Переменные
APP_NAME = zloy-billing
COMPOSE_FILE = docker-compose.yml

# Цвета для вывода
GREEN = \033[0;32m
YELLOW = \033[1;33m
RED = \033[0;31m
NC = \033[0m # No Color

help: ## Показать справку
	@echo "$(GREEN)Доступные команды:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-15s$(NC) %s\n", $$1, $$2}'


start: ## Запустить приложение в Docker
	@echo "$(GREEN)Запуск в Docker...$(NC)"
	docker-compose up -d --build

stop: ## Остановить Docker контейнеры
	@echo "$(GREEN)Остановка Docker контейнеров...$(NC)"
	docker-compose -f $(COMPOSE_FILE) down

logs: ## Показать логи Docker контейнеров
	@echo "$(GREEN)Логи Docker контейнеров:$(NC)"
	docker-compose -f $(COMPOSE_FILE) logs -f

status: ## Показать статус сервисов
	@echo "$(GREEN)Статус сервисов:$(NC)"
	docker-compose -f $(COMPOSE_FILE) ps

logs-app: ## Логи только приложения
	@echo "$(GREEN)Логи приложения:$(NC)"
	docker-compose -f $(COMPOSE_FILE) logs -f billing

logs-db: ## Логи баз данных
	@echo "$(GREEN)Логи PostgreSQL:$(NC)"
	docker-compose -f $(COMPOSE_FILE) logs postgres
	@echo "$(GREEN)Логи MongoDB:$(NC)"
	docker-compose -f $(COMPOSE_FILE) logs mongo

tidy: ## Очистка зависимостей Go
	@echo "$(GREEN)Очистка зависимостей...$(NC)"
	go mod tidy

# === Команды миграций ===

migrate-up: ## Применить все миграции
	@echo "$(GREEN)Применение миграций PostgreSQL...$(NC)"
	@for file in migrations/postgres/*.sql; do \
		echo "Выполняю: $$file"; \
		docker exec -i zloy-postgres-1 psql -U postgres -d zl0y_billing < "$$file"; \
	done
	@echo "$(GREEN)Миграции применены успешно!$(NC)"

migrate-down: ## Откатить последнюю миграцию (очистить все таблицы)
	@echo "$(YELLOW)Внимание: Это удалит все данные!$(NC)"
	@echo "$(RED)Очистка базы данных...$(NC)"
	docker exec -i zloy-postgres-1 psql -U postgres -d zl0y_billing -c "DROP TABLE IF EXISTS users CASCADE;"
	@echo "$(GREEN)База данных очищена!$(NC)"

migrate-status: ## Показать состояние базы данных
	@echo "$(GREEN)Статус базы данных PostgreSQL:$(NC)"
	@echo "$(YELLOW)Список таблиц:$(NC)"
	docker exec -i zloy-postgres-1 psql -U postgres -d zl0y_billing -c "\dt"
	@echo "$(YELLOW)Список индексов:$(NC)"
	docker exec -i zloy-postgres-1 psql -U postgres -d zl0y_billing -c "\di"

migrate-reset: ## Пересоздать базу данных с миграциями
	@echo "$(RED)Пересоздание базы данных...$(NC)"
	$(MAKE) migrate-down
	$(MAKE) migrate-up
	@echo "$(GREEN)База данных пересоздана!$(NC)"

migrate-create: ## Создать новую миграцию (использование: make migrate-create NAME=migration_name)
	@if [ -z "$(NAME)" ]; then \
		echo "$(RED)Ошибка: Укажите имя миграции. Пример: make migrate-create NAME=add_user_email$(NC)"; \
		exit 1; \
	fi
	@NEXT_NUM=$$(ls migrations/postgres/*.sql 2>/dev/null | wc -l | xargs expr 1 +); \
	FILENAME="migrations/postgres/$$(printf "%03d" $$NEXT_NUM)_$(NAME).sql"; \
	echo "-- Migration: $$(basename $$FILENAME)" > $$FILENAME; \
	echo "-- Description: $(NAME)" >> $$FILENAME; \
	echo "" >> $$FILENAME; \
	echo "-- Добавьте ваш SQL код здесь" >> $$FILENAME; \
	echo "$(GREEN)Создана миграция: $$FILENAME$(NC)"

# === Команды базы данных ===

db-connect: ## Подключиться к PostgreSQL
	@echo "$(GREEN)Подключение к PostgreSQL...$(NC)"
	docker exec -it zloy-postgres-1 psql -U postgres -d zl0y_billing

db-mongo: ## Подключиться к MongoDB
	@echo "$(GREEN)Подключение к MongoDB...$(NC)"
	docker exec -it zloy-mongo-1 mongosh zl0y_billing
	