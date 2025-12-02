# Insider Message API

A message management and automation service built with Go. This application allows you to create messages, manage automated message sending, and track sent messages.

## Features

- 📨 Create and manage messages
- ⚡ Automated message sending with configurable intervals
- 🔄 Toggle automation on/off via API
- 📊 Track all sent messages
- 🗄️ PostgreSQL for persistent storage
- ⚡ Redis for caching

## Tech Stack

- **Language:** Go 1.24
- **Database:** PostgreSQL 16
- **Cache:** Redis 7
- **Containerization:** Docker & Docker Compose

## Project Structure
```text
insider-message/
├── api/
│   ├── cache/          # Redis cache implementation
│   ├── database/       # PostgreSQL connection
│   ├── handler/        # HTTP handlers
│   ├── job/            # Background job/automation
│   ├── model/          # Data models
│   ├── repository/     # Data access layer
│   └── service/        # Business logic
├── cmd/
│   ├── migrate/        # Database migration tool
│   └── server/         # Main application entry
├── config/             # Configuration loader
├── migrations/         # SQL migration files
├── config.yaml         # Application configuration
├── docker-compose.yaml
├── Dockerfile
└── swagger.yaml        # API documentation
```

## Prerequisites

- Docker & Docker Compose
- Go 1.24+ (for local development)

## Quick Start with Docker

### 1. Clone the repository

```bash
git clone https://github.com/kayaramazan/insider-message.git
cd insider-message
```
### 2. Run database migrations

```bash
docker-compose run --rm migrate
```

### 3. Start all services

```bash
docker-compose up -d
```

This will start:
- PostgreSQL on port `5432`
- Redis on port `6379`
- Application on port `8080`



### 4. Verify the application is running

```bash
curl http://localhost:8080/api/messages
```

## Local Development

### 1. Start dependencies

```bash
docker-compose up -d postgres redis
```

### 2. Run migrations

```bash
go run cmd/migrate/main.go up
```

### 3. Start the application

```bash
go run cmd/server/main.go
```

## Configuration

Configuration is managed via `config.yaml`:


### Environment Variables

You can override configuration using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | PostgreSQL user | `postgres` |
| `DB_PASSWORD` | PostgreSQL password | `Aa123456` |
| `DB_NAME` | Database name | `postgres` |
| `DB_SSLMODE` | SSL mode | `disable` |
| `REDIS_HOST` | Redis host | `localhost` |
| `REDIS_PORT` | Redis port | `6379` |
| `REDIS_PASSWORD` | Redis password | `` |

## API Endpoints

### Toggle Automation

```bash
curl -X PUT http://localhost:8080/api/automation/toggle
```

**Response:**
```json
{
  "Status": "Accepted",
  "Automation Status": true
}
```

### Create Message

```bash
curl -X POST http://localhost:8080/api/message \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Hello, this is a test message!",
    "phone": "+905551234567"
  }'
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "content": "Hello, this is a test message!",
  "phone": "+905551234567",
  "status": 1,
  "created_at": "2025-12-02T10:30:00Z"
}
```

### Get All Sent Messages

```bash
curl http://localhost:8080/api/messages
```

**Response:**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "content": "Hello, this is a test message!",
    "phone": "+905551234567",
    "status": 2,
    "created_at": "2025-12-02T10:30:00Z"
  }
]
```

## Message Status

| Status Code | Description |
|-------------|-------------|
| `1` | Pending |
| `2` | Sent |

## API Documentation

Swagger documentation is available in `swagger.yaml`. You can view it using:

- [Swagger Editor](https://editor.swagger.io/) - Paste the content of `swagger.yaml`
- [Swagger UI](https://swagger.io/tools/swagger-ui/)

## Database Migrations

### Run migrations up

```bash
# With Docker
docker-compose run --rm migrate

# Locally
go run cmd/migrate/main.go up
```

### Rollback migrations

```bash
go run cmd/migrate/main.go down
```

## Troubleshooting

### Container won't start

```bash
# Check logs
docker-compose logs app

# Rebuild containers
docker-compose up -d --build
```

### Database connection issues

```bash
# Check if PostgreSQL is healthy
docker-compose ps

# Check PostgreSQL logs
docker-compose logs postgres
```

### Reset everything

```bash
docker-compose down -v
docker-compose up -d
docker-compose run --rm migrate
```

## License

MIT License
```

Bu README dosyası şunları içeriyor:

- ✅ Proje açıklaması ve özellikler
- ✅ Teknoloji stack'i
- ✅ Proje yapısı
- ✅ Docker ile hızlı başlangıç
- ✅ Tüm Docker komutları tablosu
- ✅ Lokal geliştirme adımları
- ✅ Yapılandırma detayları
- ✅ Environment variables
- ✅ API endpoint'leri ve örnekler
- ✅ Migration komutları
- ✅ Troubleshooting bölümü
