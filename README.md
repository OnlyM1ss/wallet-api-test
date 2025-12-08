# Wallet API Test

A production-ready REST API for user account and wallet management, built with Go and designed for financial transaction processing.

## Description

**Wallet API Test** is a modern web application that provides a comprehensive API for managing user accounts and digital wallets. The API handles user registration, authentication, and secure wallet operations including deposits and withdrawals. It's built with Go using the Gin web framework and PostgreSQL for data persistence, making it suitable for payment systems, e-wallet platforms, and financial applications.

The project demonstrates clean architecture principles, implements proper error handling, and includes comprehensive configuration management for different deployment environments.

## Features

- **User Management**
  - User registration with email and password authentication
  - Secure password handling
  - User profile retrieval
  - JWT-based authentication

- **Wallet Operations**
  - Create wallets for users
  - Deposit and withdraw funds
  - Real-time balance tracking
  - Transaction validation

- **System Reliability**
  - Health check endpoint for monitoring
  - Graceful shutdown handling
  - Connection pooling for database optimization
  - PostgreSQL for data persistence

- **Developer-Friendly**
  - RESTful API design following best practices
  - Comprehensive API documentation with examples
  - Docker and Docker Compose for easy setup
  - Environment-based configuration
  - Structured logging with Uber Zap

## Key Components

### 1. **Application Layer** (`internal/app/`)
- **App** - Main application orchestrator
  - Initializes database connections and dependency injection
  - Configures HTTP routes and middleware
  - Manages server startup and graceful shutdown
  - Creates database tables automatically

### 2. **Domain Layer** (`internal/domain/`)
- **Services** - Business logic layer
  - `UserService` - User creation, authentication, and retrieval
  - `WalletService` - Wallet operations and balance management
  - `Entities` - Domain models (User, Wallet)
  - `Repositories` - Data access interfaces

### 3. **Infrastructure Layer** (`internal/infrastructure/`)
- **Database** - PostgreSQL implementations
  - `UserRepository` - User CRUD operations
  - `WalletRepository` - Wallet CRUD operations
  - Connection management and SQL execution
- **HTTP** - Request/Response handlers
  - `UserHandler` - User endpoints (create, login, get)
  - `WalletHandler` - Wallet endpoints (create, process operations, get balance)

### 4. **Configuration** (`internal/config/`)
- **Config** - Environment-based settings
  - Server configuration (port, timeouts)
  - Database credentials (PostgreSQL)
  - Redis settings (optional)
  - JWT configuration
  - Log level settings

### 5. **Utilities** (`internal/pkg/`)
- **Logger** - Structured logging using Uber Zap
  - Configurable log levels
  - Structured logging with context

## How It Works (Basic Overview)

### Architecture Flow

```
HTTP Request
    ↓
[Gin Router] → [HTTP Handler]
                    ↓
            [Domain Service] → [Business Logic]
                    ↓
            [Repository] → [PostgreSQL Database]
```

### Typical Workflow

1. **User Registration**
   - Client sends email, username, and password via POST `/api/v1/users`
   - Handler validates input and calls UserService
   - Service creates user record in PostgreSQL with UUID
   - User record is returned with creation timestamp

2. **User Authentication**
   - Client sends email and password to POST `/api/v1/login`
   - Service validates credentials against database
   - JWT token is generated and returned with user info

3. **Wallet Operations**
   - Client creates wallet via POST `/api/v1/wallet/create`
   - Client performs deposits/withdrawals via POST `/api/v1/wallet`
   - Service validates operation (amount, sufficient balance for withdrawals)
   - Database transaction updates wallet balance
   - Response confirms successful operation

4. **Data Retrieval**
   - Client can retrieve user info via GET `/api/v1/users/:id`
   - Client can check wallet balance via GET `/api/v1/wallet/:walletId`
   - Health check available at GET `/health`

## API/Interfaces

### User Endpoints

| Method | Endpoint | Description | Request Body |
|--------|----------|-------------|--------------|
| POST | `/api/v1/users` | Create new user | `{ "email": "string", "username": "string", "password": "string" }` |
| POST | `/api/v1/login` | Authenticate user | `{ "email": "string", "password": "string" }` |
| GET | `/api/v1/users/:id` | Get user by UUID | (none) |

