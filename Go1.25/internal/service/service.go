package service

import (
	"github.com/ArtemChadaev/SeeThisGame/internal/domain"
	"github.com/ArtemChadaev/SeeThisGame/internal/repository"
	"github.com/redis/go-redis/v9"
)

// Service объединяет в себе все интерфейсы сервисов из ядра (domain)
type Service struct {
	domain.AuthorizationService
	domain.UserSettingsService
}

func NewService(repos *repository.Repository, redis *redis.Client) *Service {
	// Инициализируем конкретные реализации логики
	userSettingsService := NewUserSettingsService(repos.UserSettingsRepository, redis)
	authService := NewAuthService(repos.AuthorizationRepository, userSettingsService)

	return &Service{
		AuthorizationService: authService,
		UserSettingsService:  userSettingsService,
	}
}
