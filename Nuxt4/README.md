[На русском](./README.ru.md)

# Go 1.25 Backend Development

## 📋 Table of Contents

  - [Project Goal](https://www.google.com/search?q=%23project-goal)
  - [Key Features](https://www.google.com/search?q=%23key-features)
  - [Tech Stack](https://www.google.com/search?q=%23tech-stack)
  - [Project Architecture](https://www.google.com/search?q=%23project-architecture)
  - [Project Structure](https://www.google.com/search?q=%23project-structure)
  - [Database](https://www.google.com/search?q=%23database)
  - [Installation and Run](https://www.google.com/search?q=%23installation-and-run)
  - [Configuration](https://www.google.com/search?q=%23configuration)
  - [Development Stages](https://www.google.com/search?q=%23development-stages)

## 🎯 Project Goal

Building a high-performance and scalable server application to handle all game logic, interact with the database, and provide an API for the client side of the "Choose Me" game.

## 🚀 Key Features

  - **Client API:** Provision of a RESTful API for all client requests (registration, fetching world data, player actions).
  - **User Management:** Logic for registration, authorization, and session management with multi-device support.
  - **Clan System:** A full-fledged clan system with roles and custom names.
  - **Game Logic:**
      - Generation and management of the game world state.
      - Procedural character generation (based on JSON tags).
      - Handling of game events and player actions.
      - Calculation of decay mechanics, leveling, and NPC interactions.
  - **n8n Integration:** Interaction with the n8n service to trigger image generation workflows.
  - **Payment Processing:** Integration with payment gateways.

## 🛠 Tech Stack

### Core Technologies

| Technology | Version | Purpose |
|------------|--------|------------|
| **Go** | 1.25.1 | Main programming language |
| **Gin** | 1.10.1 | Web framework for building RESTful APIs |
| **PostgreSQL** | - | Primary relational database |
| **Redis** | 9.14.0 | In-memory DB for caching and sessions |
| **Docker** | - | Application containerization |

### Libraries and Dependencies

#### Database Interaction

  - **sqlx** (1.4.0) - Extension for database/sql with convenient methods.
  - **lib/pq** (1.10.9) - PostgreSQL driver.
  - **go-redis** (9.14.0) - Client for Redis.

#### Authentication and Security

  - **jwt/v5** (5.3.0) - JSON Web Tokens for authentication.
  - **uuid** (1.6.0) - Generation of unique identifiers.

#### Configuration and Logging

  - **Viper** (1.21.0) - Application configuration management.
  - **Logrus** (1.9.3) - Structured logging.
  - **godotenv** (1.5.1) - Loading environment variables from a .env file.

## 🏗 Project Architecture

The project is built on **Clean Architecture** principles, separated into three main layers:

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

  - **Repository Pattern** - Abstraction of data access.
  - **Dependency Injection** - Injecting dependencies via constructors.
  - **Middleware Pattern** - Processing requests through a middleware chain.
  - **Clean Architecture** - Separation into independent layers.

## 📁 Project Structure

```
Go1.25/
├── cmd/
│   └── main.go                      # Application entry point
│
├── pkg/                             # Main application code
│   ├── handler/                     # HTTP handlers (Gin routes)
│   │   ├── handler.go              # Routes initialization
│   │   ├── auth.go                 # Authentication (login, register)
│   │   ├── user_settings.go        # User settings
│   │   ├── middleware.go           # JWT middleware, CORS
│   │   └── response.go             # Standardized responses
│   │
│   ├── service/                     # Business logic
│   │   ├── service.go              # Services initialization
│   │   ├── auth.go                 # Authentication logic
│   │   └── user_settings.go        # User settings logic
│   │
│   └── repository/                  # Database interaction
│       ├── repository.go           # Repositories initialization
│       ├── postgres.go             # PostgreSQL connection
│       ├── redis.go                # Redis connection
│       ├── auth_postgres.go        # Auth repository
│       └── user_setting_postgres.go # Settings repository
│
├── configs/
│   └── config.yml                   # Application configuration
│
├── migrate/                         # Database migrations
│   ├── 000001_init.up.sql          # Table creation
│   └── 000001_init.down.sql        # Rollback migrations
│
├── Dockerfile                       # Multi-stage Docker build
├── .env                            # Environment variables (not in git)
├── .gitignore
├── go.mod                          # Project dependencies
├── go.sum                          # Dependencies checksums
├── server.go                       # HTTP server
├── user.go                         # User model
├── errors.go                       # Custom errors
└── README.md
```

## 🗄 Database

### PostgreSQL Schema

#### Main Tables

**users** - System users

```sql
- id (SERIAL PRIMARY KEY)
- email (VARCHAR UNIQUE)
- password_hash (VARCHAR)
```

**user\_refresh\_tokens** - Refresh tokens for different devices

```sql
- id (SERIAL PRIMARY KEY)
- user_id (INT FK → users)
- token (VARCHAR UNIQUE)
- expires_at (TIMESTAMPTZ)
- name_device (VARCHAR)
- device_info (VARCHAR)
```

**user\_settings** - User settings and profile

```sql
- user_id (INT PRIMARY KEY FK → users)
- name (VARCHAR)
- icon (VARCHAR)
- coin (INT)
- date_of_registration (TIMESTAMPTZ)
- paid_subscription (BOOLEAN)
- date_of_paid_subscription (TIMESTAMPTZ)
```

#### Clan System

**clan** - Clans

```sql
- id (SERIAL PRIMARY KEY)
- name (VARCHAR UNIQUE)
- description (TEXT)
- other (JSONB)
```

**roles** - System roles (1-5, where 1 is the highest)

```sql
- id (SMALLINT PRIMARY KEY)
- name (VARCHAR)
```

**clan\_members** - Clan participants

```sql
- clan_id (INT FK → clan)
- user_id (INT FK → users)
- role_id (SMALLINT FK → roles)
- PRIMARY KEY (clan_id, user_id)
```

**clan\_role\_names** - Custom role names for each clan

```sql
- clan_id (INT FK → clan)
- role_id (SMALLINT FK → roles)
- custom_name (VARCHAR)
- PRIMARY KEY (clan_id, role_id)
```

#### Game Entities

**cards** - Character cards

```sql
- id (SERIAL PRIMARY KEY)
- user_id (INT FK → users)
- name (VARCHAR)
- description (TEXT)
- other (JSONB)
```

**items** - Game items

```sql
- id (SERIAL PRIMARY KEY)
- name (VARCHAR)
- description (TEXT)
- HaveCard (BOOLEAN)
- other (JSONB)
```

### Redis

Used for:

  - **Caching** - Frequently requested data
  - **Sessions** - JWT tokens and refresh tokens
  - **Rate limiting** - Limiting request frequency

## 🚀 Installation and Run

### Prerequisites

  - Go 1.25.1 or higher
  - PostgreSQL 14+
  - Redis 7+
  - Docker and Docker Compose (optional)

### Local Development

1.  **Clone the repository**

<!-- end list -->

```bash
git clone <repository-url>
cd Go1.25
```

2.  **Install dependencies**

<!-- end list -->

```bash
go mod download
```

3.  **Configure environment variables**

Create a `.env` file in the project root:

```env
DB_PASSWORD=your_postgres_password
REDIS_PASSWORD=your_redis_password
JWT_SECRET=your_jwt_secret_key
```

4.  **Start PostgreSQL and Redis**

<!-- end list -->

```bash
# Using Docker Compose (recommended)
docker-compose up -d postgres redis
```

5.  **Apply migrations**

<!-- end list -->

```bash
# Use migrate CLI or execute SQL manually
psql -U postgres -d postgres -f migrate/000001_init.up.sql
```

6.  **Run the application**

<!-- end list -->

```bash
go run cmd/main.go
```

The server will start at `http://localhost:8080`

### Docker Deployment

The project uses a multi-stage Docker build to minimize image size.

1.  **Build the image**

<!-- end list -->

```bash
docker build -t go-game-backend:latest .
```

2.  **Run the container**

<!-- end list -->

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
  username: "postgres"    # PostgreSQL username
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

1.  **Environment Setup** - Installing Go, setting up the workspace, project initialization.
2.  **Architecture Design** - Defining project structure, modules, and database schema.
3.  **User API Implementation** - Developing endpoints for registration, authorization, and profile management.
4.  **Database Integration** - Configuring connections to PostgreSQL and Redis, implementing models and repositories.
5.  **Clan System** - Implementing the basic clan structure with roles.

### 🔄 In Progress

6.  **Core Game Logic Development** - Creating mechanisms for world and character generation, and simulating their life.
7.  **Gameplay API Creation** - Developing endpoints for interacting with the game world.

### 📋 Planned

8.  **n8n Integration** - Setting up interaction for image generation.
9.  **Payment Systems Integration** - Connecting payment gateways.
10. **Testing** - Writing unit and integration tests to verify the correctness of the API and game logic.
11. **Performance Optimization** - Profiling and optimizing bottlenecks.
12. **Deployment** - Preparing for deployment to the production server.

-----

## 📚 Additional Documentation

  - [DEVELOPMENT\_PLAN.md](https://www.google.com/search?q=./DEVELOPMENT_PLAN.md) - Detailed development plan and roadmap.
  - [API Documentation](https://www.google.com/search?q=./docs/API.md) - API endpoints documentation (in development).

## 🤝 Contribution

The project is under active development. When making changes, please follow the established architecture and design patterns.
