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
	"github.com/lib/pq"
)

const (
	salt                  = "asdagedrhftyki518sadf5as8"
	signingKey            = "awsg8s#@4Sf86DS#$2dF"
	accessTokenTTL        = time.Minute * 15
	refreshTokenTTL       = time.Hour * 24 * 365
	updateRefreshTokenTTL = time.Hour * 24 * 90
)

type tokenClaims struct {
	jwt.RegisteredClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	repo            domain.AuthorizationRepository // Используем интерфейс из domain
	settingsService domain.UserSettingsService     // Ссылка на сервис настроек через интерфейс
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

func (s *AuthService) CreateUser(user domain.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)

	id, err := s.repo.CreateUser(user)
	if err != nil {
		// Проверяем ошибку на нарушение уникальности (Unique Violation) в Postgres
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return 0, domain.ErrUserAlreadyExists
		}
		// Все остальные системные ошибки оборачиваем в InternalServerError
		return 0, domain.NewInternalServerError(err)
	}

	userName := strings.Split(user.Email, "@")[0]
	if err := s.settingsService.CreateInitialUserSettings(id, userName); err != nil {
		return 0, domain.NewInternalServerError(err)
	}

	return id, nil
}

func (s *AuthService) GenerateTokens(email, password string) (domain.ResponseTokens, error) {
	userId, err := s.repo.GetUser(email, generatePasswordHash(password))
	if err != nil {
		// Если пользователь не найден в БД, возвращаем типизированную ошибку
		return domain.ResponseTokens{}, domain.ErrInvalidCredentials
	}

	return s.createTokens(userId)
}

func (s *AuthService) GenerateTokensForUser(userId int) (domain.ResponseTokens, error) {
	return s.createTokens(userId)
}

func (s *AuthService) createTokens(userId int) (domain.ResponseTokens, error) {
	// 1. Создаем Access Token (JWT)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserId: userId,
	})

	accessToken, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return domain.ResponseTokens{}, domain.NewInternalServerError(err)
	}

	// 2. Создаем Refresh Token
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

func (s *AuthService) GetAccessToken(refreshToken string) (domain.ResponseTokens, error) {
	refresh, err := s.repo.GetRefreshToken(refreshToken)
	if err != nil {
		return domain.ResponseTokens{}, domain.ErrInvalidToken
	}

	if time.Now().After(refresh.ExpiresAt) {
		_ = s.repo.DeleteRefreshToken(refresh.ID)
		return domain.ResponseTokens{}, domain.ErrInvalidToken
	}

	// Создаем новый Access Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		refresh.UserID,
	})

	accessToken, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return domain.ResponseTokens{}, domain.NewInternalServerError(err)
	}

	// Если Refresh Token скоро истечет, обновляем и его (Rotating Refresh Tokens)
	currentRefreshToken := refresh.Token
	if refresh.ExpiresAt.Before(time.Now().Add(updateRefreshTokenTTL)) {
		newRefresh, err := s.newRefreshToken(refresh.UserID)
		if err != nil {
			return domain.ResponseTokens{}, err
		}

		if err := s.repo.UpdateToken(refreshToken, newRefresh); err != nil {
			return domain.ResponseTokens{}, err
		}
		currentRefreshToken = newRefresh.Token
	}

	return domain.ResponseTokens{
		AccessToken:  accessToken,
		RefreshToken: currentRefreshToken,
	}, nil
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(signingKey), nil
	})

	if err != nil {
		return 0, domain.ErrInvalidToken
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok || !token.Valid {
		return 0, domain.ErrInvalidToken
	}

	return claims.UserId, nil
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
