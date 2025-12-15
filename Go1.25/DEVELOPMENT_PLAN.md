[На русском](./DEVELOPMENT_PLAN.ru.md)

# Backend Development Plan (Go 1.25)

## 📊 Current Project State

### Overall Readiness: ~35%

The project is in the initial development stage. The basic infrastructure, authentication system, and core data models have been implemented.

## ✅ Implemented Features

### 1. Infrastructure and Architecture

- ✅ **Clean Architecture** - Separation into layers: handler → service → repository
- ✅ **Dependency Injection** - Dependency injection via constructors
- ✅ **Configuration** - Viper for settings management (config.yml + .env)
- ✅ **Logging** - Structured logging via Logrus (JSON format)
- ✅ **Docker** - Multi-stage Dockerfile for production deployment
- ✅ **DB Migrations** - Migration system for PostgreSQL

### 2. Authentication System

- ✅ **User Registration** - Endpoint for creating new accounts
- ✅ **Authorization** - Login with JWT token issuance
- ✅ **JWT Authentication** - Access and Refresh tokens
- ✅ **Multi-device support** - Support for tokens across different devices
- ✅ **Middleware** - JWT middleware for endpoint protection
- ✅ **Password hashing** - Secure password storage

### 3. User Management

- ✅ **User Settings** - CRUD operations for user settings
- ✅ **User Profile** - Name, icon, coin balance
- ✅ **Premium Subscription** - Support for paid subscriptions with dates

### 4. Database

- ✅ **PostgreSQL** - Main DB with full schema
- ✅ **Redis** - Caching and session storage
- ✅ **Clan Schema** - Tables for the clan system with roles
- ✅ **Card Schema** - Tables for game cards and items
- ✅ **Triggers** - Automatic `updated_at` updates

### 5. API Endpoints

#### Implemented Endpoints:

**Authentication:**
- `POST /auth/register` - Registration
- `POST /auth/login` - Login
- `POST /auth/refresh` - Refresh token

**User Settings:**
- `GET /api/user/settings` - Get settings
- `PUT /api/user/settings` - Update settings

## 🔄 In Progress

### 1. Game Logic (0%)

**Priority: HIGH**

- [ ] **Game World**
  - [ ] Models for locations and the world
  - [ ] Procedural world generation
  - [ ] API for retrieving world information
  
- [ ] **Characters (NPCs)**
  - [ ] Procedural generation based on JSON tags
  - [ ] Attributes and stats system
  - [ ] Character life simulation
  - [ ] API for interaction with characters

- [ ] **Events and Quests**
  - [ ] Game event system
  - [ ] Quest generation
  - [ ] Player choice processing
  - [ ] Calculation of consequences

### 2. Clan System (20%)

**Priority: MEDIUM**

- [x] DB Schema for clans
- [ ] API for creating clans
- [ ] API for managing members
- [ ] API for managing roles
- [ ] Permission system
- [ ] Clan chat (possibly via WebSocket)

### 3. Gacha Mechanics (0%)

**Priority: HIGH**

- [ ] Card rarity system
- [ ] Drop probabilities
- [ ] API for opening packs
- [ ] Card inventory
- [ ] Duplicate system

## 📋 Planned Features

### Phase 1: Game Core (2-3 months)

#### 1.1 Game World and Characters

**Goal:** Create a living game world with procedurally generated characters

**Tasks:**
- Develop a procedural world generation system
- Implement a character generator based on JSON tags
- Create an NPC life simulation system
- Develop mechanics for player interaction with the world

**Technical Details:**
- Use JSONB in PostgreSQL for storing dynamic data
- Implement caching for active characters in Redis
- Create background workers for world simulation

#### 1.2 Events and Quests

**Goal:** Implement a system of dynamic events and quests

**Tasks:**
- Create an event engine
- Develop a system of choices and consequences
- Implement character leveling mechanics
- Create a decay and degradation system

#### 1.3 Gacha System

**Goal:** Implement mechanics for obtaining new characters

**Tasks:**
- Develop a rarity and probability system
- Create an API for opening packs
- Implement inventory and card management
- Add a card trading/selling system

### Phase 2: Integrations (1-2 months)

#### 2.1 n8n Integration

**Goal:** Automation of image generation for characters

**Tasks:**
- Configure webhook endpoints for n8n
- Create a queue system for image generation
- Implement processing of results from n8n
- Add fallback mechanisms

**Technical Details:**
- Use Redis for queues
- Implement retry logic
- Add generation status monitoring

#### 2.2 Payment System

**Goal:** Integration with payment gateways for monetization

