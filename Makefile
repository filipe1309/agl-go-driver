.PHONY: help run test up-db

# run tests
test:
	@echo "🟢 Running tests..."
	go test ./internal/... -v -coverprofile=coverage.out

coverage:
	@echo "🟢 Running coverage..."
	go tool cover -html=coverage.out

build:
	@echo "🏗 Building code..."
	./scripts/shell/build.sh

run:
	@echo "🏁 Running code..."

run-api:
	@echo "🏁 Running API..."
	. ./scripts/shell/env.sh
	go run cmd/api/main.go

run-worker:
	@echo "🏁 Running Worker..."
	. ./scripts/shell/env.sh
	go run cmd/worker/main.go

up-db:
	@echo "🏁 Starting database..."
	docker run --name imersao-postgres -e POSTGRES_PASSWORD=1234 -p 5432:5432 -d postgres

up-queue:
	@echo "🏁 Starting queue..."
	docker run --name imersao-rabbit --hostname aprenda-golang -p 5672:5672 -p 15672:15672 -d rabbitmq:3-management

up-all: up-db up-queue

help:
	@echo "📖 Available commands:"
	@echo "  make run"
	@echo "  make test"
	@echo "  make help"
