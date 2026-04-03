# ZarishSphere FHIR Engine Makefile

.PHONY: help build run test clean docker-build docker-up docker-down migrate lint format

# Default target
help:
	@echo "ZarishSphere FHIR Engine - Development Commands"
	@echo ""
	@echo "Development:"
	@echo "  build      - Build the application"
	@echo "  run        - Run the application"
	@echo "  test       - Run tests"
	@echo "  lint       - Run linter"
	@echo "  format     - Format code"
	@echo "  clean      - Clean build artifacts"
	@echo ""
	@echo "Database:"
	@echo "  migrate    - Run database migrations"
	@echo "  db-reset   - Reset database"
	@echo ""
	@echo "Docker:"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-up    - Start Docker Compose"
	@echo "  docker-down  - Stop Docker Compose"
	@echo "  docker-logs  - Show Docker logs"
	@echo ""
	@echo "Production:"
	@echo "  deploy-staging   - Deploy to staging"
	@echo "  deploy-production - Deploy to production"

# Development targets
build:
	@echo "Building ZarishSphere FHIR Engine..."
	mkdir -p build
	go build -o build/zs-core-fhir-engine ./cmd/fhir-engine

run: build
	@echo "Starting ZarishSphere FHIR Engine..."
	./build/zs-core-fhir-engine serve

test:
	@echo "Running tests..."
	go test -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint:
	@echo "Running linter..."
	golangci-lint run

format:
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

clean:
	@echo "Cleaning build artifacts..."
	rm -rf build/
	rm -f coverage.out coverage.html
	go clean

# Database targets
migrate:
	@echo "Running database migrations..."
	psql -h localhost -U postgres -d zs_fhir -f deploy/migrations/001_create_tables.sql

db-reset:
	@echo "Resetting database..."
	dropdb -h localhost -U postgres zs_fhir || true
	createdb -h localhost -U postgres zs_fhir
	$(MAKE) migrate

# Docker targets
docker-build:
	@echo "Building Docker image..."
	docker build -t zarishsphere/zs-fhir-engine:latest .

docker-up:
	@echo "Starting Docker Compose..."
	docker-compose up -d

docker-down:
	@echo "Stopping Docker Compose..."
	docker-compose down

docker-logs:
	@echo "Showing Docker logs..."
	docker-compose logs -f

docker-restart: docker-down docker-up

# Production targets
deploy-staging:
	@echo "Deploying to staging..."
	# Add staging deployment commands here

deploy-production:
	@echo "Deploying to production..."
	# Add production deployment commands here

# Development utilities
install-deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

update-deps:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

security-scan:
	@echo "Running security scan..."
	gosec ./...

benchmark:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# Documentation
docs:
	@echo "Generating documentation..."
	godoc -http=:6060 &
	@echo "Documentation available at http://localhost:6060"

# Quick start for development
dev-setup: install-deps docker-up
	@echo "Setting up development environment..."
	sleep 10
	$(MAKE) migrate
	@echo "Development environment ready!"
	@echo "FHIR Engine: http://localhost:8080"
	@echo "Keycloak: http://localhost:8081"
	@echo "Grafana: http://localhost:3000"
	@echo "Prometheus: http://localhost:9090"

# CI/CD helpers
ci-test:
	@echo "Running CI tests..."
	go test -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

ci-build:
	@echo "CI build..."
	$(MAKE) build
	$(MAKE) lint
	$(MAKE) security-scan

# Health checks
health:
	@echo "Checking application health..."
	curl -f http://localhost:8080/health || echo "FHIR Engine is not healthy"

health-detailed:
	@echo "Checking detailed health..."
	curl -f http://localhost:8080/health/detailed

# Load testing
load-test:
	@echo "Running load tests..."
	hey -n 1000 -c 10 http://localhost:8080/fhir/R5/Patient

# Backup and restore
backup:
	@echo "Creating backup..."
	docker exec postgres pg_dump -U postgres zs_fhir > backup_$(shell date +%Y%m%d_%H%M%S).sql

restore:
	@echo "Restoring from backup..."
	@read -p "Enter backup file: " backup_file; \
	docker exec -i postgres psql -U postgres zs_fhir < $$backup_file

# Monitoring
logs:
	@echo "Following application logs..."
	tail -f logs/zs-fhir-engine.log

metrics:
	@echo "Fetching metrics..."
	curl http://localhost:8080/metrics

# Version info
version:
	@echo "ZarishSphere FHIR Engine"
	@echo "Go version: $(shell go version)"
	@echo "Git commit: $(shell git rev-parse --short HEAD)"
	@echo "Build time: $(shell date)"

# Release management
release-patch:
	@echo "Creating patch release..."
	$(MAKE) test
	$(MAKE) lint
	$(MAKE) docker-build
	# Add release commands here

release-minor:
	@echo "Creating minor release..."
	$(MAKE) test
	$(MAKE) lint
	$(MAKE) docker-build
	# Add release commands here

# Environment setup
setup-dev:
	@echo "Setting up development environment..."
	@if ! command -v go &> /dev/null; then \
		echo "Go is not installed. Please install Go 1.26.1 or higher."; \
		exit 1; \
	fi
	@if ! command -v docker &> /dev/null; then \
		echo "Docker is not installed. Please install Docker."; \
		exit 1; \
	fi
	@if ! command -v docker-compose &> /dev/null; then \
		echo "Docker Compose is not installed. Please install Docker Compose."; \
		exit 1; \
	fi
	$(MAKE) install-deps
	@echo "Development environment setup complete!"

# Quick commands for common tasks
quick-test: test lint
quick-build: build lint
quick-deploy: docker-build docker-up
