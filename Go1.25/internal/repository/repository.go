package repository

import (
	"time"

	"github.com/ArtemChadaev/SeeThisGame/internal/domain"
	"github.com/jmoiron/sqlx"
)

type Autorization interface {
	// User Management
	CreateUser(user domain.rest) (int, error)
	GetUser(username, password string) (int, error)
	GetUserEmailFromId(id int) (string, error)
	UpdateUserPassword(user domain.rest) error

	// Token Management
	GetUserIdByRefreshToken(refreshToken string) (int, error)
	CreateToken(refreshToken domain.rest) error
	GetRefreshToken(refreshToken string) (domain.rest, error)
	UpdateToken(oldRefreshToken string, refreshToken domain.rest) error
	DeleteRefreshToken(tokenId int) error
	DeleteAllUserRefreshTokens(userId int) error
	GetRefreshTokens(userId int) ([]domain.rest, error)

	// OAuth Management
	CreateOAuthUser(user domain.rest) (int, error)
	GetUserByOAuth(provider, oauthID string) (domain.rest, error)
	GetUserByEmail(email string) (domain.rest, error)
}

type UserSettings interface {
	// Profile Management
	CreateUserSettings(settings domain.rest) error
	GetUserSettings(userId int) (domain.rest, error)
	UpdateUserSettings(settings domain.rest) error

	// Economy
	UpdateUserCoin(userId int, coin int) error

	// Subscription
	BuyPaidSubscription(userId int, time time.Time) error
	DeactivateExpiredSubscriptions() (int64, error)
}

type Repository struct {
	Autorization
	UserSettings
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Autorization: NewAuthPostgres(db),
		UserSettings: NewUserSettingsPostgres(db),
	}
}
