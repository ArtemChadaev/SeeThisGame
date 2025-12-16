package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ArtemChadaev/SeeThisGame/internal/domain"
	"github.com/ArtemChadaev/SeeThisGame/internal/repository"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

type OAuthService struct {
	repo         repository.Autorization
	authService  *AuthService
	googleConfig *oauth2.Config
	githubConfig *oauth2.Config
}

func NewOAuthService(repo repository.Autorization, authService *AuthService) *OAuthService {
	return &OAuthService{
		repo:        repo,
		authService: authService,
		googleConfig: &oauth2.Config{
			ClientID:     viper.GetString("oauth.google.clientID"), // Will be loaded from env/secrets
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

	// Generate state token (in production should be random and stored in session/cookie)
	// For simplicity we use a static state, but this should be improved
	state := "random-state-string"
	return config.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

func (s *OAuthService) HandleCallback(provider, code string) (domain.rest, error) {
	var config *oauth2.Config
	switch provider {
	case "google":
		config = s.googleConfig
	case "github":
		config = s.githubConfig
	default:
		return domain.rest{}, errors.New("unsupported provider")
	}

	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		return domain.rest{}, err
	}

	userInfo, err := s.getUserInfo(provider, token)
	if err != nil {
		return domain.rest{}, err
	}

	return s.authenticateOAuthUser(provider, userInfo)
}

func (s *OAuthService) getUserInfo(provider string, token *oauth2.Token) (domain.rest, error) {
	var userInfoURL string
	switch provider {
	case "google":
		userInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"
	case "github":
		userInfoURL = "https://api.github.com/user"
	}

	client := http.Client{}
	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return domain.rest{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		return domain.rest{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.rest{}, err
	}

	var oauthUser domain.rest
	oauthUser.Provider = domain.rest.OAuthProvider(provider)

	if provider == "google" {
		var googleUser struct {
			ID      string `json:"id"`
			Email   string `json:"email"`
			Name    string `json:"name"`
			Picture string `json:"picture"`
		}
		if err := json.Unmarshal(body, &googleUser); err != nil {
			return domain.rest{}, err
		}
		oauthUser.ID = googleUser.ID
		oauthUser.Email = googleUser.Email
		oauthUser.Name = googleUser.Name
		oauthUser.Picture = googleUser.Picture
	} else if provider == "github" {
		var githubUser struct {
			ID        int    `json:"id"`
			Login     string `json:"login"`
			Name      string `json:"name"`
			AvatarURL string `json:"avatar_url"`
			Email     string `json:"email"` // Might be empty if private
		}
		if err := json.Unmarshal(body, &githubUser); err != nil {
			return domain.rest{}, err
		}
		oauthUser.ID = fmt.Sprintf("%d", githubUser.ID)
		oauthUser.Name = githubUser.Name
		if oauthUser.Name == "" {
			oauthUser.Name = githubUser.Login
		}
		oauthUser.Picture = githubUser.AvatarURL
		oauthUser.Email = githubUser.Email

		// If email is not public, we need to fetch it separately
		if oauthUser.Email == "" {
			email, err := s.getGitHubEmail(token.AccessToken)
			if err == nil {
				oauthUser.Email = email
			}
		}
	}

	return oauthUser, nil
}

func (s *OAuthService) getGitHubEmail(accessToken string) (string, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := client.Do(req)
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
	if len(emails) > 0 {
		return emails[0].Email, nil
	}
	return "", errors.New("no email found")
}

func (s *OAuthService) authenticateOAuthUser(provider string, userInfo domain.rest) (domain.rest, error) {
	// 1. Try to find user by OAuth ID
	user, err := s.repo.GetUserByOAuth(provider, userInfo.ID)
	if err == nil {
		// User found, generate tokens
		return s.authService.GenerateTokensForUser(user.ID)
	}

	// 2. If not found by OAuth, try to find by email (to link accounts)
	if userInfo.Email != "" {
		user, err = s.repo.GetUserByEmail(userInfo.Email)
		if err == nil {
			// User exists with this email, link OAuth account?
			// For now, we just return error or maybe we should auto-link?
			// Let's auto-link logic: update user with oauth info
			// TODO: Implement account linking
			// For now, we treat it as found user and generate tokens
			return s.authService.GenerateTokensForUser(user.ID)
		}
	}

	// 3. Create new user
	providerStr := provider
	oauthIDStr := userInfo.ID
	newUser := domain.rest{
		Email:         userInfo.Email,
		Password:      "", // No password for OAuth
		OAuthProvider: &providerStr,
		OAuthID:       &oauthIDStr,
	}

	id, err := s.repo.CreateOAuthUser(newUser)
	if err != nil {
		return domain.rest{}, err
	}

	// Create initial settings
	if err := s.authService.settingsService.CreateInitialUserSettings(id, userInfo.Name); err != nil {
		// Log error but continue?
	}

	return s.authService.GenerateTokensForUser(id)
}
