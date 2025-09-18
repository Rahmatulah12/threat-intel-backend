# Class Diagram - Zentara Threat Intelligence Backend System

## Arsitektur Clean Architecture dengan Domain-Driven Design

```mermaid
classDiagram
    %% Domain Layer - Core Business Logic
    class User {
        +UUID ID
        +string Email
        +string PasswordHash
        +UserRole Role
        +bool IsActive
        +time.Time CreatedAt
        +time.Time UpdatedAt
        +NewUser(email, password, role) *User
        +ValidatePassword(password) bool
        +HasPermission(requiredRole) bool
    }

    class UserRole {
        <<enumeration>>
        +RoleAdmin
        +RoleAnalyst
        +RoleViewer
    }

    class Order {
        +UUID ID
        +UUID UserID
        +string ItemID
        +int Quantity
        +OrderStatus Status
        +time.Time CreatedAt
        +time.Time UpdatedAt
        +User User
    }

    class OrderStatus {
        <<enumeration>>
        +OrderStatusPending
        +OrderStatusConfirmed
        +OrderStatusCompleted
        +OrderStatusCancelled
    }

    class OrderAggregate {
        +Order Order
        +NewOrder(userID, itemID, quantity) *OrderAggregate
        +Confirm()
        +Complete()
        +Cancel()
    }

    %% Domain Interfaces (Ports)
    class UserRepository {
        <<interface>>
        +Save(user *User) error
        +FindByID(id UUID) (*User, error)
        +FindByEmail(email string) (*User, error)
    }

    class OrderRepository {
        <<interface>>
        +Save(order *Order) error
        +FindByID(id UUID) (*Order, error)
        +FindByUserID(userID UUID) ([]*Order, error)
    }

    %% Application Layer - Use Cases
    class AuthService {
        -userRepo UserRepository
        -jwtService *JWTService
        +NewAuthService(userRepo, jwtService) *AuthService
        +Login(req LoginRequest) (*AuthResponse, error)
        +Register(req RegisterRequest) (*AuthResponse, error)
        +RefreshToken(refreshToken string) (*AuthResponse, error)
    }

    class OrderService {
        -orderRepo OrderRepository
        -userRepo UserRepository
        +NewOrderService(orderRepo, userRepo) *OrderService
        +CreateOrder(userID UUID, req CreateOrderRequest) (*OrderResponse, error)
        +GetOrder(orderID, userID UUID) (*Order, error)
        +GetUserOrders(userID UUID) ([]*Order, error)
    }

    class LoginRequest {
        +string Email
        +string Password
    }

    class RegisterRequest {
        +string Email
        +string Password
        +UserRole Role
    }

    class AuthResponse {
        +string AccessToken
        +string RefreshToken
        +*User User
    }

    class CreateOrderRequest {
        +string ItemID
        +int Quantity
    }

    class OrderResponse {
        +string OrderID
        +OrderStatus Status
    }

    %% Infrastructure Layer - External Integrations
    class PostgresUserRepository {
        -db *gorm.DB
        +NewUserRepository(db) *PostgresUserRepository
        +Save(user *User) error
        +FindByID(id UUID) (*User, error)
        +FindByEmail(email string) (*User, error)
    }

    class PostgresOrderRepository {
        -db *gorm.DB
        +NewOrderRepository(db) *PostgresOrderRepository
        +Save(order *Order) error
        +FindByID(id UUID) (*Order, error)
        +FindByUserID(userID UUID) ([]*Order, error)
    }

    class Database {
        +NewConnection(config Config) (*gorm.DB, error)
        +Migrate(db *gorm.DB) error
    }

    class DatabaseConfig {
        +string Host
        +string Port
        +string User
        +string Password
        +string DBName
        +string SSLMode
    }

    class JWTService {
        -secretKey []byte
        -accessTokenTTL time.Duration
        -refreshTokenTTL time.Duration
        +NewService(secretKey string) *JWTService
        +GenerateAccessToken(userID UUID, role UserRole) (string, error)
        +GenerateRefreshToken(userID UUID) (string, error)
        +ValidateAccessToken(tokenString string) (*Claims, error)
        +ValidateRefreshToken(tokenString string) (UUID, error)
    }

    class Claims {
        +UUID UserID
        +UserRole Role
        +jwt.RegisteredClaims
    }

    class RedisClient {
        -rdb *redis.Client
        +NewClient(addr, password string, db int) *RedisClient
        +Set(ctx, key string, value interface{}, expiration time.Duration) error
        +Get(ctx, key string) (string, error)
        +Del(ctx, keys ...string) error
        +Exists(ctx, keys ...string) (int64, error)
        +Ping(ctx) error
        +Close() error
    }

    class NewRelicMonitor {
        -app *newrelic.Application
        +NewMonitor(licenseKey, appName string) (*NewRelicMonitor, error)
        +GetApplication() *newrelic.Application
        +RecordCustomEvent(eventType string, params map[string]interface{})
        +RecordCustomMetric(name string, value float64)
        +Shutdown()
    }

    %% Interface Layer - HTTP Handlers & Middleware
    class Handler {
        -authService *AuthService
        -orderService *OrderService
        -logger *logrus.Logger
        +NewHandler(authService, orderService, logger) *Handler
        +Health(c *gin.Context)
        +Login(c *gin.Context)
        +Register(c *gin.Context)
        +RefreshToken(c *gin.Context)
        +CreateOrder(c *gin.Context)
        +GetOrder(c *gin.Context)
        +GetUserOrders(c *gin.Context)
    }

    class Middleware {
        -jwtService *JWTService
        -logger *logrus.Logger
        +NewMiddleware(jwtService, logger) *Middleware
        +CORS() gin.HandlerFunc
        +Logger() gin.HandlerFunc
        +Auth() gin.HandlerFunc
        +RequireRole(requiredRole UserRole) gin.HandlerFunc
        +RateLimit() gin.HandlerFunc
        +NewRelic(app *newrelic.Application) gin.HandlerFunc
    }

    class Router {
        -handler *Handler
        -middleware *Middleware
        +NewRouter(handler, middleware) *Router
        +Setup(app *newrelic.Application) *gin.Engine
    }

    %% Configuration
    class Config {
        +ServerConfig Server
        +DatabaseConfig Database
        +RedisConfig Redis
        +JWTConfig JWT
        +NewRelicConfig NewRelic
        +Load() *Config
    }

    class ServerConfig {
        +string Port
        +string Host
    }

    class RedisConfig {
        +string Addr
        +string Password
        +int DB
    }

    class JWTConfig {
        +string SecretKey
    }

    class NewRelicConfig {
        +string LicenseKey
        +string AppName
    }

    %% Relationships - Domain Layer
    User ||--|| UserRole : has
    Order ||--|| OrderStatus : has
    Order }|--|| User : belongs_to
    OrderAggregate ||--|| Order : contains

    %% Relationships - Application Layer Dependencies
    AuthService --> UserRepository : uses
    AuthService --> JWTService : uses
    OrderService --> OrderRepository : uses
    OrderService --> UserRepository : uses

    %% Relationships - Infrastructure Layer Implementations
    PostgresUserRepository ..|> UserRepository : implements
    PostgresOrderRepository ..|> OrderRepository : implements
    PostgresUserRepository --> Database : uses
    PostgresOrderRepository --> Database : uses
    Database --> DatabaseConfig : uses

    %% Relationships - Interface Layer Dependencies
    Handler --> AuthService : uses
    Handler --> OrderService : uses
    Middleware --> JWTService : uses
    Router --> Handler : uses
    Router --> Middleware : uses

    %% Relationships - Configuration
    Config --> ServerConfig : contains
    Config --> DatabaseConfig : contains
    Config --> RedisConfig : contains
    Config --> JWTConfig : contains
    Config --> NewRelicConfig : contains

    %% Relationships - JWT & Claims
    JWTService --> Claims : creates
    Claims --> UserRole : contains

    %% Relationships - Request/Response DTOs
    AuthService --> LoginRequest : accepts
    AuthService --> RegisterRequest : accepts
    AuthService --> AuthResponse : returns
    OrderService --> CreateOrderRequest : accepts
    OrderService --> OrderResponse : returns
```

