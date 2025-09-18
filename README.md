# Zentara Threat Intelligence Backend System

A secure, scalable threat intelligence backend system built with Go, implementing Clean Architecture principles with comprehensive monitoring and observability.

## 🏗️ Architecture

This system follows **Domain-Driven Design (DDD)**, **Clean Architecture**, and **Hexagonal Architecture** patterns:

- **Domain Layer**: Core business logic, entities, and value objects
- **Application Layer**: Use cases and business orchestration
- **Infrastructure Layer**: Database, Redis, JWT, New Relic integrations
- **Interface Layer**: HTTP handlers, middleware, and API routes

## 🚀 Features

### Core Functionality
- **JWT Authentication** with access & refresh tokens
- **Role-based Access Control** (Admin, Analyst, Viewer)
- **Order Management** for threat intelligence data
- **Rate Limiting** and security middleware
- **Comprehensive Logging** with structured JSON format

### Infrastructure
- **PostgreSQL** with GORM for data persistence
- **Redis** for caching and session management
- **New Relic APM** for monitoring and observability
- **Docker & Kubernetes** ready deployments
- **CI/CD Pipeline** with GitHub Actions

## 🛠️ Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin HTTP framework
- **Database**: PostgreSQL with GORM
- **Cache**: Redis
- **Monitoring**: New Relic APM
- **Authentication**: JWT tokens
- **Containerization**: Docker & Docker Compose
- **Orchestration**: Kubernetes
- **CI/CD**: GitHub Actions

## 📁 Project Structure

```
threat-intel-backend/
├── cmd/                    # Application entry point
├── domain/                 # Domain entities and business logic
├── application/            # Use cases and services
├── infrastructure/         # External integrations
│   ├── postgres/          # Database layer
│   ├── redis/             # Cache layer
│   ├── jwt/               # Authentication
│   └── newrelic/          # Monitoring
├── interfaces/            # HTTP handlers and middleware
├── configs/               # Configuration management
├── deployments/           # Kubernetes manifests
└── .github/workflows/     # CI/CD pipelines
```

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 15+
- Redis 7+

### Local Development

1. **Clone the repository**
```bash
git clone <repository-url>
cd threat-intel-backend
```

2. **Set up environment variables**
```bash
cp .env.example .env
# Edit .env with your configuration
```

3. **Run with Docker Compose**
```bash
docker-compose up --build
```

4. **Or run locally**
```bash
# Install dependencies
go mod download

# Run the application
go run cmd/main.go
```

The API will be available at `http://localhost:8080`

### API Documentation
Swagger documentation is available at: `http://localhost:8080/swagger/index.html`

## 🔐 Authentication

### Register a new user
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "role": "viewer"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### Create an order
```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-access-token>" \
  -d '{
    "item_id": "intel-basic",
    "quantity": 1
  }'
```

## 🐳 Docker Deployment

### Build and run with Docker Compose
```bash
docker-compose up --build
```

### Build Docker image
```bash
docker build -t threat-intel-backend .
```

## ☸️ Kubernetes Deployment

### Deploy to Kubernetes
```bash
# Apply all manifests
kubectl apply -f deployments/

# Check deployment status
kubectl get pods -n threat-intel
```

### Scale the application
```bash
kubectl scale deployment threat-intel-api --replicas=5 -n threat-intel
```

## 📊 Monitoring with New Relic

1. **Set up New Relic account** and get your license key
2. **Configure environment variables**:
   ```bash
   NEW_RELIC_LICENSE_KEY=your_license_key
   NEW_RELIC_APP_NAME=zentara-threat-intel-api
   ```
3. **Monitor your application** through New Relic dashboard

### Key Metrics Tracked
- API response times and throughput
- Database query performance
- Error rates and exceptions
- Custom business metrics
- Infrastructure metrics

## 🧪 Testing

### Run tests
```bash
go test -v ./...
```

### Run tests with coverage
```bash
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🔒 Security Features

- **JWT Authentication** with secure token handling
- **Password Hashing** using bcrypt
- **Role-based Access Control** with permission hierarchy
- **Rate Limiting** to prevent abuse
- **Input Validation** and sanitization
- **CORS** configuration
- **Security Headers** middleware

## 📈 Performance & Scalability

- **Horizontal Pod Autoscaling** in Kubernetes
- **Redis Caching** for improved performance
- **Connection Pooling** for database efficiency
- **Structured Logging** for observability
- **Health Checks** for reliability

## 🚀 CI/CD Pipeline

The GitHub Actions pipeline includes:
- **Automated Testing** with PostgreSQL and Redis services
- **Code Linting** with golangci-lint
- **Docker Image Building** and pushing to registry
- **Deployment** to staging and production environments
- **New Relic Deployment Markers** for release tracking

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite
6. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 👨‍💻 Author

**Rahmatullah Sidik**  
Backend/Infrastructure Engineer

---

## 🎯 Architecture Decisions

### Why Clean Architecture?
- **Separation of Concerns**: Each layer has a single responsibility
- **Testability**: Business logic is isolated and easily testable
- **Flexibility**: Easy to swap implementations without affecting business logic
- **Maintainability**: Clear boundaries make the codebase easier to understand and modify

### Why Go?
- **Performance**: Excellent performance for concurrent operations
- **Simplicity**: Clean syntax and powerful standard library
- **Concurrency**: Built-in goroutines for handling multiple requests
- **Deployment**: Single binary deployment with minimal dependencies

### Why PostgreSQL + Redis?
- **PostgreSQL**: ACID compliance, complex queries, and reliability
- **Redis**: High-performance caching and session storage
- **Complementary**: PostgreSQL for persistence, Redis for speed

This system is designed to handle **10k+ concurrent users** with proper scaling and monitoring in place.