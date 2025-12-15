# Technical Specification for "Choose Me" Game Development

## 1. Program Name and Application Area

**Program Name:** "Choose Me".
**Application Area:** Software product intended for use in the field of interactive entertainment.
**Application Object:** Web platform used on personal computers and mobile devices of users through a web browser.

## 2. Development Basis

**Document on which development is based:** Project concept.
**Development Topic Name:** "Development of a life simulator game with RPG elements and procedural generation".

## 3. Development Purpose

**Functional Purpose:** Creating a game world populated by unique, procedurally generated characters controlled by AI. Providing the player with the role of observer and manager who influences the world and its inhabitants through game events and mechanics.
**Operational Purpose:** Providing 24/7 access for users to their game world, as well as providing a personal account for managing the account and in-game purchases.

## 4. Technical Requirements for the Program

### 4.1. Functional Characteristics Requirements

**Functions Performed for the Player:**

- Registration and authorization.
- Creating and customizing their world (choosing era, style).
- Obtaining characters through game mechanics (gacha, events).
- Sending characters to various locations for tasks, training, or living.
- Interacting with the world through an event system (quests) with choice options.
- Observing the life of characters, their interactions and development.
- Using special items to influence character generation.
- Purchasing in-game currency and premium status.

**Input Data Organization:** Data entered by users (registration data, event choices), world settings.
**Output Data Organization:** Web pages with visualization of the game world, characters, events; AI-generated images.

### 4.2. Reliability Requirements

- Ensuring stable website operation in 24/7 mode.
- **Decay Mechanics:** To maintain activity, a world decay system is implemented when the player is absent for a long time (1 day, week, month, year). A "Freeze" item is provided to prevent decay.
- Control and validation of user-entered information.
- Automatic database backup should be performed at least once a day.

### 4.3. Operating Conditions

- **Client Side:** Availability of a PC or mobile device with Internet access and a modern web browser installed (Google Chrome, Mozilla Firefox, Safari, Edge).
- **Personnel Qualifications:** Platform administration requires a specialist with system administration skills.

### 4.4. Technical Equipment Composition and Parameters Requirements

- **Server Side:** Web server (2-core processor, 4 GB RAM, 50 GB SSD).
- **Client Side:** User devices must ensure stable web browser operation.

### 4.5. Information and Software Compatibility Requirements

- **Server Software:** Linux OS, Nginx web server, PostgreSQL DBMS, Redis.
- **Development Stack:** Backend – Go, Frontend – Nuxt 4.
- **AI Integration:** ComfyUI, pollinations and others for image generation, n8n for automation and connection with AI models.
- **Information Protection:** Traffic encryption via HTTPS protocol.

## 5. Technical and Economic Indicators

- **Monetization:** Sale of in-game currency, premium status, and special items.
- **Economic Advantages:** Creating a unique gaming experience through deep AI integration and procedural generation, which increases player engagement and loyalty.

## 6. Development Stages and Phases

### Stage 1: Creating Web Platform and Basic Infrastructure

1.  **Website Development:** Landing page, registration/authorization, personal account, payment system integration.
2.  **Design and Interface:** Layout design, tutorial page creation.
3.  **Backend Infrastructure:** Server setup on Go and PostgreSQL database.

### Stage 2: Creating Prototype and Game Core

1.  **Basic Character Generation:** Based on JSON tags with portrait generation via AI.
2.  **World Creation:** Implementation of basic world structure and 2-3 locations.
3.  **Main Gameplay Loop:** "Summoning" characters and sending them to locations.

### Stage 3: Expanding Game Mechanics

1.  **NPC Interaction System:** Communication and relationship system.
2.  **Leveling and Skills:** Adding character development system.
3.  **Event System:** Creating quest editor with choices.

### Stage 4: Content, Polish and Launch

1.  **Content Expansion:** More locations, events, character types.
2.  **AI Enhancement:** Improving character behavior.
3.  **Polish:** Adding animations, interface debugging, testing and balancing.

## 7. Control and Acceptance Procedure

- **Types of Testing:** Unit, integration and acceptance testing.
- **Work Acceptance:** Performed after verifying that the implemented functionality meets the requirements set forth in this specification.