### Wallet Endpoints

| Method | Endpoint | Description | Request Body |
|--------|----------|-------------|--------------|
| POST | `/api/v1/wallet/create` | Create new wallet | `{ "userId": "uuid" }` |
| POST | `/api/v1/wallet` | Process operation | `{ "walletId": "uuid", "operationType": "DEPOSIT\|WITHDRAW", "amount": "int64" }` |
| GET | `/api/v1/wallet/:walletId` | Get wallet balance | (none) |

### System Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check and DB connectivity status |

### Response Format

**Success Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "username": "johndoe",
  "created_at": "2025-12-07T20:58:10Z"
}
```

**Error Response:**
```json
{
  "error": "Invalid email format"
}
```

## Requirements/Installation

### Prerequisites

- **Go**: 1.23 or higher
- **PostgreSQL**: 12 or higher
- **Redis**: 6 or higher (optional, for future enhancements)
- **Docker & Docker Compose** (optional, for containerized deployment)

### Dependencies

Core libraries used in the project:

| Library | Version | Purpose |
|---------|---------|---------|
| `github.com/gin-gonic/gin` | v1.9.1 | HTTP web framework |
| `github.com/lib/pq` | v1.10.9 | PostgreSQL driver |
| `github.com/jmoiron/sqlx` | v1.3.5 | Database utilities |
| `github.com/spf13/viper` | v1.16.0 | Configuration management |
| `go.uber.org/zap` | v1.24.0 | Structured logging |
| `github.com/google/uuid` | v1.3.0 | UUID generation |

### Installation Steps

#### Option 1: Local Development

1. **Clone the repository**
   ```bash
   git clone https://github.com/OnlyM1ss/wallet-api-test.git
   cd wallet-api-test
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up database**
   ```bash
   # Ensure PostgreSQL is running
   createdb mydb
   # Or use Docker Compose (Option 2 below)
   ```

4. **Configure environment variables**
   ```bash
   # Copy and modify config.env as needed
   export DB_HOST=localhost
   export DB_PORT=5432
   export DB_USER=postgres
   export DB_PASSWORD=postgres
   export DB_NAME=mydb
   export SERVER_PORT=8080
   ```

5. **Run the application**
   ```bash
   go run ./cmd/api/main.go
   ```

#### Option 2: Using Docker Compose (Recommended)

1. **Clone the repository**
   ```bash
   git clone https://github.com/OnlyM1ss/wallet-api-test.git
   cd wallet-api-test
   ```

2. **Start services**
   ```bash
   docker-compose up
   ```

   This will start:
   - **Application** on port 8080
   - **PostgreSQL** on port 5432
   - **Redis** on port 6379

3. **Verify health**
   ```bash
   curl http://localhost:8080/health
   ```

#### Building Docker Image

```bash
# Build the image
docker build -t wallet-api:latest .

# Run container
docker run -p 8080:8080 \
  -e DB_HOST=postgres \
  -e DB_PORT=5432 \
  -e DB_USER=postgres \
  -e DB_PASSWORD=postgres \
  -e DB_NAME=mydb \
  wallet-api:latest
```

## Usage Example

### Example 1: Basic User and Wallet Workflow

