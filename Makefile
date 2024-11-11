# Simple Makefile for a Go project
CWD := ${shell pwd}

# COLORS
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
RESET  := $(shell tput -Txterm sgr0)

TARGET_MAX_CHAR_NUM=20

.PHONY: vendor test

## Show help
help:
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  ${YELLOW}%-$(TARGET_MAX_CHAR_NUM)s${RESET} ${GREEN}%s${RESET}\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)


## Create DB container
docker-dev-up:
	@if docker compose  -f docker-compose.dev.yml up -d 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose"; \
	fi

## Shutdown DB container
docker-dev-down:
	@if docker compose -f docker-compose.dev.yml down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose"; \
		docker-compose  -f docker-compose.dev.yml down; \
	fi

## Create migration script example: make migrate-create name=create_users_table
migrate-create:
	migrate create -ext sql -dir db/migrations -seq $(name)

## Run migration script example: make migrate-up
migrate-up:
	migrate -path db/migrations -database "postgres://postgres:postgres@localhost:5432/task-services?sslmode=disable" up

## Rollback migration script example: make migrate-down
migrate-down:
	migrate -path db/migrations -database "postgres://postgres:postgres@localhost:5432/task-services?sslmode=disable" down

## Generate swagger documentation build into docs folder -> swagger/index.html
# output into docs.json, docs.yaml, docs.go
swagger:
	swag init -g cmd/main.go --outputTypes go,yaml,json  --output docs 

## Run the project
run:
	make swagger
	go run cmd/main.go

## Setup first time of the project
install:
	brew update

	bash < <(curl -sSL https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
	gvm install go1.18
	gvm use go1.18 --default


	brew install golang-migrate
	go install -v github.com/golang/mock/mockgen@latest
	
	go mod tidy

## generate mocks for task-service
mock-task-service:
	@echo "Generating mocks for task-service..."
	@mockgen -source=./internal/domains/tasks/interfaces/index.go -destination=./internal/mocks/tasks/interfaces/index.go -package=mocks
	@mockgen -source=./internal/domains/tasks/usecases/index.go -destination=./internal/mocks/tasks/usecases/index.go -package=mocks

## test the project
test:
	go test -timeout 30s -coverprofile=coverage.out ./...