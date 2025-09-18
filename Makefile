# Makefile for Threat Intelligence Backend

.PHONY: help build run test clean docker-build docker-run k8s-deploy

# Variables
APP_NAME=threat-intel-backend
DOCKER_IMAGE=threat-intel-backend:latest
NAMESPACE=threat-intel

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building $(APP_NAME)..."
	go build -o bin/$(APP_NAME) cmd/main.go

run: ## Run the application locally
	@echo "Running $(APP_NAME)..."
	go run cmd/main.go

test: ## Run tests
	@echo "Running tests..."
	go test -v -race ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

docker-run: ## Run with Docker Compose
	@echo "Running with Docker Compose..."
	docker-compose up --build

docker-down: ## Stop Docker Compose
	@echo "Stopping Docker Compose..."
	docker-compose down

k8s-deploy: ## Deploy to Kubernetes
	@echo "Deploying to Kubernetes..."
	kubectl apply -f deployments/

k8s-delete: ## Delete from Kubernetes
	@echo "Deleting from Kubernetes..."
	kubectl delete -f deployments/

k8s-logs: ## View application logs in Kubernetes
	@echo "Viewing logs..."
	kubectl logs -f deployment/threat-intel-api -n $(NAMESPACE)

k8s-status: ## Check Kubernetes deployment status
	@echo "Checking deployment status..."
	kubectl get pods,svc,deploy -n $(NAMESPACE)

dev-setup: ## Set up development environment
	@echo "Setting up development environment..."
	cp .env.example .env
	@echo "Please edit .env file with your configuration"

swagger: ## Generate Swagger documentation
	@echo "Generating Swagger docs..."
	swag init -g cmd/main.go

format: ## Format code
	@echo "Formatting code..."
	go fmt ./...

security-scan: ## Run security scan
	@echo "Running security scan..."
	gosec ./...

all: clean deps lint test build ## Run all checks and build