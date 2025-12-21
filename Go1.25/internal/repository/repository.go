package repository

import (
	"github.com/ArtemChadaev/SeeThisGame/internal/domain"
	"github.com/jmoiron/sqlx"
)

// Repository объединяет в себе все интерфейсы репозиториев из ядра (domain)
type Repository struct {
	domain.AuthorizationRepository
	domain.UserSettingsRepository
	domain.GameUserRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		// Здесь мы инициализируем конкретные реализации (например, из postgres)
		AuthorizationRepository: NewAuthPostgres(db),
		UserSettingsRepository:  NewUserSettingsPostgres(db),
		GameUserRepository:      NewGameUserPostgres(db),
	}
}
