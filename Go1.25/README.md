# Backend Development with Go 1.25

## 📋 Table of Contents

- [Project Goal](#project-goal)
- [Key Features](#key-features)
- [Technology Stack](#technology-stack)
- [Project Architecture](#project-architecture)
- [Project Structure](#project-structure)
- [Database](#database)
- [Installation and Setup](#installation-and-setup)
- [Configuration](#configuration)
- [Development Stages](#development-stages)

## 🎯 Project Goal

Creating a high-performance and scalable server application to handle all game logic, database interactions, and provide API for the "Choose Me" game client.

## 🚀 Key Features

- **Client API:** Providing RESTful API for all client requests (registration, world data retrieval, player actions)
- **User Management:** Logic for registration, authorization, session management with multi-device support
- **Clan System:** Full-featured clan system with roles and custom names
- **Game Logic:**
  - Generation and management of game world state
  - Procedural character generation (based on JSON tags)
  - Processing game events and player actions
  - Calculation of decay mechanics, leveling, NPC interactions
- **n8n Integration:** Interaction with n8n service to trigger image generation workflows
- **Payment Processing:** Integration with payment gateways

## 🛠 Technology Stack

### Core Technologies

| Technology     | Version | Purpose                             |
| -------------- | ------- | ----------------------------------- |
| **Go**         | 1.25.1  | Main programming language           |
| **Gin**        | 1.10.1  | Web framework for creating RESTful API |
| **PostgreSQL** | -       | Main relational database            |
| **Redis**      | 9.14.0  | In-memory DB for caching and sessions |
| **Docker**     | -       | Application containerization        |

### Libraries and Dependencies

#### Database Operations
- **sqlx** (1.4.0) - Extension for database/sql with convenient methods
- **lib/pq** (1.10.9) - PostgreSQL driver
- **go-redis** (9.14.0) - Redis client

#### Authentication and Security
- **jwt/v5** (5.3.0) - JSON Web Tokens for authentication
- **uuid** (1.6.0) - Unique identifier generation

#### Configuration and Logging
- **Viper** (1.21.0) - Application configuration management
- **Logrus** (1.9.3) - Structured logging
- **godotenv** (1.5.1) - Loading environment variables from .env file

## 🏗 Project Architecture

The project is built on **Clean Architecture** principles with separation into three main layers:

```
┌─────────────────────────────────────────┐
│         HTTP Layer (Gin)                │
│  ┌───────────────────────────────────┐  │
│  │         Handlers                  │  │
│  │  • auth.go                        │  │
│  │  • user_settings.go               │  │
│  │  • middleware.go                  │  │
│  └───────────────────────────────────┘  │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│         Business Logic Layer            │
│  ┌───────────────────────────────────┐  │
│  │         Services                  │  │
│  │  • auth.go                        │  │
│  │  • user_settings.go               │  │
│  └───────────────────────────────────┘  │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│         Data Access Layer               │
│  ┌───────────────────────────────────┐  │
│  │       Repositories                │  │
│  │  • auth_postgres.go               │  │
│  │  • user_setting_postgres.go       │  │
│  │  • redis.go                       │  │
│  └───────────────────────────────────┘  │
└─────────────────┬───────────────────────┘
                  │
        ┌─────────┴─────────┐
        ▼                   ▼
  ┌──────────┐        ┌──────────┐
  │PostgreSQL│        │  Redis   │
  └──────────┘        └──────────┘
```

### Design Patterns

- **Repository Pattern** - Data access abstraction
- **Dependency Injection** - Dependency injection through constructors
- **Middleware Pattern** - Request processing through middleware chain
- **Clean Architecture** - Separation into independent layers

## 📁 Project Structure

```
Go1.25/
├── cmd/
│   └── main.go                      # Application entry point
│
├── pkg/                             # Main application code
│   ├── handler/                     # HTTP handlers (Gin routes)
│   │   ├── handler.go              # Route initialization
│   │   ├── auth.go                 # Authentication (login, register)
│   │   ├── user_settings.go        # User settings
│   │   ├── middleware.go           # JWT middleware, CORS
│   │   └── response.go             # Standardized responses
│   │
│   ├── service/                     # Business logic
│   │   ├── service.go              # Service initialization
│   │   ├── auth.go                 # Authentication logic
│   │   └── user_settings.go        # User settings logic
│   │
│   └── repository/                  # Database operations
│       ├── repository.go           # Repository initialization
│       ├── postgres.go             # PostgreSQL connection
│       ├── redis.go                # Redis connection
│       ├── auth_postgres.go        # Authentication repository
│       └── user_setting_postgres.go # Settings repository
│
├── configs/
│   └── config.yml                   # Application configuration
│
├── migrate/                         # Database migrations
│   ├── 000001_init.up.sql          # Table creation
│   └── 000001_init.down.sql        # Migration rollback
│
├── Dockerfile                       # Multi-stage Docker build
├── .env                            # Environment variables (not in git)
├── .gitignore
├── go.mod                          # Project dependencies
├── go.sum                          # Dependency checksums
├── server.go                       # HTTP server
├── user.go                         # User model
├── errors.go                       # Custom errors
└── README.md
```

## Database

### [PostgreSQL](SCHEME_POSTRESQL.md)

### Redis

Used for:
- **Caching** - Frequently requested data
- **Sessions** - JWT tokens and refresh tokens
- **Rate limiting** - Request frequency limiting

## 🚀 Installation and Setup

### Prerequisites

- Go 1.25.1 or higher
- PostgreSQL 14+
- Redis 7+
- Docker and Docker Compose

### Local Development

1. **Clone repository**
```bash
git clone <repository-url>
cd Go1.25
```

2. **Install dependencies**
```bash
go mod download
```

3. **Configure environment variables**

Create `.env` file in project root:
```env
DB_PASSWORD=your_postgres_password
REDIS_PASSWORD=your_redis_password
JWT_SECRET=your_jwt_secret_key
```

4. **Start PostgreSQL and Redis**
```bash
# Using Docker Compose (recommended)
docker-compose up -d postgres redis
```

5. **Apply migrations**
```bash
# Use migrate CLI or execute SQL manually
psql -U postgres -d postgres -f migrate/000001_init.up.sql
```

6. **Run application**
```bash
go run cmd/main.go
```

Server will start on `http://localhost:8080`

### Docker Deployment

The project uses multi-stage Docker build to minimize image size.

1. **Build image**
```bash
docker build -t go-game-backend:latest .
```

2. **Run container**
```bash
docker run -d \
  --name game-backend \
  -p 8080:8080 \
  --env-file .env \
  go-game-backend:latest
```

## ⚙️ Configuration

### config.yml

```yaml
port: "8080"              # HTTP server port

db:
  username: "postgres"    # PostgreSQL user
  host: "localhost"       # PostgreSQL host
  port: "5432"           # PostgreSQL port
  database: "postgres"    # Database name
  sslmode: "disable"     # SSL mode

redis:
  addr: "localhost:6379" # Redis address
  db: 0                  # Redis database number
```

### Environment Variables (.env)

```env
DB_PASSWORD=          # PostgreSQL password
REDIS_PASSWORD=       # Redis password (if set)
JWT_SECRET=           # Secret key for JWT
```

## 📝 Development Stages

### ✅ Completed

1. **Environment Setup** - Go installation, workspace setup, project initialization
2. **Architecture Design** - Project structure definition, modules and database schema
3. **User API Implementation** - Endpoint development for registration, authorization and profile management
4. **Database Integration** - PostgreSQL and Redis connection setup, models and repositories implementation
5. **Clan System** - Basic clan structure implementation with roles

### 🔄 In Progress

6. **Game Logic Core Development** - Creating world and character generation mechanisms, simulating their life
7. **Game Process API Creation** - Developing endpoints for game world interaction

### 📋 Planned

8. **n8n Integration** - Setting up interaction for image generation
9. **Payment System Integration** - Connecting payment gateways
10. **Testing** - Writing unit and integration tests to verify API and game logic correctness
11. **Performance Optimization** - Profiling and optimizing bottlenecks
12. **Deployment** - Preparing for production server deployment

---

## 📚 Additional Documentation

- [DEVELOPMENT_PLAN.md](./DEVELOPMENT_PLAN.md) - Detailed development plan and roadmap
- [API Documentation](./docs/API.md) - API endpoints documentation (in development)

## 🤝 Contributing

The project is in active development. When making changes, follow the established architecture and design patterns.

## 📄 License

[Specify project license]
