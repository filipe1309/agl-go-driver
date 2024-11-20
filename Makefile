.PHONY: help run test up-db

# run tests
test:
	@echo "🟢 Running tests..."
	go test ./internal/... -v

# run node
run:
	@echo "🏁 Running code..."

up-db:
	@echo "🏁 Starting database..."
	docker run --name imersao-postgres -e POSTGRES_PASSWORD=1234 -p 5432:5432 -d postgres

up-queue:
	@echo "🏁 Starting queue..."
	docker run --name imersao-rabbit --hostname aprenda-golang -p 5672:5672 -d rabbitmq:3

help:
	@echo "📖 Available commands:"
	@echo "  make run"
	@echo "  make test"
	@echo "  make help"
