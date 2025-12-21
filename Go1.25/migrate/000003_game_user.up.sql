CREATE TABLE game_users (
    -- Генерируем UUID для каждого персонажа/профиля
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- Ссылка на основной аккаунт (Сайт)
    user_id INT NOT NULL,
    nickname VARCHAR(255) NOT NULL,

    -- Все настройки персонажа/профиля
    settings JSONB DEFAULT '{}',

    -- Всё о мире игры
    world_state JSONB NOT NULL DEFAULT '{}',
    
    CONSTRAINT fk_account
        FOREIGN KEY (user_id) 
        REFERENCES users(id)
        ON DELETE CASCADE
);
