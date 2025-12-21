package repository

import (
	"context"
	"encoding/json"
	"github.com/ArtemChadaev/SeeThisGame/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type GameUserRepository struct {
	db *sqlx.DB
}

func NewGameUserPostgres(db *sqlx.DB) *GameUserRepository {
	return &GameUserRepository{db: db}
}

// --- Связи и ID ---

// GetUserIDByGameID возвращает ID владельца (аккаунта) по UUID персонажа
func (r *GameUserRepository) GetUserIDByGameID(ctx context.Context, gameUserID uuid.UUID) (int, error) {
	var userID int
	query := "SELECT user_id FROM game_users WHERE id = $1"
	err := r.db.Get(&userID, query, gameUserID)
	return userID, err
}

// GetGameUsersByUserID возвращает список всех персонажей, привязанных к одному аккаунту
func (r *GameUserRepository) GetGameUsersByUserID(ctx context.Context, userID int) ([]domain.GameUser, error) {
	var gameUsers []domain.GameUser
	query := "SELECT * FROM game_users WHERE user_id = $1"
	err := r.db.Select(&gameUsers, query, userID)
	return gameUsers, err
}

// --- Управление персонажем (Lifecycle) ---

// CountGameUsers возвращает текущее количество персонажей у пользователя (для проверки лимита в 5 штук)
func (r *GameUserRepository) CountGameUsers(ctx context.Context, userID int) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM game_users WHERE user_id = $1"
	err := r.db.Get(&count, query, userID)
	return count, err
}

// CreateEmptyGameUser создает "пустую" запись персонажа (только ID аккаунта, UUID и никнейм)
// Настройки и мир на этом этапе будут иметь значения по умолчанию из БД (DEFAULT JSONB)
func (r *GameUserRepository) CreateEmptyGameUser(ctx context.Context, userID int, nickname string) (uuid.UUID, error) {
	newUUID := uuid.New()
	query := "INSERT INTO game_users (id, user_id, nickname) VALUES ($1, $2, $3)"
	_, err := r.db.Exec(query, newUUID, userID, nickname)
	return newUUID, err
}

// GetGameUser возвращает полные данные персонажа (включая JSONB поля) по его UUID
func (r *GameUserRepository) GetGameUser(ctx context.Context, gameUserID uuid.UUID) (domain.GameUser, error) {
	var gameUser domain.GameUser
	query := "SELECT * FROM game_users WHERE id = $1"
	err := r.db.Get(&gameUser, query, gameUserID)
	return gameUser, err
}

// DeleteGameUser полностью удаляет персонажа и всю его историю (ON DELETE CASCADE)
func (r *GameUserRepository) DeleteGameUser(ctx context.Context, gameUserID uuid.UUID) error {
	query := "DELETE FROM game_users WHERE id = $1"
	_, err := r.db.Exec(query, gameUserID)
	return err
}

// --- Работа с настройками (Settings) ---

// UpdateSettings заменяет текущий JSONB объект настроек на новый
func (r *GameUserRepository) UpdateSettings(ctx context.Context, gameUserID uuid.UUID, settings domain.GameUserSettings) error {
	// 1. Маршалим (превращаем) структуру в JSON
	settingsRaw, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	// 2. Отправляем в базу уже готовые байты JSON
	query := "UPDATE game_users SET settings = $1 WHERE id = $2"
	_, err = r.db.ExecContext(ctx, query, settingsRaw, gameUserID)
	return err
}

// GetSettings возвращает только настройки персонажа (чтобы не тянуть всю строку из БД)
func (r *GameUserRepository) GetSettings(ctx context.Context, gameUserID uuid.UUID) (domain.GameUserSettings, error) {
	var settings domain.GameUserSettings
	query := "SELECT settings FROM game_users WHERE id = $1"
	err := r.db.Get(&settings, query, gameUserID)
	return settings, err
}

// --- Работа с миром (World State) ---

// UpdateWorldState обновляет данные о мире, иерархии и сеттинге персонажа
func (r *GameUserRepository) UpdateWorldState(ctx context.Context, gameUserID uuid.UUID, world domain.WorldState) error {
	// 1. Маршалим структуру мира
	worldRaw, err := json.Marshal(world)
	if err != nil {
		return err
	}

	// 2. Отправляем в базу
	query := "UPDATE game_users SET world_state = $1 WHERE id = $2"
	_, err = r.db.ExecContext(ctx, query, worldRaw, gameUserID)
	return err
}

// GetWorldState возвращает только данные о мире персонажа
func (r *GameUserRepository) GetWorldState(ctx context.Context, gameUserID uuid.UUID) (domain.WorldState, error) {
	var world domain.WorldState
	query := "SELECT world_state FROM game_users WHERE id = $1"
	err := r.db.Get(&world, query, gameUserID)
	return world, err
}

