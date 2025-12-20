[На русском](./DEVELOPMENT_PLAN.ru.md)

# "Choose Me" Game Development Plan - Frontend

## Project Overview

Development of the frontend part of the "Choose Me" game based on Nuxt 4 using the modern Vue.js ecosystem tech stack.

## Application Architecture

### Tech Stack

#### Core Technologies

  - **Nuxt 4** (v4.1.2) - Meta-framework for Vue.js
  - **Vue 3** (v3.5.18) - Composition API, Reactivity System
  - **TypeScript** (v5.6.3) - Strict typing
  - **Vite** - Fast build and HMR

#### State Management

  - **Pinia** (v3.0.3) - Official state manager
  - **pinia-plugin-persistedstate** - State persistence in localStorage

#### UI/UX

  - **Nuxt UI** (v4.0.0) - Ready-made components
  - **Iconify** (Lucide icons) - Icons
  - **Custom CSS** - Custom styles

#### Validation & Types

  - **Zod** (v4.1.11) - Runtime validation and schemas

#### Development Tools

  - **ESLint** + **Prettier** - Code quality
  - **Nuxt DevTools** - Debugging
  - **Bun** - Package manager

### Application Structure

```
app/
├── assets/          # Static assets (CSS, images)
├── components/      # Reusable Vue components
├── composables/     # Composables (logic, hooks)
├── layouts/         # Page layouts
├── pages/           # Pages (file-based routing)
├── stores/          # Pinia stores
└── middleware/      # Routing middleware (planned)
```

## Detailed Development Plan

### Phase 1: Infrastructure and Basic Setup ✅

**Status:** Completed

**Completed Tasks:**

  - [x] Nuxt 4 project initialization
  - [x] TypeScript configuration setup
  - [x] Nuxt UI installation and setup
  - [x] Pinia setup with persistence
  - [x] ESLint and Prettier configuration
  - [x] Zod setup for validation
  - [x] API proxy setup
  - [x] Docker configuration

**Result:**

  - Fully configured development environment
  - Ready project structure
  - Configured code quality tools

-----

### Phase 2: Basic Architecture ✅

**Status:** Completed

**Completed Tasks:**

  - [x] Creation of basic directory structure
  - [x] Pinia stores setup (user, token)
  - [x] Creation of `useApiFetch` composable
  - [x] Default layout creation
  - [x] AppHeader component creation
  - [x] File-based routing setup

**Components:**

#### Stores

1.  **token.ts** - JWT token management

      - Storage of access/refresh tokens
      - Automatic token refresh
      - Persistence in localStorage

2.  **user.ts** - User data management

      - User profile
      - Balance and currency
      - Persistence in localStorage

#### Composables

1.  **useApiFetch.ts** - Wrapper over $fetch
      - Automatic token injection
      - Error handling
      - Type-safe requests

**Result:**

  - Ready architecture for API interaction
  - State management system
  - Reusable logic

-----

### Phase 3: Authentication 🔄

**Status:** In Progress

**Current Progress:**

  - [x] `/login` page
  - [x] LoginForm component
  - [x] RegisterForm component
  - [x] Basic form validation
  - [ ] Backend API integration
  - [ ] Token handling
  - [ ] Middleware for route protection
  - [ ] Auth error handling

**Tasks:**

#### 3.1 API Integration

  - [ ] Connect `/auth/login` endpoint
  - [ ] Connect `/auth/register` endpoint
  - [ ] Connect `/auth/refresh` endpoint
  - [ ] Server response handling
  - [ ] Saving tokens to store

#### 3.2 Route Protection

  - [ ] Create `auth.ts` middleware
  - [ ] Check token existence
  - [ ] Automatic redirect to `/login`
  - [ ] Check token validity

#### 3.3 UX Improvements

  - [ ] Loading states
  - [ ] Error handling (toast notifications)
  - [ ] Form validation with Zod
  - [ ] Automatic login after registration

**Dependencies:**

  - Backend API endpoints for auth
  - Defined JWT token structure

-----

### Phase 4: User Profile 📋

**Status:** Planned

**Tasks:**

#### 4.1 Profile Page

  - [ ] Create `/profile/[id]` dynamic route
  - [ ] Profile information display
  - [ ] Avatar component
  - [ ] Player statistics display

#### 4.2 Profile Editing

  - [ ] Profile edit form
  - [ ] Avatar upload
  - [ ] Data validation
  - [ ] Saving changes via API

#### 4.3 Balance and Currency

  - [ ] Balance display component
  - [ ] Transaction history
  - [ ] History filtering
  - [ ] Pagination

