CREATE TABLE prompt_histories (
    -- Идентификаторы
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,          -- Владелец истории
    folder_id INT NOT NULL,         -- Просто ID группы (сессии), без внешней таблицы
    number SMALLINT NOT NULL,       -- Порядок сообщения в этой группе
    
    -- Контент
    prompt TEXT NOT NULL,
    generated VARCHAR(255) NOT NULL,  -- состояние
    
    -- Окружение (может меняться внутри одной "папки")
    place VARCHAR(100),
    time VARCHAR(50),
    season VARCHAR(50),
    style VARCHAR(100),
    
    -- Технические параметры генерации
    model_text VARCHAR(100),
    model_image VARCHAR(100),
    width SMALLINT,
    height SMALLINT,
    loras JSONB DEFAULT '[]',       -- Список использованных Lora
    safe BOOLEAN DEFAULT TRUE,
    
    -- Системное
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE history (
    prompt_history_id INT PRIMARY KEY REFERENCES prompt_histories(id) ON DELETE CASCADE,
    data JSONB NOT NULL
);

CREATE INDEX idx_ph_user_folder ON prompt_histories(user_id, folder_id);