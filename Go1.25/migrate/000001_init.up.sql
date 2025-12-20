CREATE TABLE users
(
    id            SERIAL       NOT NULL UNIQUE,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    CONSTRAINT users_pk PRIMARY KEY (id)
);
-- Токены для разных устройств
CREATE TABLE user_refresh_tokens
(
    id          SERIAL PRIMARY KEY,
    user_id     INT          NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token       VARCHAR(255) NOT NULL UNIQUE,
    expires_at  TIMESTAMPTZ  NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    name_device VARCHAR(255),
    device_info VARCHAR(255) -- Полезно для отладки
);
-- Функция для обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';
-- Тригер функции выше
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-- Настройки пользователя
CREATE TABLE user_settings
(
    user_id                   INT         NOT NULL UNIQUE,
    name                      VARCHAR(255) DEFAULT 'Alex',
    icon                      VARCHAR(255),
    coin                      INT DEFAULT 0,
    date_of_registration      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    paid_subscription         boolean DEFAULT FALSE,
--     Дата ОКОНЧАНИЯ её
    date_of_paid_subscription TIMESTAMPTZ,
    CONSTRAINT user_settings_pk PRIMARY KEY (user_id),
    CONSTRAINT fk_user_settings_user_id
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);