package service

import (
	"context"
	"errors"
	"time"

	"github.com/ArtemChadaev/SeeThisGame/internal/domain"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

const (
	// mockPaymentToken — простой токен для имитации успешной оплаты в pet-проекте
	mockPaymentToken = "mock-success-payment-token"

	// checkExpiredSubscriptionsInterval — интервал очистки просроченных подписок
	checkExpiredSubscriptionsInterval = 10 * time.Minute

	// dayCoins — количество монеток, выдаваемых ежедневно
	dayCoins = 3
)

type UserSettingsService struct {
	repo  domain.UserSettingsRepository // Используем интерфейс из domain
	redis *redis.Client
}

func NewUserSettingsService(repo domain.UserSettingsRepository, redis *redis.Client) *UserSettingsService {
	service := &UserSettingsService{
		repo:  repo,
		redis: redis,
	}

	// Запускаем фоновую задачу для проверки подписок
	go service.startSubscriptionChecker()

	return service
}

// CreateInitialUserSettings создает начальные настройки для нового пользователя.
func (s *UserSettingsService) CreateInitialUserSettings(userId int, name string) error {
	settings := domain.UserSettings{ // Используем конкретную структуру
		UserID:             userId,
		Name:               name,
		DateOfRegistration: time.Now(),
	}
	return s.repo.CreateUserSettings(settings)
}

// GetByUserID возвращает настройки пользователя по его ID.
func (s *UserSettingsService) GetByUserID(userId int) (domain.UserSettings, error) {
	return s.repo.GetUserSettings(userId)
}

// UpdateInfo обновляет имя и иконку пользователя.
func (s *UserSettingsService) UpdateInfo(userId int, name, icon string) error {
	settings, err := s.repo.GetUserSettings(userId)
	if err != nil {
		return err
	}

	settings.Name = name
	if icon != "" {
		settings.Icon = &icon
	}

	return s.repo.UpdateUserSettings(settings)
}

// ChangeCoins изменяет баланс монет пользователя (добавляет или списывает).
func (s *UserSettingsService) ChangeCoins(userId, coin int) error {
	settings, err := s.repo.GetUserSettings(userId)
	if err != nil {
		return err // Ошибка будет обработана выше (например, UserNotFound)
	}

	newBalance := settings.Coin + coin
	if newBalance < 0 {
		return errors.New("insufficient coins") // Или domain.ErrNoCoins
	}

	return s.repo.UpdateUserCoin(userId, newBalance)
}

// ActivateSubscription активирует или продлевает подписку.
func (s *UserSettingsService) ActivateSubscription(userId, daysToAdd int, paymentToken string) error {
	if paymentToken != mockPaymentToken {
		return errors.New("payment failed")
	}

	settings, err := s.repo.GetUserSettings(userId)
	if err != nil {
		return err
	}

	var newExpirationDate time.Time

	// Если подписка активна, продлеваем её от даты окончания, иначе — от текущего момента
	if settings.PaidSubscription && settings.DateOfPaidSubscription != nil && settings.DateOfPaidSubscription.After(time.Now()) {
		newExpirationDate = settings.DateOfPaidSubscription.AddDate(0, 0, daysToAdd)
	} else {
		newExpirationDate = time.Now().AddDate(0, 0, daysToAdd)
	}

	return s.repo.BuyPaidSubscription(userId, newExpirationDate)
}

// GetGrantDailyReward выдает ежедневную награду, используя Redis для контроля.
func (s *UserSettingsService) GetGrantDailyReward(userId int) error {
	// Ключ уникален для каждого дня
	key := "daily_rewards:" + time.Now().UTC().Format("2006-01-02")

	// Атомарно проверяем и добавляем пользователя в Redis Set
	added, err := s.redis.SAdd(context.Background(), key, userId).Result()
	if err != nil {
		return err
	}

	if added == 0 {
		return errors.New("reward already claimed today")
	}

	// Устанавливаем TTL для автоматической очистки ключа
	s.redis.Expire(context.Background(), key, 24*time.Hour)

	// Обновляем монеты в БД
	if err := s.ChangeCoins(userId, dayCoins); err != nil {
		s.redis.SRem(context.Background(), key, userId) // Откатываем Redis при ошибке БД
		return err
	}

	return nil
}

// startSubscriptionChecker — фоновый процесс для деактивации истекших подписок.
func (s *UserSettingsService) startSubscriptionChecker() {
	ticker := time.NewTicker(checkExpiredSubscriptionsInterval)
	defer ticker.Stop()

	logrus.Infof("Фоновая задача: проверка подписок каждые %v", checkExpiredSubscriptionsInterval)

	for range ticker.C {
		rowsAffected, err := s.repo.DeactivateExpiredSubscriptions()
		if err != nil {
			logrus.Errorf("Ошибка при деактивации подписок: %v", err)
			continue
		}

		if rowsAffected > 0 {
			logrus.Infof("Деактивировано %d просроченных подписок", rowsAffected)
		}
	}
}
