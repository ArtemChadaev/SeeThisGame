# Description

## Location

- **Version**: 15
- **Type**: alpine
- **Port**: 5432

New tables are populated in [migrate](./migrate)

# PostgreSQL Schemas

## Main Tables

**users** - System users
```sql
- id (SERIAL PRIMARY KEY)
- email (VARCHAR UNIQUE)
- password_hash (VARCHAR)
```

**user_refresh_tokens** - Refresh tokens for different devices
```sql
- id (SERIAL PRIMARY KEY)
- user_id (INT FK → users)
- token (VARCHAR UNIQUE)
- expires_at (TIMESTAMPTZ)
- name_device (VARCHAR)
- device_info (VARCHAR)
```

**user_settings** - User settings and profile
```sql
- user_id (INT PRIMARY KEY FK → users)
- name (VARCHAR)
- icon (VARCHAR)
- coin (INT)
- date_of_registration (TIMESTAMPTZ)
- paid_subscription (BOOLEAN)
- date_of_paid_subscription (TIMESTAMPTZ)
```

## Clan System

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

**clan_members** - Clan members
```sql
- clan_id (INT FK → clan)
- user_id (INT FK → users)
- role_id (SMALLINT FK → roles)
- PRIMARY KEY (clan_id, user_id)
```

**clan_role_names** - Custom role names for each clan
```sql
- clan_id (INT FK → clan)
- role_id (SMALLINT FK → roles)
- custom_name (VARCHAR)
- PRIMARY KEY (clan_id, role_id)
```

## Game Entities

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

## Entities for Creating/Using AI in n8n

**prompt_histories** - Description of history/event sent to n8n (JSON sent in [Response: n8n → Backend Callback](../n8n/README.md))
```sql
- id (SERIAL PRIMARY KEY)
- prompt (TEXT NOT NULL)
- folder_id (INT NOT NULL)
- number (SMALLINT NOT NULL)
- language (VARCHAR)
- model_text (VARCHAR)
- model_image (VARCHAR)
- width (SMALLINT)
- height (SMALLINT)
- place (VARCHAR)
- time (VARCHAR)
- season (VARCHAR)
- style (VARCHAR)
- loras (VARCHAR)
- safe (BOOLEAN)
- generated (VARCHAR)
```

**history** - Returned history, JSON: [Response: n8n → Backend Callback](../n8n/README.md)
```sql
- id (INT UNIQUE FK → prompt_histories)
- data (JSONB NOT NULL)
  - []
    - number (SMALLINT)
    - text (TEXT)
    - image (TEXT)
```
