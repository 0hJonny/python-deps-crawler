
# Переменные
NAMESPACE := python-deps-crawler
KUBECTL := kubectl
KAFKA_POD := $(shell $(KUBECTL) get pods -n $(NAMESPACE) -l app=kafka -o jsonpath='{.items[0].metadata.name}')

PROTO_DIR := api/proto
PROTO_OUT_DIR := pkg/proto
PROTO_FILES := $(wildcard $(PROTO_DIR)/*.proto)

# Цвета для вывода
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[0;33m
BLUE := \033[0;34m
NC := \033[0m # No Color

.PHONY: help up down init-kafka-topics proto-gen proto-clean

help: ## Помощь
	@echo "$(BLUE)Доступные команды:$(NC)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

up: ## Запустить инфраструктуру
	@echo "$(BLUE)Запуск инфраструктуры...$(NC)"
	@$(KUBECTL) apply -f deployments/k8s/develop/namespace.yaml
	@$(KUBECTL) apply -f deployments/k8s/develop/configmap.yaml
	@$(KUBECTL) apply -f deployments/k8s/develop/secrets.yaml
	@echo "$(YELLOW)Запуск PostgreSQL...$(NC)"
	@$(KUBECTL) apply -f deployments/k8s/develop/postgres.yaml
	@echo "$(YELLOW)Запуск Kafka...$(NC)"
	@$(KUBECTL) apply -f deployments/k8s/develop/kafka.yaml
	@echo "$(YELLOW)Запуск Redis...$(NC)"
	@$(KUBECTL) apply -f deployments/k8s/develop/redis.yaml
	@echo "$(YELLOW)Запуск Prometheus...$(NC)"
	@$(KUBECTL) apply -f deployments/k8s/develop/prometheus.yaml
	@echo "$(YELLOW)Запуск среды разработки Go...$(NC)"
	@$(KUBECTL) apply -f deployments/k8s/develop/go-dev.yaml
	@echo "$(GREEN)Инфраструктура запущена!$(NC)"
	@echo "$(YELLOW)Ожидание готовности подов...$(NC)"
	@$(KUBECTL) wait --for=condition=ready pod -l app=kafka -n $(NAMESPACE) --timeout=300s
	@$(KUBECTL) wait --for=condition=ready pod -l app=postgres -n $(NAMESPACE) --timeout=300s
	@$(KUBECTL) wait --for=condition=ready pod -l app=redis -n $(NAMESPACE) --timeout=300s
	@echo "$(GREEN)Все поды готовы!$(NC)"

down: ## Остановить всю инфраструктуру
	@echo "$(RED)Остановка инфраструктуры...$(NC)"
	@$(KUBECTL) delete -f deployments/k8s/develop/go-dev.yaml --ignore-not-found=true
	@$(KUBECTL) delete -f deployments/k8s/develop/prometheus.yaml --ignore-not-found=true
	@$(KUBECTL) delete -f deployments/k8s/develop/redis.yaml --ignore-not-found=true
	@$(KUBECTL) delete -f deployments/k8s/develop/kafka.yaml --ignore-not-found=true
	@$(KUBECTL) delete -f deployments/k8s/develop/postgres.yaml --ignore-not-found=true
	@echo "$(YELLOW)Очистка PVC...$(NC)"
	@$(KUBECTL) delete pvc -n $(NAMESPACE) --all --ignore-not-found=true
	@$(KUBECTL) delete -f deployments/k8s/develop/secrets.yaml --ignore-not-found=true
	@$(KUBECTL) delete -f deployments/k8s/develop/configmap.yaml --ignore-not-found=true
	@echo "$(RED)Удаление namespace...$(NC)"
	@$(KUBECTL) delete namespace $(NAMESPACE) --ignore-not-found=true
	@echo "$(GREEN)Инфраструктура остановлена!$(NC)"

init-kafka-topics: ## Создать топики Kafka для микросервисов
	@echo "$(BLUE)Создание топиков Kafka...$(NC)"
	@if [ -z "$(KAFKA_POD)" ]; then \
		echo "$(RED)Ошибка: Kafka pod не найден!$(NC)"; \
		exit 1; \
	fi
	@echo "$(YELLOW)Создание топика dependency.package.fetch...$(NC)"
	@$(KUBECTL) exec -n $(NAMESPACE) $(KAFKA_POD) -- /opt/kafka/bin/kafka-topics.sh \
		--bootstrap-server localhost:9092 \
		--create \
		--if-not-exists \
		--topic dependency.package.fetch \
		--partitions 3 \
		--replication-factor 1
	@echo "$(YELLOW)Создание топика dependency.package.response...$(NC)"
	@$(KUBECTL) exec -n $(NAMESPACE) $(KAFKA_POD) -- /opt/kafka/bin/kafka-topics.sh \
		--bootstrap-server localhost:9092 \
		--create \
		--if-not-exists \
		--topic dependency.package.response \
		--partitions 3 \
		--replication-factor 1
	@echo "$(YELLOW)Создание топика dependency.analysis.request...$(NC)"
	@$(KUBECTL) exec -n $(NAMESPACE) $(KAFKA_POD) -- /opt/kafka/bin/kafka-topics.sh \
		--bootstrap-server localhost:9092 \
		--create \
		--if-not-exists \
		--topic dependency.analysis.request \
		--partitions 3 \
		--replication-factor 1
	@echo "$(YELLOW)Создание топика dependency.analysis.response...$(NC)"
	@$(KUBECTL) exec -n $(NAMESPACE) $(KAFKA_POD) -- /opt/kafka/bin/kafka-topics.sh \
		--bootstrap-server localhost:9092 \
		--create \
		--if-not-exists \
		--topic dependency.analysis.response \
		--partitions 3 \
		--replication-factor 1
	@echo "$(YELLOW)Создание топика dependency.graph.request...$(NC)"
	@$(KUBECTL) exec -n $(NAMESPACE) $(KAFKA_POD) -- /opt/kafka/bin/kafka-topics.sh \
		--bootstrap-server localhost:9092 \
		--create \
		--if-not-exists \
		--topic dependency.graph.request \
		--partitions 3 \
		--replication-factor 1
	@echo "$(YELLOW)Создание топика dependency.graph.response...$(NC)"
	@$(KUBECTL) exec -n $(NAMESPACE) $(KAFKA_POD) -- /opt/kafka/bin/kafka-topics.sh \
		--bootstrap-server localhost:9092 \
		--create \
		--if-not-exists \
		--topic dependency.graph.response \
		--partitions 3 \
		--replication-factor 1
	@echo "$(YELLOW)Создание топика dependency.status.request...$(NC)"
	@$(KUBECTL) exec -n $(NAMESPACE) $(KAFKA_POD) -- /opt/kafka/bin/kafka-topics.sh \
		--bootstrap-server localhost:9092 \
		--create \
		--if-not-exists \
		--topic dependency.status.request \
		--partitions 3 \
		--replication-factor 1
	@echo "$(YELLOW)Создание топика dependency.status.response...$(NC)"
	@$(KUBECTL) exec -n $(NAMESPACE) $(KAFKA_POD) -- /opt/kafka/bin/kafka-topics.sh \
		--bootstrap-server localhost:9092 \
		--create \
		--if-not-exists \
		--topic dependency.status.response \
		--partitions 3 \
		--replication-factor 1
	@echo "$(GREEN)Все топики созданы!$(NC)"

proto-gen: ## Сгенерировать Go код из proto файлов
	@echo "Генерация Go кода из proto файлов..."
	@mkdir -p $(PROTO_OUT_DIR)
	@for proto in $(PROTO_FILES); do \
		name=$$(basename $$proto .proto); \
		mkdir -p $(PROTO_OUT_DIR)/$$name; \
		echo "Обработка $$proto..."; \
		protoc --go_out=$(PROTO_OUT_DIR)/$$name --go_opt=paths=source_relative \
		       --go-grpc_out=$(PROTO_OUT_DIR)/$$name --go-grpc_opt=paths=source_relative \
		       --proto_path=$(PROTO_DIR) $$proto; \
	done
	@echo "Proto файлы скомпилированы!"

proto-clean: ## Очистить сгенерированные proto файлы
	@echo "Очистка сгенерированных файлов..."
	@rm -rf $(PROTO_OUT_DIR)/*
	@echo "Очистка завершена!"

.PHONY: pgadmin-start
pgadmin-start: ## Запустить pgadmin
	@$(KUBECTL) apply -f deployments/k8s/develop/pgadmin.yaml
	@echo "$(YELLOW)Запуск pgAdmin...$(NC)"