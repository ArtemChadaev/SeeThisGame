package service

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ArtemChadaev/SeeThisGame/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

const (
	salt                  = "asdagedrhftyki518sadf5as8"
	signingKey            = "awsg8s#@4Sf86DS#$2dF"
	accessTokenTTL        = time.Minute * 15
	refreshTokenTTL       = time.Hour * 24 * 365
	updateRefreshTokenTTL = time.Hour * 24 * 90
)


type AuthService struct {
	repo            domain.AuthorizationRepository
	settingsService domain.UserSettingsService
}

func NewAuthService(repo domain.AuthorizationRepository, settingsService domain.UserSettingsService) *AuthService {
	return &AuthService{
		repo:            repo,
		settingsService: settingsService,
	}
}

// --- Помощники (Helpers) ---

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func generateRefreshToken() (string, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(tokenBytes), nil
}

func (s *AuthService) newRefreshToken(userId int) (domain.RefreshToken, error) {
	token, err := generateRefreshToken()
	if err != nil {
		return domain.RefreshToken{}, err
	}

	return domain.RefreshToken{
		UserID:    userId,
		Token:     token,
		ExpiresAt: time.Now().Add(refreshTokenTTL),
	}, nil
}

// --- Основные методы ---

// CreateUser — создание пользователя с сайта (без персонажа)
func (s *AuthService) CreateUser(user domain.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)

	id, err := s.repo.CreateUser(user)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return 0, domain.ErrUserAlreadyExists
		}
		return 0, domain.NewInternalServerError(err)
	}

	userName := strings.Split(user.Email, "@")[0]
	if err := s.settingsService.CreateInitialUserSettings(id, userName); err != nil {
		return 0, domain.NewInternalServerError(err)
	}

	return id, nil
}

// GenerateTokens — логин с сайта. Выдает токены без GameUserID.
func (s *AuthService) GenerateTokens(email, password string) (domain.ResponseTokens, error) {
	userId, err := s.repo.GetUser(email, generatePasswordHash(password))
	if err != nil {
		return domain.ResponseTokens{}, domain.ErrInvalidCredentials
	}

	return s.createTokens(userId, uuid.Nil) // UUID персонажа пустой
}

// GenerateGameTokens — вызывается при выборе персонажа. "Апгрейдит" токены.
func (s *AuthService) GenerateGameTokens(userId int, gameUserId uuid.UUID) (domain.ResponseTokens, error) {
    //TODO: можно добавить проверку, принадлежит ли gameUserId этому userId через repo
	return s.createTokens(userId, gameUserId)
}

// createTokens — универсальный приватный метод для сборки пары токенов.
func (s *AuthService) createTokens(userId int, gameUserId uuid.UUID) (domain.ResponseTokens, error) {
	// Используем domain.MyClaims
	claims := &domain.MyClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID:     userId,
		GameUserID: gameUserId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return domain.ResponseTokens{}, domain.NewInternalServerError(err)
	}

	refresh, err := s.newRefreshToken(userId)
	if err != nil {
		return domain.ResponseTokens{}, err
	}

	if err = s.repo.CreateToken(refresh); err != nil {
		return domain.ResponseTokens{}, domain.NewInternalServerError(err)
	}

	return domain.ResponseTokens{
		AccessToken:  accessToken,
		RefreshToken: refresh.Token,
	}, nil
}

// GetAccessToken — обновление токенов. 
// Принимает опциональный gameUserId, чтобы сохранить контекст игры при обновлении.
func (s *AuthService) GetAccessToken(refreshToken string, gameUserId uuid.UUID) (domain.ResponseTokens, error) {
	refresh, err := s.repo.GetRefreshToken(refreshToken)
	if err != nil {
		return domain.ResponseTokens{}, domain.ErrInvalidToken
	}

	if time.Now().After(refresh.ExpiresAt) {
		_ = s.repo.DeleteRefreshToken(refresh.ID)
		return domain.ResponseTokens{}, domain.ErrInvalidToken
	}

	// Создаем новый Access Token через общий метод
	// Теперь он будет знать про gameUserId, если тот был передан
	tokens, err := s.createTokens(refresh.UserID, gameUserId)
	if err != nil {
		return domain.ResponseTokens{}, err
	}

	// Если Refresh Token скоро истечет, обновляем и его
	if refresh.ExpiresAt.Before(time.Now().Add(updateRefreshTokenTTL)) {
		newRefresh, err := s.newRefreshToken(refresh.UserID)
		if err != nil {
			return domain.ResponseTokens{}, err
		}

		if err := s.repo.UpdateToken(refreshToken, newRefresh); err != nil {
			return domain.ResponseTokens{}, err
		}
		tokens.RefreshToken = newRefresh.Token
	} else {
		// Если рефреш не обновляли, возвращаем старый
		tokens.RefreshToken = refresh.Token
	}

	return tokens, nil
}

// ParseToken — теперь возвращает *domain.MyClaims
func (s *AuthService) ParseToken(accessToken string) (*domain.MyClaims, error) {
	claims := &domain.MyClaims{}

	token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(signingKey), nil
	})

	if err != nil || !token.Valid {
		return nil, domain.ErrInvalidToken
	}

	return claims, nil
}

func (s *AuthService) UnAuthorize(refreshToken string) error {
	refresh, err := s.repo.GetRefreshToken(refreshToken)
	if err != nil {
		return err
	}
	return s.repo.DeleteRefreshToken(refresh.ID)
}

func (s *AuthService) UnAuthorizeAll(email, password string) error {
	id, err := s.repo.GetUser(email, generatePasswordHash(password))
	if err != nil {
		return err
	}
	return s.repo.DeleteAllUserRefreshTokens(id)
}
