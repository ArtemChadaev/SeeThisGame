package domain

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthorizationRepository interface {
	// User Management
	CreateUser(user User) (int, error)
	GetUser(username, password string) (int, error)
	GetUserEmailFromId(id int) (string, error)
	UpdateUserPassword(user User) error

	// Token Management
	GetUserIdByRefreshToken(token string) (int, error)
	CreateToken(token RefreshToken) error
	GetRefreshToken(token string) (RefreshToken, error)
	UpdateToken(oldToken string, newToken RefreshToken) error
	DeleteRefreshToken(tokenId int) error
	DeleteAllUserRefreshTokens(userId int) error
	GetRefreshTokens(userId int) ([]RefreshToken, error)
}

type AuthorizationService interface {
	CreateUser(user User) (int, error)
	GenerateTokens(email, password string) (ResponseTokens, error)
	GetAccessToken(refreshToken string, gameUserId uuid.UUID) (ResponseTokens, error)
	ParseToken(accessToken string) (*MyClaims, error)
	UnAuthorize(refreshToken string) error
	UnAuthorizeAll(email, password string) error
	GenerateGameTokens(userId int, gameUserId uuid.UUID) (ResponseTokens, error)
}

// MyClaims описывает содержимое JWT токена для аутентификации в игре.
type MyClaims struct {
    jwt.RegisteredClaims
    UserID     int       `json:"user_id"`
    GameUserID uuid.UUID `json:"game_user_id,omitempty"` // omitempty, если на сайте он не нужен
}

type ResponseTokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type User struct {
	ID            int     `json:"-" db:"id"`
	Email         string  `json:"email" binding:"required"`
	Password      string  `json:"password" binding:"required"`
	OAuthProvider *string `json:"oauth_provider,omitempty" db:"oauth_provider"`
	OAuthID       *string `json:"oauth_id,omitempty" db:"oauth_id"`
}

type RefreshToken struct {
	ID         int       `db:"id"`
	UserID     int       `db:"user_id"`
	Token      string    `db:"token"`
	ExpiresAt  time.Time `db:"expires_at"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	NameDevice *string   `db:"name_device"`
	DeviceInfo *string   `db:"device_info"`
}

type UserSettings struct {
	UserID                 int        `json:"id" db:"user_id"`
	Name                   string     `json:"name" db:"name"`
	Icon                   *string    `json:"icon" db:"icon"`
	Coin                   int        `json:"coin" db:"coin"`
	DateOfRegistration     time.Time  `json:"dateOfRegistration" db:"date_of_registration"`
	PaidSubscription       bool       `json:"paidSubscription" db:"paid_subscription"`
	DateOfPaidSubscription *time.Time `json:"dateOfPaidSubscription" db:"date_of_paid_subscription"`
}
