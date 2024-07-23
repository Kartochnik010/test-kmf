## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## run: run the project
.PHONY: run
run:
	@go run ./cmd/app

## build: format swagger and build the project
.PHONY: build
build:
	make swag
	go build -o ./bin/main ./cmd/app

.PHONY: migrate
## migrate: apply up migrations stored
migrate:
	@docker compose -f ${DOCKER_COMPOSE_FILE} --profile tools run --rm migrate up

.PHONY: migrate-down
## migrate-down: apply down migrations stored
migrate-down:
	@docker compose -f ${DOCKER_COMPOSE_FILE} --profile tools run --rm migrate drop

.PHONY: test
## test: run tests
test: 
	@go test -v ./...

## down: remove containers create by docker compose
down:
	@docker compose down

## up: compose project up
up:
	@docker compose up -d && docker ps

## swag: generate swagger documentation
swag:
	swag init -g cmd/app/main.go -o ./docs --parseDependency