```bash
# 1. Check if service is running
curl -X GET http://localhost:8080/health

# Response:
# {"status":"healthy"}

# 2. Create a user
USER_RESPONSE=$(curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@example.com",
    "username": "alice",
    "password": "securepass123"
  }')

echo $USER_RESPONSE
# Response:
# {
#   "id": "550e8400-e29b-41d4-a716-446655440000",
#   "email": "alice@example.com",
#   "username": "alice",
#   "created_at": "2025-12-07T20:58:10Z"
# }

# 3. Login to get token
LOGIN_RESPONSE=$(curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@example.com",
    "password": "securepass123"
  }')

echo $LOGIN_RESPONSE
# Response includes JWT token and user info

# 4. Create wallet for the user
WALLET_RESPONSE=$(curl -X POST http://localhost:8080/api/v1/wallet/create \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "550e8400-e29b-41d4-a716-446655440000"
  }')

echo $WALLET_RESPONSE
# Response:
# {
#   "id": "660e8400-e29b-41d4-a716-446655440000",
#   "user_id": "550e8400-e29b-41d4-a716-446655440000",
#   "balance": 0,
#   "created_at": "2025-12-07T21:00:00Z"
# }

# 5. Deposit funds into wallet
curl -X POST http://localhost:8080/api/v1/wallet \
  -H "Content-Type: application/json" \
  -d '{
    "walletId": "660e8400-e29b-41d4-a716-446655440000",
    "operationType": "DEPOSIT",
    "amount": 5000
  }'

# Response:
# {
#   "message": "operation completed successfully",
#   "walletId": "660e8400-e29b-41d4-a716-446655440000",
#   "operationType": "DEPOSIT",
#   "amount": 5000
# }

# 6. Check wallet balance
curl -X GET http://localhost:8080/api/v1/wallet/660e8400-e29b-41d4-a716-446655440000

# Response:
# {
#   "id": "660e8400-e29b-41d4-a716-446655440000",
#   "user_id": "550e8400-e29b-41d4-a716-446655440000",
#   "balance": 5000,
#   "created_at": "2025-12-07T21:00:00Z",
#   "updated_at": "2025-12-07T21:05:00Z"
# }

# 7. Withdraw funds
curl -X POST http://localhost:8080/api/v1/wallet \
  -H "Content-Type: application/json" \
  -d '{
    "walletId": "660e8400-e29b-41d4-a716-446655440000",
    "operationType": "WITHDRAW",
    "amount": 1000
  }'

# 8. Verify withdrawal
curl -X GET http://localhost:8080/api/v1/wallet/660e8400-e29b-41d4-a716-446655440000
# Balance should now be 4000
```

### Example 2: Python Script Integration

```python
import requests
import json
from pprint import pprint

BASE_URL = "http://localhost:8080"

# Create user
user_data = {
    "email": "bob@example.com",
    "username": "bob",
    "password": "password456"
}

response = requests.post(f"{BASE_URL}/api/v1/users", json=user_data)
user = response.json()
user_id = user["id"]

print("Created User:")
pprint(user)

# Create wallet
wallet_data = {"userId": user_id}
response = requests.post(f"{BASE_URL}/api/v1/wallet/create", json=wallet_data)
wallet = response.json()
wallet_id = wallet["id"]

print("\nCreated Wallet:")
pprint(wallet)

# Deposit
deposit_data = {
    "walletId": wallet_id,
    "operationType": "DEPOSIT",
    "amount": 10000
}

response = requests.post(f"{BASE_URL}/api/v1/wallet", json=deposit_data)
print("\nDeposit Response:")
pprint(response.json())

# Get wallet balance
response = requests.get(f"{BASE_URL}/api/v1/wallet/{wallet_id}")
print("\nWallet Balance:")
pprint(response.json())
```

## Project Structure

```
wallet-api-test/
├── cmd/
│   └── api/
│       └── main.go                    # Application entry point
├── internal/
│   ├── app/
│   │   └── app.go                     # Main application orchestrator
│   ├── config/
│   │   └── config.go                  # Configuration management
│   ├── domain/
│   │   ├── entities/                  # Domain models
│   │   ├── repositories/              # Repository interfaces
│   │   ├── services/                  # Business logic
│   │   └── user/                      # User domain
│   ├── infrastructure/
│   │   ├── database/postgres/         # PostgreSQL implementations
│   │   └── http/handlers/             # HTTP handlers
│   └── pkg/
│       └── logger/                    # Logging utilities
├── migrations/                        # Database migrations
├── Dockerfile                         # Docker image definition
├── docker-compose.yml                 # Docker Compose configuration
├── config.env                         # Environment configuration
├── go.mod                             # Go module definition
├── go.sum                             # Dependency checksums
├── API_EXAMPLES.md                    # Detailed API examples
└── README.md                          # This file
```

## Configuration

The application uses environment variables and YAML configuration files for settings. 

