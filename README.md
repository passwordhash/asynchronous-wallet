# Asynchronous Wallet Service

Wallet service built with Go, providing reliable financial operations with thread-safety and transaction integrity.

---

##  Stack

- **Go 1.24**: Modern, concurrent programming language
- **Gin**: High-performance HTTP web framework
- **PostgreSQL**: Reliable, ACID-compliant database
- **PGX**: Optimized PostgreSQL driver for Go
- **Docker & Docker Compose**: Containerization and orchestration
- **Testing**: Unit tests with mocks and integration tests

## Architecture

The project follows clean architecture principles:

- **Handler Layer**: HTTP request handling using Gin
- **Service Layer**: Business logic implementation
- **Repository Layer**: Data access and persistence
- **Entity Layer**: Domain models and core business objects

## DB schema

```sql
CREATE TABLE wallets (
    id UUID PRIMARY KEY NOT NULL,
    balance BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

## Getting Started

### Requirements

- Docker and Docker Compose
- Go 1.24+ 
- Make

### Running with Docker

1. Clone the repository:
```bash
git clone https://github.com/passwordhash/asynchronous-wallet.git
cd asynchronous-wallet
```

2. Start the services:
```bash
make container-start
```
or 
```bash
docker compose -f ./docker-compose.yml --env-file config.env up -d
```

The application will be available at http://localhost:5782.

### Local Development

1. Clone the repository:
```bash
git clone https://github.com/passwordhash/asynchronous-wallet.git
cd asynchronous-wallet
```

2. Start infrastructure:
```bash
make container-infra
```

3. Run application:
```bash
make run
```

## Testing

The project includes comprehensive tests:

- **Unit Tests**: Test individual components in isolation using mocks
- **Integration Tests**: Test the complete flow of operations
- **Concurrency Tests**: Verify thread-safety under high load

Run unit tests with:
```bash
make unit-test
```

Run integration tests with (requires running the application in docker check [Running with Docker](#running-with-docker)
):
```bash
make integration-test
```

## API Endpoints

### Wallet Operations

- **POST /api/v1/wallet**
  - Deposit or withdraw funds
  - Request body: `{"wallet_id": "uuid", "amount": 100}` (positive for deposit, negative for withdrawal)

### Wallet Information

- **GET /api/v1/wallets/:id**
  - Get wallet balance
  - Returns: `{"balance": 100}`

