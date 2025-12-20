# Описание

## Данные

- **Version**: 15
- **Type image**: alpine
- **Port**: 5432

Все таблицы в [migrate](../migrate)

# PostgreSQL схемы

## Пользователь

**users**
```sql
- id (SERIAL PRIMARY KEY)
- email (VARCHAR UNIQUE)
- password_hash (VARCHAR)
```

**user_refresh_tokens**
```sql
- id (SERIAL PRIMARY KEY)
- user_id (INT FK → users)
- token (VARCHAR UNIQUE)
- expires_at (TIMESTAMPTZ)
- name_device (VARCHAR)
- device_info (VARCHAR)
```

**user_settings**
```sql
- user_id (INT PRIMARY KEY FK → users)
- name (VARCHAR)
- icon (VARCHAR)
- coin (INT)
- date_of_registration (TIMESTAMPTZ)
- paid_subscription (BOOLEAN)
- date_of_paid_subscription (TIMESTAMPTZ)
```

## Дополнительные активности пользователя


## Игра (Часто дополнятся будет)

### Пользователь игры

### Персонажи и история

**character** - Проще говоря какой то персонаж, есть определенные поля и много JSONB
```postgresql
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    
    -- Основные параметры
    name VARCHAR(32) NOT NULL,
    level INTEGER DEFAULT 1 NOT NULL CHECK (level > 0),
    rarity SMALLINT DEFAULT 0 NOT NULL CHECK (rarity >= 0), -- Твоя система редкости
    experience BIGINT DEFAULT 0 NOT NULL,
    
    -- Системные поля
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    -- Гибкие данные
    stats JSONB NOT NULL DEFAULT '{}',       -- Сила, мана и т.д.
    appearance JSONB NOT NULL DEFAULT '{}',  -- Визуал
    
    -- "Память" и "Отношения"
    memory JSONB NOT NULL DEFAULT '[]',       -- Хронология событий персонажа
    relationships JSONB NOT NULL DEFAULT '{}' -- Связи с другими ID
```

### Вещи в игре

**item_templates** - Основа предметов 
```postgresql
    id SERIAL PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    base_weight NUMERIC(10, 3) NOT NULL, -- Вес по умолчанию
    category SMALLINT NOT NULL,          -- 1: Оружие, 2: Еда и т.д.
    is_customizable BOOLEAN DEFAULT FALSE, -- Можно ли менять этот предмет
    static_data JSONB DEFAULT '{}'       -- Описание, иконка, базовая цена
```

**items** - Сами предметы
```postgresql
    id SERIAL PRIMARY KEY,
    owner_id UUID NOT NULL, 
    template_id INTEGER REFERENCES item_templates(id),
    
    -- Динамический вес: если предмет изменен игроком, пишем сюда.
    -- Если NULL — берем base_weight из шаблона.
    current_weight NUMERIC(10, 3), 
    
    -- Кастомное имя (если игрок переименовал предмет)
    custom_name VARCHAR(64),
    
    -- Специфические данные конкретного экземпляра
    properties JSONB NOT NULL DEFAULT '{}'
```

### Клан (реализован в ооочень далеком будущем)

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

## Entities for Creating/Using AI in n8n

**prompt_histories** - Description of history/event sent to n8n (JSON sent in [Response: n8n → Backend Callback](../../n8n/README.md))
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

**history** - Returned history, JSON: [Response: n8n → Backend Callback](../../n8n/README.md)
```sql
- id (INT UNIQUE FK → prompt_histories)
- data (JSONB NOT NULL)
  - []
    - number (SMALLINT)
    - text (TEXT)
    - image (TEXT)
```
