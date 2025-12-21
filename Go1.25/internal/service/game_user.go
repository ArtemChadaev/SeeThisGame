package service

import (
	"context"
	"errors"

	"github.com/ArtemChadaev/SeeThisGame/internal/domain"
	"github.com/google/uuid"
)

type GameUserService struct {
	repo domain.GameUserRepository
}

func NewGameUserService(repo domain.GameUserRepository) *GameUserService {
	return &GameUserService{
		repo: repo,
	}
}

// --- Связи ---

// GetOwnerID находит владельца персонажа.
func (s *GameUserService) GetOwnerID(ctx context.Context, gameUserID uuid.UUID) (int, error) {
	return s.repo.GetUserIDByGameID(ctx, gameUserID)
}

// --- Создание персонажа ---

// InitialCreateCharacter проверяет лимит (5 шт) и создает пустую запись персонажа.
func (s *GameUserService) InitialCreateCharacter(ctx context.Context, userID int, nickname string) (uuid.UUID, error) {
	// 1. Проверяем текущее количество персонажей
	count, err := s.repo.CountGameUsers(ctx, userID)
	if err != nil {
		return uuid.Nil, domain.NewInternalServerError(err)
	}

	// 2. Ограничение на 5 персонажей
	if count >= 5 {
		return uuid.Nil, errors.New("limit_reached: достигнут максимум персонажей (5)")
	}

	// 3. Создаем пустого персонажа
	gameUserID, err := s.repo.CreateEmptyGameUser(ctx, userID, nickname)
	if err != nil {
		return uuid.Nil, domain.NewInternalServerError(err)
	}

	return gameUserID, nil
}

// --- Настройка по порядку ---

// SetInitialSettings обновляет UI настройки (тема, язык) после создания.
func (s *GameUserService) SetInitialSettings(ctx context.Context, gameUserID uuid.UUID, settings domain.GameUserSettings) error {
	return s.repo.UpdateSettings(ctx, gameUserID, settings)
}

// DefineWorld устанавливает сеттинг и иерархию, завершая создание мира.
func (s *GameUserService) DefineWorld(ctx context.Context, gameUserID uuid.UUID, world domain.WorldState) error {
	return s.repo.UpdateWorldState(ctx, gameUserID, world)
}

// --- Работа с существующим персонажем ---

// GetFullProfile возвращает все данные GameUser (Nickname, Level, Settings, World).
func (s *GameUserService) GetFullProfile(ctx context.Context, gameUserID uuid.UUID) (domain.GameUser, error) {
	return s.repo.GetGameUser(ctx, gameUserID)
}

// UpdateSettings полностью заменяет настройки пользователя (тема, язык и т.д.)
// Теперь это не только тема, а весь объект GameUserSettings.
func (s *GameUserService) UpdateSettings(ctx context.Context, gameUserID uuid.UUID, settings domain.GameUserSettings) error {
	// Мы просто пробрасываем структуру в репозиторий.
	if err := s.repo.UpdateSettings(ctx, gameUserID, settings); err != nil {
		return domain.NewInternalServerError(err)
	}
	return nil
}

// UpdateWorldState полностью обновляет состояние мира персонажа
// Сюда входит: сеттинг, тон, вся иерархия (роль, статус, фракция) и список особенностей.
func (s *GameUserService) UpdateWorldState(ctx context.Context, gameUserID uuid.UUID, world domain.WorldState) error {
	// Репозиторий должен выполнить: UPDATE game_users SET world_state = $1 WHERE id = $2
	if err := s.repo.UpdateWorldState(ctx, gameUserID, world); err != nil {
		return domain.NewInternalServerError(err)
	}
	return nil
}

// GetMyCharacters возвращает список всех персонажей, принадлежащих пользователю.
func (s *GameUserService) GetMyCharacters(ctx context.Context, userID int) ([]domain.GameUser, error) {
	// Вызываем метод репозитория
	characters, err := s.repo.GetGameUsersByUserID(ctx, userID)
	if err != nil {
		// Оборачиваем системную ошибку в нашу внутреннюю ошибку для логирования
		return nil, domain.NewInternalServerError(err)
	}

	// Если у пользователя нет персонажей, вернется пустой слайс (не nil)
	return characters, nil
}