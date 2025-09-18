# Zentara Threat Intelligence Backend System

## Overview
This project is a **Threat Intelligence Backend System** built with **Go (Golang)**, designed as part of the Zentara Backend/Infrastructure Engineer Technical Assessment. The system demonstrates skills in:
- Secure backend API development
- Authentication & Authorization
- Infrastructure as Code (IaC)
- Containerization & Orchestration (Docker & Kubernetes)
- CI/CD pipeline design
- Application of software architecture patterns (DDD, Onion, Hexagonal, Clean)
- **Monitoring & Observability with New Relic**

The project is designed to be **secure, scalable, observable, and cloud-ready**.

---

## Features

### Core API Features
- **Authentication System**
  - JWT-based authentication with access & refresh tokens
  - Role-based access control (Admin, Analyst, Viewer)
  - Password hashing with bcrypt
  - Login rate limiting & brute force protection

- **Security Analytics**
  - User activity logging & audit trails
  - API rate limiting per user role
  - Request/response logging for compliance

### Technical Stack
- **Language:** Go 1.21+
- **Framework:** Gin (HTTP framework)
- **Database:** PostgreSQL (GORM as ORM)
- **Caching/Session:** Redis
- **Authentication:** JWT with proper secret management
- **Validation:** Input validation & sanitization
- **Logging:** Structured logging (JSON format)
- **Configuration:** Environment-based config management
- **Monitoring:** New Relic APM integrated for metrics, transactions, and tracing

### Infrastructure
- **Containerization:** Multi-stage Dockerfile & Docker Compose
- **CI/CD:** GitHub Actions with automated testing, linting, and deployment
- **Kubernetes:** Scalable deployment with PostgreSQL StatefulSet, Redis, ConfigMaps, Secrets
- **Monitoring:** Health checks, New Relic dashboards, structured logs

---

## Repository Structure
```
root
├── cmd/                  # Application entrypoint
├── domain/               # Entities, Value Objects, Aggregates
│   └── order.go
├── application/          # Use cases (services, orchestrating domain logic)
│   └── order_service.go
├── infrastructure/       # DB, Redis, JWT, logging, monitoring
│   ├── postgres/
│   ├── redis/
│   ├── jwt/
│   └── newrelic/         # New Relic integration
├── interfaces/           # HTTP handlers, gRPC, CLI
│   └── http/
├── configs/              # Config files (env, yaml)
├── deployments/          # Kubernetes manifests
├── scripts/              # Automation scripts
└── README.md
```

---

## Example Use Case: User Ordering

### Scenario
A **user places an order for cybersecurity threat intelligence data** via the backend system.

### Flow
1. **User Registration/Login**  
   - The user registers and logs in with email & password.
   - JWT access & refresh tokens are issued.

2. **Placing an Order**  
   - Authenticated user sends a `POST /orders` request with:
     ```json
     {
       "item_id": "intel-basic",
       "quantity": 1
     }
     ```

3. **Domain Validation**  
   - The system validates that `item_id` exists.
   - Business rules (e.g., user role permissions) are enforced.

4. **Persistence**  
   - The order is stored in PostgreSQL.
   - Redis may cache recent orders for performance.

5. **Audit Logging & Monitoring**  
   - User activity is logged for compliance.
   - New Relic records API transaction performance, throughput, and error rates.

6. **Response**  
   - The API responds:
     ```json
     {
       "order_id": "12345",
       "status": "confirmed"
     }
     ```

---

## Architecture Approaches

This project applies **DDD, Onion, Hexagonal, and Clean Architecture principles** to ensure scalability, maintainability, and testability.

### Domain-Driven Design (DDD)
- **Entities & Value Objects:** Define business models (e.g., `User`, `Order`, `ThreatIndicator`).
- **Aggregates:** Group related entities with a root (e.g., `OrderAggregate`).
- **Domain Events:** Trigger events like `OrderPlacedEvent` for asynchronous processing.
- **Repositories:** Interface-based persistence abstraction.

### Onion Architecture
- **Domain Layer (Core):** Business logic, independent of frameworks.
- **Application Layer:** Use cases, services, orchestration of domain objects.
- **Infrastructure Layer:** Database, external APIs, caching, monitoring integrations.
- **Interfaces Layer:** HTTP, gRPC, CLI adapters.

### Hexagonal Architecture (Ports & Adapters)
- **Ports:** Interfaces for domain operations (e.g., `OrderRepository`, `NotificationService`).
- **Adapters:** Concrete implementations (PostgreSQL adapter, Redis cache, New Relic adapter, SMTP email).
- **Inbound Adapters:** HTTP handlers (Gin controllers).
- **Outbound Adapters:** Database, message broker, monitoring, external services.

### Clean Architecture
- Emphasis on **use cases as central** business logic.
- Dependency Rule: **Dependencies only point inward.**
- High-level policies (business rules) are not dependent on low-level details (database, frameworks).

---

## Deployment

### Local Development
```bash
# Run with Docker Compose
docker-compose up --build
```

### Kubernetes Deployment
```bash
kubectl apply -f deployments/
```

### CI/CD
- GitHub Actions workflow executes:
  - Linting & testing
  - Docker image build & push to container registry
  - Deploy to staging/production environments
  - New Relic deployment markers for release tracking

---

## Monitoring with New Relic
- **APM Integration:** Each request/transaction is instrumented and sent to New Relic.
- **Custom Metrics:** Track order processing latency, DB query duration, Redis cache hits/misses.
- **Dashboards:** Visualize API performance, throughput, and error rates.
- **Alerts:** Configured thresholds trigger alerts on anomalies.

Example environment variable setup:
```env
NEW_RELIC_LICENSE_KEY=your_newrelic_license_key
NEW_RELIC_APP_NAME=zentara-threat-intel-api
```

---

## Documentation
- **API Docs:** Swagger/OpenAPI available at `/swagger/index.html`
- **Deployment Guide:** See [docs/deployment.md]
- **Architecture Diagram:** See [docs/architecture.png]
- **Monitoring Guide:** See [docs/monitoring.md]

---

## Evaluation Preparation
Be ready to discuss:
1. Why **Gin + PostgreSQL + Redis** were chosen
2. How JWT & role-based access control were implemented
3. Strategies for handling **10k+ concurrent users**
4. CI/CD design (GitHub Actions workflows with New Relic release tracking)
5. Kubernetes HA design with scaling
6. Caching strategies for read-heavy endpoints
7. Monitoring & debugging approach (New Relic APM + structured logs)
8. (Optional) LLM integration for threat intelligence enrichment

---

## Author
**Rahmatulah Sidik**  
Backend/Infrastructure Engineer C