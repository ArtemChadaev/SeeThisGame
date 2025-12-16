package service

import (
	"github.com/ArtemChadaev/SeeThisGame/internal/domain"
	"github.com/ArtemChadaev/SeeThisGame/internal/repository"
	"github.com/redis/go-redis/v9"
)

type Autorization interface {
	// Authentication
	CreateUser(user domain.rest) (int, error)
	GenerateTokens(email, password string) (tokens domain.rest, err error)
	GetAccessToken(refreshToken string) (tokens domain.rest, err error)
	ParseToken(accessToken string) (int, error)
	UnAuthorize(refreshToken string) error
	UnAuthorizeAll(email, password string) error
}

type UserSettings interface {
	// Profile Management
	CreateInitialUserSettings(userId int, name string) error
	GetByUserID(userId int) (domain.rest, error)
	UpdateInfo(userId int, name, icon string) error

	// Economy
	ChangeCoins(userId, amount int) error

	// Subscription
	ActivateSubscription(userId, daysToAdd int, paymentToken string) error

	// Rewards
	GetGrantDailyReward(userId int) error
}
type OAuth interface {
	GetAuthURL(provider string) (string, error)
	HandleCallback(provider, code string) (domain.rest, error)
}

type Service struct {
	Autorization
	UserSettings
	OAuth
}

func NewService(repos *repository.Repository, redis *redis.Client) *Service {
	userSettingsService := NewUserSettingsService(repos.UserSettings, redis)

	authService := NewAuthService(repos.Autorization, userSettingsService)
	oauthService := NewOAuthService(repos.Autorization, authService)

	return &Service{
		Autorization: authService,
		UserSettings: userSettingsService,
		OAuth:        oauthService,
	}
}