**Tasks:**
- Select a payment gateway (Stripe, PayPal, or local)
- Implement webhook handlers for payments
- Create a system for purchasing premium subscriptions
- Implement in-game currency purchases
- Add transaction history

**Security:**
- Webhook signature validation
- Idempotency of operations
- Logging of all transactions

### Phase 3: Optimization and Scaling (1 month)

#### 3.1 Performance

**Tasks:**
- Application profiling (pprof)
- SQL query optimization
- DB index tuning
- Redis operation optimization
- Connection pooling implementation

#### 3.2 Caching

**Tasks:**
- Caching strategy for different data types
- Cache invalidation mechanisms
- Read-through cache implementation
- Hit rate monitoring

#### 3.3 Rate Limiting

**Tasks:**
- Rate limiting middleware implementation
- DDoS protection
- Throttling for expensive operations

### Phase 4: Testing and Quality (ongoing)

#### 4.1 Unit Tests

**Goal:** Minimum 70% test coverage

**Tasks:**
- Tests for all service layers
- Tests for repository layers
- Dependency mocking
- Use testify for assertions

#### 4.2 Integration Tests

**Tasks:**
- API endpoint tests
- DB operation tests
- Redis operation tests
- E2E tests for critical flows

#### 4.3 CI/CD

**Tasks:**
- GitHub Actions configuration
- Automatic test execution
- Linting (golangci-lint)
- Automatic deploy to staging

### Phase 5: Monitoring and Observability (1 month)

#### 5.1 Metrics

**Tasks:**
- Prometheus integration
- Custom metrics creation
- Grafana dashboard setup
- Alerts for critical metrics

#### 5.2 Tracing

**Tasks:**
- OpenTelemetry integration
- Distributed tracing for requests
- Performance profiling

#### 5.3 Logging

**Tasks:**
- Centralized logging (ELK stack)
- Structured logs
- Log levels management
- Correlation of logs with traces

## 🔧 Technical Debt and Improvements

### High Priority

- [ ] **Input Validation** - Add validation for all endpoints
- [ ] **Error handling** - Standardize error processing
- [ ] **API Documentation** - Create OpenAPI/Swagger specification
- [ ] **Graceful shutdown** - Correct server shutdown handling
- [ ] **Health checks** - Endpoints for service health checks

### Medium Priority

- [ ] **Request ID** - Add request ID for tracing
- [ ] **CORS configuration** - Proper CORS setup for production
- [ ] **Pagination** - Implement pagination for lists
- [ ] **Sorting & Filtering** - Add sorting and filtering capabilities
- [ ] **API Versioning** - Versioning of API endpoints

### Low Priority

- [ ] **GraphQL** - Consider adding a GraphQL API
- [ ] **WebSocket** - For real-time updates
- [ ] **gRPC** - For internal services communication
- [ ] **Cache Warming** - Automatic cache warming on startup

## 📈 Success Metrics

### Performance

- Response time < 100ms for 95% of requests
- Throughput > 1000 RPS
- Database query time < 50ms
- Redis operations < 10ms

### Reliability

- Uptime > 99.9%
- Error rate < 0.1%
- Zero data loss
- Successful recovery from failures

### Code Quality

- Test coverage > 70%
- Zero critical security vulnerabilities
- Code review for all changes
- Documentation for all public APIs

## 🗓 Estimated Timeline

| Phase | Duration | Status |
|------|------|--------|
| **Infrastructure** | Month 1 | ✅ Completed |
| **Authentication** | Month 1-2 | ✅ Completed |
| **Game Core** | Month 3-5 | 🔄 In Progress |
| **Integrations** | Month 6-7 | 📋 Planned |
| **Optimization** | Month 8 | 📋 Planned |
| **Testing** | Ongoing | 🔄 In Progress |
| **Monitoring** | Month 9 | 📋 Planned |
| **Production Ready** | Month 10 | 🎯 Goal |

## 🎯 Next Steps (Next Sprint)

### Sprint Goals

1. **Start Game Core Development**
   - Create models for the game world
   - Implement basic character generation
   - Develop API for retrieving world info

2. **Improve Existing Code**
   - Add input validation
   - Write unit tests for auth service
   - Create OpenAPI documentation

3. **Prepare Infrastructure**
   - Configure CI/CD pipeline
   - Add health check endpoints
   - Implement graceful shutdown

### Specific Tasks

- [ ] Create `World` and `Location` models
- [ ] Implement `CharacterGenerator` service
- [ ] Add endpoints for the game world
- [ ] Write tests for `AuthService`
- [ ] Create Swagger documentation
- [ ] Configure GitHub Actions
- [ ] Add `/health` and `/ready` endpoints

---

**Last Update:** 28.11.2025

**Document Version:** 1.0