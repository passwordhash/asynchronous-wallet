ENV_FILE ?= config.env

.EXPORT_ALL_VARIABLES:
include $(ENV_FILE)

# === Containers ===
.PHONY: container-start container-down

container-start:
	@docker compose -f ./docker-compose.yml up -d

container-down:
	@docker compose -f ./docker-compose.yml down

container-infra:
	@docker compose -f ./docker-compose.yml up postgres migrate -d

# === Local dev ===
.PHONY: start-sandbox run dev

run:
	@go run cmd/http_server/main.go -config=./configs/local.yml

dev: container-infra run

# === Tests ===
.PHONY: unit-test integration-test

unit-test:
	@go test ./internal/...

integration-test:
	APP_OUT_PORT=$(APP_OUT_PORT) go test ./tests/...