### Environment Variables

```bash
# Server
SERVER_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=mydb
DB_SSLMODE=disable

# Redis (optional)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT
JWT_SECRET_KEY=your-secret-key-change-in-production
JWT_EXPIRES_IN=3600

# Logging
LOG_LEVEL=info
```

### Configuration File (config.env)

The application can also read settings from `config.env` file in YAML format:

```yaml
server:
  port: "8080"
  readTimeout: 30
  writeTimeout: 30
  idleTimeout: 120

database:
  host: "localhost"
  port: "5432"
  user: "postgres"
  password: "postgres"
  name: "mydb"
  sslMode: "disable"
```

## API Documentation

For detailed API endpoint documentation with curl and PowerShell examples, see [API_EXAMPLES.md](./API_EXAMPLES.md).

The documentation includes:
- Health check endpoint
- User CRUD operations
- Authentication (login)
- Wallet creation and balance retrieval
- Wallet operations (deposit/withdraw)
- Complete workflow example
- PowerShell integration examples

## Development

### Project Structure Principles

This project follows **Clean Architecture** principles:

- **Domain Layer** - Business logic independent of frameworks
- **Application Layer** - Orchestration and dependency injection
- **Infrastructure Layer** - Database and HTTP implementations
- **External Interfaces** - HTTP handlers and routes

### Adding New Features

1. Define domain models in `internal/domain/entities/`
2. Create repository interfaces in `internal/domain/repositories/`
3. Implement business logic in `internal/domain/services/`
4. Add database implementation in `internal/infrastructure/database/`
5. Create HTTP handlers in `internal/infrastructure/http/handlers/`
6. Register routes in `internal/app/app.go`

## Deployment

### Docker Deployment

```bash
# Build and push to registry
docker build -t myregistry/wallet-api:1.0.0 .
docker push myregistry/wallet-api:1.0.0

# Run in production
docker run -d \
  -p 8080:8080 \
  -e DB_HOST=postgres-server \
  -e DB_PORT=5432 \
  -e DB_USER=app_user \
  -e DB_PASSWORD=$DB_PASSWORD \
  -e DB_NAME=wallet_db \
  -e JWT_SECRET_KEY=$JWT_SECRET \
  --name wallet-api \
  myregistry/wallet-api:1.0.0
```

### Kubernetes Deployment

Create a `deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: wallet-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: wallet-api
  template:
    metadata:
      labels:
        app: wallet-api
    spec:
      containers:
      - name: wallet-api
        image: myregistry/wallet-api:1.0.0
        ports:
        - containerPort: 8080
        env:
        - name: DB_HOST
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: db_host
```

## Contributing

Contributions are welcome! Please follow these guidelines:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Style

- Follow Go conventions and idioms
- Use `gofmt` for code formatting
- Write clear, concise comments
- Add tests for new features

## License

This project is provided as-is for educational and testing purposes.

## Support & Documentation

- **API Documentation**: See [API_EXAMPLES.md](./API_EXAMPLES.md) for detailed endpoint examples
- **Configuration**: See configuration section above
- **Docker Setup**: Use `docker-compose up` for quick local development

## Troubleshooting

### Database Connection Issues

**Error**: `Failed to connect to database`

**Solution**:
- Ensure PostgreSQL is running: `docker-compose up postgres`
- Check credentials in `config.env`
- Verify database exists: `createdb mydb`

### Port Already in Use

**Error**: `Address already in use`

**Solution**:
```bash
# Change port in config.env
SERVER_PORT=8081

# Or kill process using port 8080
lsof -ti:8080 | xargs kill -9
```

### JWT Token Invalid

**Error**: `Invalid token`

**Solution**:
- Ensure `JWT_SECRET_KEY` is set correctly
- Check token expiration time
- Regenerate token with login endpoint

## Future Enhancements

- [ ] Redis caching for performance optimization
- [ ] Transaction history tracking
- [ ] Advanced authentication (OAuth2, 2FA)
- [ ] Webhook system for wallet events
- [ ] Analytics and reporting dashboard
- [ ] Rate limiting and API versioning improvements
