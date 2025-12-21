package domain

import (
	"context"

	"github.com/google/uuid"
)

type GameUserRepository interface {
	// --- Связи и ID ---
	
	// GetUserIDByGameID возвращает ID владельца (аккаунта) по UUID персонажа
	GetUserIDByGameID(ctx context.Context, gameUserID uuid.UUID) (int, error)
	
	// GetGameUsersByUserID возвращает список всех персонажей, привязанных к одному аккаунту
	GetGameUsersByUserID(ctx context.Context, userID int) ([]GameUser, error)

	// --- Управление персонажем (Lifecycle) ---
	
	// CountGameUsers возвращает текущее количество персонажей у пользователя (для проверки лимита в 5 штук)
	CountGameUsers(ctx context.Context, userID int) (int, error)
	
	// CreateEmptyGameUser создает "пустую" запись персонажа (только ID аккаунта, UUID и никнейм)
	// Настройки и мир на этом этапе будут иметь значения по умолчанию из БД (DEFAULT JSONB)
	CreateEmptyGameUser(ctx context.Context, userID int, nickname string) (uuid.UUID, error)
	
	// GetGameUser возвращает полные данные персонажа (включая JSONB поля) по его UUID
	GetGameUser(ctx context.Context, gameUserID uuid.UUID) (GameUser, error)
	
	// DeleteGameUser полностью удаляет персонажа и всю его историю (благодаря ON DELETE CASCADE)
	DeleteGameUser(ctx context.Context, gameUserID uuid.UUID) error

	// --- Работа с настройками (Settings) ---

	// UpdateSettings заменяет текущий JSONB объект настроек на новый
	UpdateSettings(ctx context.Context, gameUserID uuid.UUID, settings GameUserSettings) error
	
	// GetSettings возвращает только настройки персонажа (чтобы не тянуть всю строку из БД)
	GetSettings(ctx context.Context, gameUserID uuid.UUID) (GameUserSettings, error)

	// --- Работа с миром (World State) ---

	// UpdateWorldState обновляет данные о мире, иерархии и сеттинге персонажа
	UpdateWorldState(ctx context.Context, gameUserID uuid.UUID, world WorldState) error
	
	// GetWorldState возвращает только данные о мире персонажа
	GetWorldState(ctx context.Context, gameUserID uuid.UUID) (WorldState, error)
}

type GameUserService interface {
	// --- Связи ---

	// GetOwnerID находит владельца персонажа. Нужно для проверки прав доступа:
	GetOwnerID(ctx context.Context, gameUserID uuid.UUID) (int, error)

	// --- Создание персонажа ---

	// InitialCreateCharacter — Первый шаг. Проверяет, нет ли уже 5 персонажей. Если лимит превышен — возвращает ошибку. Если нет — создает "пустого" героя.
	InitialCreateCharacter(ctx context.Context, userID int, nickname string) (uuid.UUID, error)

	// --- Настройка по порядку ---

	// SetInitialSettings — Второй шаг. Применяет базовые настройки (тема, язык) после того, как "пустой" герой был успешно создан.
	SetInitialSettings(ctx context.Context, gameUserID uuid.UUID, settings GameUserSettings) error

	// DefineWorld — Третий шаг. Финализирует создание, устанавливая сеттинг (киберпанк и т.д.) и иерархию. После этого персонаж считается полностью готовым.
	DefineWorld(ctx context.Context, gameUserID uuid.UUID, world WorldState) error

	// --- Работа с существующим персонажем ---

	// GetFullProfile собирает все данные в одну структуру для отображения в игре.
	GetFullProfile(ctx context.Context, gameUserID uuid.UUID) (GameUser, error)

	// UpdateSettings полностью заменяет настройки пользователя (тема, язык и т.д.)
	UpdateSettings(ctx context.Context, gameUserID uuid.UUID, settings GameUserSettings) error

	// UpdateWorldState полностью обновляет состояние мира персонажа
	UpdateWorldState(ctx context.Context, gameUserID uuid.UUID, world WorldState) error

	// GetMyCharacters возвращает список всех персонажей, принадлежащих пользователю.
	GetMyCharacters(ctx context.Context, userID int) ([]GameUser, error)
}

type GameUserSettings struct {
    UI struct {
        Theme string `json:"theme"`
        Lang  string `json:"lang"`
    } `json:"ui"`
}

type WorldState struct {
    // Сеттинг: киберпанк, стимпанк, фэнтези
    Setting    string   `json:"setting"`    
    // Тон повествования: нуар, героика, хоррор
    Tone       string   `json:"tone"`       
    
    // Иерархия и статус
    Hierarchy struct {
        Role     string `json:"role"`     // гг: наемник, король, изгой
        Status   string `json:"status"`   // отношение общества: уважаемый/разыскиваемый
				Faction  string `json:"faction"`  // принадлежность к группе/организации
    } `json:"hierarchy"`

    // Дополнительные детали мира
    Features []string `json:"features"` // ["магия запрещена", "высокая преступность"]
}

// Хз пока что не добавил связи с User
type GameUser struct {
    ID       uuid.UUID        `db:"id"`
		UserID     int       `db:"user_id"`
    Nickname string           `db:"nickname"`
    Settings GameUserSettings `db:"settings"`
    WorldState WorldState `db:"world_state"`
}