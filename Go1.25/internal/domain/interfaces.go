package domain

import (
	"time"
)

// --- REPOSITORY INTERFACES (Контракты для работы с данными) ---

type UserSettingsRepository interface {
	CreateUserSettings(settings UserSettings) error
	GetUserSettings(userId int) (UserSettings, error)
	UpdateUserSettings(settings UserSettings) error
	UpdateUserCoin(userId int, coin int) error
	BuyPaidSubscription(userId int, expiry time.Time) error
	DeactivateExpiredSubscriptions() (int64, error)
}

// --- SERVICE INTERFACES (Контракты бизнес-логики) ---

type UserSettingsService interface {
	CreateInitialUserSettings(userId int, name string) error
	GetByUserID(userId int) (UserSettings, error)
	UpdateInfo(userId int, name, icon string) error
	ChangeCoins(userId, amount int) error
	ActivateSubscription(userId, daysToAdd int, paymentToken string) error
	GetGrantDailyReward(userId int) error
}