## Penjelasan Arsitektur

### 1. **Domain Layer (Inti Bisnis)**
- **User**: Entity utama untuk pengguna dengan role-based access control
- **Order**: Entity untuk pesanan threat intelligence data
- **UserRole & OrderStatus**: Value objects untuk enum
- **OrderAggregate**: Domain aggregate untuk business logic order
- **Repository Interfaces**: Port untuk akses data

### 2. **Application Layer (Use Cases)**
- **AuthService**: Orchestrasi autentikasi dan autorisasi
- **OrderService**: Orchestrasi business logic untuk order management
- **Request/Response DTOs**: Data transfer objects untuk API

### 3. **Infrastructure Layer (Implementasi Eksternal)**
- **PostgresUserRepository & PostgresOrderRepository**: Implementasi repository dengan GORM
- **JWTService**: Service untuk JWT token management
- **RedisClient**: Client untuk caching dan session management
- **NewRelicMonitor**: Monitoring dan observability
- **Database**: Database connection dan migration

### 4. **Interface Layer (HTTP API)**
- **Handler**: HTTP request handlers untuk REST API
- **Middleware**: Cross-cutting concerns (auth, logging, CORS, rate limiting)
- **Router**: Route configuration dan setup

### 5. **Configuration**
- **Config**: Centralized configuration management
- **Various Config Structs**: Specific configuration untuk setiap komponen

## Prinsip Clean Architecture yang Diterapkan

1. **Dependency Inversion**: Infrastructure layer bergantung pada interfaces yang didefinisikan di domain layer
2. **Separation of Concerns**: Setiap layer memiliki tanggung jawab yang jelas
3. **Testability**: Business logic terisolasi dan mudah di-test
4. **Flexibility**: Mudah mengganti implementasi tanpa mengubah business logic
5. **Domain-Driven Design**: Domain entities sebagai pusat dari arsitektur

## Pola Desain yang Digunakan

- **Repository Pattern**: Abstraksi akses data
- **Aggregate Pattern**: Encapsulation business logic dalam OrderAggregate
- **Dependency Injection**: Loose coupling antar komponen
- **Middleware Pattern**: Cross-cutting concerns
- **Factory Pattern**: Creation objects dengan NewXXX functions