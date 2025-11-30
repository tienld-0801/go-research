.PHONY: run install update-dependencies docker-run docker-down swagger-generate help

# Variables
DOCKER_COMPOSE_FILE := docker/docker-compose.yaml

# Default target
.DEFAULT_GOAL := help

# run project
run:
	@if [ -f .env ]; then export $$(grep -v '^#' .env | xargs); fi; air

# install dependencies
install:
	go mod tidy

# update all dependencies
update-dependencies:
	go get -u ./...

# generate swagger documentation
swagger-generate:
	swag init -g cmd/main.go -o docs

# run docker
docker-run:
	docker compose -f $(DOCKER_COMPOSE_FILE) up -d

# down docker
docker-down:
	docker compose -f $(DOCKER_COMPOSE_FILE) down

# help
help:
	@echo "Available commands:"
	@echo "  make run                 - Run the project with air (hot reload)"
	@echo "  make install             - Install dependencies (go mod tidy)"
	@echo "  make update-dependencies - Update all dependencies"
	@echo "  make swagger-generate    - Generate Swagger documentation"
	@echo "  make docker-run          - Start docker containers"
	@echo "  make docker-down         - Stop docker containers"
