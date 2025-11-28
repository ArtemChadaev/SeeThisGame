package repository

import (
	"time"

	"github.com/ArtemChadaev/go"
	"github.com/jmoiron/sqlx"
)

type Autorization interface {
	// User Management
	CreateUser(user rest.User) (int, error)
	GetUser(username, password string) (int, error)
	GetUserEmailFromId(id int) (string, error)
	UpdateUserPassword(user rest.User) error

	// Token Management
	GetUserIdByRefreshToken(refreshToken string) (int, error)
	CreateToken(refreshToken rest.RefreshToken) error
	GetRefreshToken(refreshToken string) (rest.RefreshToken, error)
	UpdateToken(oldRefreshToken string, refreshToken rest.RefreshToken) error
	DeleteRefreshToken(tokenId int) error
	DeleteAllUserRefreshTokens(userId int) error
	GetRefreshTokens(userId int) ([]rest.RefreshToken, error)

	// OAuth Management
	CreateOAuthUser(user rest.User) (int, error)
	GetUserByOAuth(provider, oauthID string) (rest.User, error)
	GetUserByEmail(email string) (rest.User, error)
}

type UserSettings interface {
	// Profile Management
	CreateUserSettings(settings rest.UserSettings) error
	GetUserSettings(userId int) (rest.UserSettings, error)
	UpdateUserSettings(settings rest.UserSettings) error

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