#### 4.4 Settings

  - [ ] Settings page
  - [ ] Notification settings
  - [ ] Privacy settings
  - [ ] Password change
  - [ ] Account deletion

**Components to create:**

  - `ProfileCard.vue` - Profile card
  - `ProfileEditor.vue` - Profile editor
  - `BalanceWidget.vue` - Balance widget
  - `TransactionHistory.vue` - Transaction history
  - `SettingsForm.vue` - Settings form

**API endpoints:**

  - `GET /api/profile/:id` - Get profile
  - `PUT /api/profile/:id` - Update profile
  - `GET /api/transactions` - Transaction history
  - `PUT /api/settings` - Update settings

-----

### Phase 5: Game Interface 📋

**Status:** Started (basic page)

**Tasks:**

#### 5.1 Game World

  - [x] Basic `/game` page
  - [ ] World map component
  - [ ] Location visualization
  - [ ] World navigation
  - [ ] Transition animations

#### 5.2 Location System

  - [ ] Location component
  - [ ] Location details
  - [ ] Available actions
  - [ ] Characters in location

#### 5.3 Character System

  - [ ] Character card
  - [ ] Detailed information
  - [ ] Character inventory
  - [ ] Characteristics and stats

#### 5.4 Events and Quests

  - [ ] Event component
  - [ ] Choice options display
  - [ ] Decision-making system
  - [ ] Consequences of choices
  - [ ] Event history

#### 5.5 Gacha Mechanics

  - [ ] Gacha interface
  - [ ] Opening animation
  - [ ] Obtained characters display
  - [ ] Gacha roll history

**Components to create:**

  - `GameWorld.vue` - Game world
  - `LocationCard.vue` - Location card
  - `CharacterCard.vue` - Character card
  - `EventDialog.vue` - Event dialog
  - `ChoiceButton.vue` - Choice button
  - `GachaInterface.vue` - Gacha interface
  - `GachaAnimation.vue` - Gacha animation

**Stores:**

  - `game.ts` - Game state
  - `characters.ts` - Player characters
  - `locations.ts` - Locations
  - `events.ts` - Events

**API endpoints:**

  - `GET /api/game/world` - World state
  - `GET /api/game/locations` - Locations
  - `GET /api/game/characters` - Characters
  - `GET /api/game/events` - Current events
  - `POST /api/game/choice` - Make a choice
  - `POST /api/game/gacha` - Gacha roll

-----

### Phase 6: Payment System 📋

**Status:** Planned

**Tasks:**

#### 6.1 Shop

  - [ ] Shop page `/shop`
  - [ ] Product catalog
  - [ ] Filtering and search
  - [ ] Shopping cart

#### 6.2 Premium Status

  - [ ] Premium subscription page
  - [ ] Plan comparison
  - [ ] Premium benefits
  - [ ] Subscription purchase

#### 6.3 In-game Currency

  - [ ] Currency bundles
  - [ ] Purchase bonuses
  - [ ] Special offers

#### 6.4 Payment Integration

  - [ ] Payment system integration
  - [ ] Payment processing
  - [ ] Payment confirmation
  - [ ] Purchase history

**Components:**

  - `ShopCatalog.vue` - Shop catalog
  - `ProductCard.vue` - Product card
  - `ShoppingCart.vue` - Shopping cart
  - `PremiumPlans.vue` - Premium plans
  - `PaymentForm.vue` - Payment form
  - `PurchaseHistory.vue` - Purchase history

**Stores:**

  - `shop.ts` - Shop state
  - `cart.ts` - Shopping cart
  - `premium.ts` - Premium status

**API endpoints:**

  - `GET /api/shop/products` - Products
  - `POST /api/shop/purchase` - Purchase
  - `GET /api/shop/history` - History
  - `POST /api/payment/create` - Create payment
  - `GET /api/payment/status/:id` - Payment status

-----

### Phase 7: Tutorial and Onboarding 📋

**Status:** Planned

**Tasks:**

#### 7.1 Tutorial

  - [ ] Interactive tutorial
  - [ ] Step-by-step instructions
  - [ ] Element highlighting
  - [ ] Tutorial progress

#### 7.2 Hints

  - [ ] Hint system
  - [ ] Contextual hints
  - [ ] Tooltips
  - [ ] Help information

#### 7.3 Achievements

  - [ ] Achievement system
  - [ ] Achievement progress
  - [ ] Achievement rewards

**Components:**

  - `TutorialOverlay.vue` - Tutorial overlay
  - `TutorialStep.vue` - Tutorial step
  - `Tooltip.vue` - Tooltip
  - `AchievementCard.vue` - Achievement card
  - `ProgressBar.vue` - Progress bar

