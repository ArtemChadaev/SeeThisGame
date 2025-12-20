package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ArtemChadaev/SeeThisGame/internal/domain"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

// oauthUserInfo — локальная структура для временного хранения данных из провайдеров
type oauthUserInfo struct {
	ID      string
	Email   string
	Name    string
	Picture string
}

type OAuthService struct {
	repo         domain.AuthorizationRepository // Используем новый интерфейс из domain
	authService  *AuthService
	googleConfig *oauth2.Config
	githubConfig *oauth2.Config
}

func NewOAuthService(repo domain.AuthorizationRepository, authService *AuthService) *OAuthService {
	return &OAuthService{
		repo:        repo,
		authService: authService,
		googleConfig: &oauth2.Config{
			ClientID:     viper.GetString("oauth.google.clientID"),
			ClientSecret: viper.GetString("oauth.google.clientSecret"),
			RedirectURL:  viper.GetString("oauth.google.redirectURL"),
			Scopes:       viper.GetStringSlice("oauth.google.scopes"),
			Endpoint:     google.Endpoint,
		},
		githubConfig: &oauth2.Config{
			ClientID:     viper.GetString("oauth.github.clientID"),
			ClientSecret: viper.GetString("oauth.github.clientSecret"),
			RedirectURL:  viper.GetString("oauth.github.redirectURL"),
			Scopes:       viper.GetStringSlice("oauth.github.scopes"),
			Endpoint:     github.Endpoint,
		},
	}
}

func (s *OAuthService) GetAuthURL(provider string) (string, error) {
	var config *oauth2.Config
	switch provider {
	case "google":
		config = s.googleConfig
	case "github":
		config = s.githubConfig
	default:
		return "", errors.New("unsupported provider")
	}

	state := "random-state-string" // В идеале генерировать динамически
	return config.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

func (s *OAuthService) HandleCallback(provider, code string) (domain.ResponseTokens, error) {
	var config *oauth2.Config
	switch provider {
	case "google":
		config = s.googleConfig
	case "github":
		config = s.githubConfig
	default:
		return domain.ResponseTokens{}, errors.New("unsupported provider")
	}

	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		return domain.ResponseTokens{}, err
	}

	userInfo, err := s.getUserInfo(provider, token)
	if err != nil {
		return domain.ResponseTokens{}, err
	}

	return s.authenticateOAuthUser(provider, userInfo)
}

func (s *OAuthService) getUserInfo(provider string, token *oauth2.Token) (oauthUserInfo, error) {
	var userInfoURL string
	switch provider {
	case "google":
		userInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"
	case "github":
		userInfoURL = "https://api.github.com/user"
	default:
		return oauthUserInfo{}, errors.New("unsupported provider")
	}

	client := http.Client{}
	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return oauthUserInfo{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		return oauthUserInfo{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return oauthUserInfo{}, err
	}

	var res oauthUserInfo

	if provider == "google" {
		var googleUser struct {
			ID      string `json:"id"`
			Email   string `json:"email"`
			Name    string `json:"name"`
			Picture string `json:"picture"`
		}
		if err := json.Unmarshal(body, &googleUser); err != nil {
			return oauthUserInfo{}, err
		}
		res = oauthUserInfo{
			ID:      googleUser.ID,
			Email:   googleUser.Email,
			Name:    googleUser.Name,
			Picture: googleUser.Picture,
		}
	} else if provider == "github" {
		var githubUser struct {
			ID        int    `json:"id"`
			Login     string `json:"login"`
			Name      string `json:"name"`
			AvatarURL string `json:"avatar_url"`
			Email     string `json:"email"`
		}
		if err := json.Unmarshal(body, &githubUser); err != nil {
			return oauthUserInfo{}, err
		}
		res = oauthUserInfo{
			ID:      fmt.Sprintf("%d", githubUser.ID),
			Name:    githubUser.Name,
			Picture: githubUser.AvatarURL,
			Email:   githubUser.Email,
		}
		if res.Name == "" {
			res.Name = githubUser.Login
		}
		if res.Email == "" {
			email, _ := s.getGitHubEmail(token.AccessToken)
			res.Email = email
		}
	}

	return res, nil
}

func (s *OAuthService) getGitHubEmail(accessToken string) (string, error) {
	req, _ := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", err
	}

	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email, nil
		}
	}
	return "", errors.New("no email found")
}

func (s *OAuthService) authenticateOAuthUser(provider string, userInfo oauthUserInfo) (domain.ResponseTokens, error) {
	// 1. Пытаемся найти по OAuth ID
	user, err := s.repo.GetUserByOAuth(provider, userInfo.ID)
	if err == nil {
		return s.authService.GenerateTokensForUser(user.ID) // Метод должен быть в AuthService
	}

	// 2. Пытаемся найти по Email (привязка аккаунта)
	if userInfo.Email != "" {
		user, err = s.repo.GetUserByEmail(userInfo.Email)
		if err == nil {
			// В реальности здесь может быть логика обновления OAuthID для существующего юзера
			return s.authService.GenerateTokensForUser(user.ID)
		}
	}

	// 3. Создаем нового пользователя
	newUser := domain.User{
		Email:         userInfo.Email,
		OAuthProvider: &provider,
		OAuthID:       &userInfo.ID,
	}

	id, err := s.repo.CreateOAuthUser(newUser)
	if err != nil {
		return domain.ResponseTokens{}, err
	}

	// Создаем начальные настройки профиля
	// Мы передаем имя и иконку, полученные от провайдера
	if err := s.authService.settingsService.CreateInitialUserSettings(id, userInfo.Name); err != nil {
		// Логируем, но не прерываем вход
	}

	return s.authService.GenerateTokensForUser(id)
}
