.DEFAULT_GOAL := help
GOCACHE ?= $(CURDIR)/work/go-build-cache
export GOCACHE

.PHONY: help setup dev test lint migrate seed worker frontend backend docker-up docker-down
help:
	@echo "make dev | test | lint | docker-up | docker-down"
setup:
	@copy .env.example .env 2>NUL || true
dev backend:
	go run ./backend/cmd/api
test:
	go test ./...
lint:
	gofmt -w backend
	migrate:
	@echo "Migrations run automatically by the Docker MySQL init mount. See migrations/."
seed:
	@echo "The API seeds the editable candidate profile in development."
worker:
	@echo "Phase 1 has no background worker yet."
frontend:
	@echo "The dashboard is served by the API at http://localhost:8080"
docker-up:
	docker compose up --build
docker-down:
	docker compose down