**Stores:**

  - `tutorial.ts` - Tutorial state
  - `achievements.ts` - Achievements

-----

### Phase 8: Testing and Optimization 📋

**Status:** Planned

**Tasks:**

#### 8.1 Unit Tests

  - [ ] Tests for stores
  - [ ] Tests for composables
  - [ ] Tests for utils
  - [ ] Vitest setup

#### 8.2 E2E Tests

  - [ ] Auth tests
  - [ ] Gameplay tests
  - [ ] Purchase tests
  - [ ] Playwright setup

#### 8.3 Performance Optimization

  - [ ] Component lazy loading
  - [ ] Image optimization
  - [ ] Code splitting
  - [ ] API request caching
  - [ ] List virtualization

#### 8.4 SEO

  - [ ] Meta tags
  - [ ] Open Graph
  - [ ] Sitemap
  - [ ] robots.txt

**Tools:**

  - Vitest - Unit tests
  - Playwright - E2E tests
  - Lighthouse - Performance audit

-----

### Phase 9: Deploy and CI/CD 📋

**Status:** Planned

**Tasks:**

#### 9.1 Docker

  - [x] Dockerfile created
  - [ ] Docker Compose for dev
  - [ ] Image optimization
  - [ ] Multi-stage build

#### 9.2 CI/CD

  - [ ] GitHub Actions / GitLab CI
  - [ ] Automated tests
  - [ ] Automated deploy
  - [ ] Versioning

#### 9.3 Production

  - [ ] Server setup
  - [ ] SSL certificates
  - [ ] CDN for static assets
  - [ ] Monitoring and logging

-----

## Development Priorities

### High Priority (MVP)

1.  ✅ Basic Infrastructure
2.  🔄 Authentication
3.  📋 Basic Game Interface
4.  📋 Event and Choice System

### Medium Priority

5.  📋 User Profile
6.  📋 Gacha Mechanics
7.  📋 Tutorial

### Low Priority

8.  📋 Payment System
9.  📋 Achievements
10. 📋 Advanced Settings

-----

## Technical Requirements

### Performance

  - First Contentful Paint \< 1.5s
  - Time to Interactive \< 3.5s
  - Lighthouse Score \> 90

### Compatibility

  - Chrome 90+
  - Firefox 88+
  - Safari 14+
  - Edge 90+
  - Mobile browsers

### Security

  - HTTPS mandatory
  - XSS protection
  - CSRF protection
  - Secure token storage
  - Rate limiting

-----

## Backend Dependencies

### Required API Endpoints

#### Authentication

  - `POST /auth/register` - Registration
  - `POST /auth/login` - Login
  - `POST /auth/refresh` - Refresh token
  - `POST /auth/logout` - Logout

#### User

  - `GET /api/user/profile` - Profile
  - `PUT /api/user/profile` - Update profile
  - `GET /api/user/balance` - Balance
  - `GET /api/user/transactions` - Transactions

#### Game

  - `GET /api/game/world` - World state
  - `GET /api/game/locations` - Locations
  - `GET /api/game/characters` - Characters
  - `GET /api/game/events` - Events
  - `POST /api/game/choice` - Choice
  - `POST /api/game/gacha` - Gacha

#### Shop

  - `GET /api/shop/products` - Products
  - `POST /api/shop/purchase` - Purchase
  - `GET /api/shop/history` - History

-----

## Success Metrics

### Technical Metrics

  - ✅ 100% TypeScript coverage
  - ✅ ESLint 0 errors
  - 📋 80%+ test coverage
  - 📋 Lighthouse score \> 90

### User Metrics

  - 📋 Load time \< 3s
  - 📋 0 critical bugs
  - 📋 Responsive on all devices

-----

## Risks and Mitigation

### Risks

1.  **Backend API Delays** - May slow down integration

      - *Mitigation:* Use mocks for development

2.  **Mobile Performance** - Complex animations

      - *Mitigation:* Progressive enhancement, optimization

3.  **Game Logic Complexity** - May be hard to maintain

      - *Mitigation:* Good architecture, documentation

-----

## Next Steps

### Immediate (Phase 3)

1.  Complete authentication integration with API
2.  Create middleware for route protection
3.  Add error handling

### Short-term (1-2 weeks)

1.  Start user profile development
2.  Create basic game interface components
3.  Setup event system

### Long-term (1+ month)

1.  Full game interface implementation
2.  Payment system integration
3.  Testing and optimization