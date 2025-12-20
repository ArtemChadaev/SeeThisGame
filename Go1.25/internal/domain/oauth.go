package domain

type OAuthService interface {
	GetAuthURL(provider string) (string, error)
	HandleCallback(provider, code string) (ResponseTokens, error)
}

// OAuthProvider represents supported OAuth providers
type OAuthProvider string

const (
	OAuthProviderGoogle OAuthProvider = "google"
	OAuthProviderGitHub OAuthProvider = "github"
	OAuthProviderLocal  OAuthProvider = "local"
)

// OAuthConfig holds OAuth 2.0 configuration for a provider
type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
	AuthURL      string
	TokenURL     string
}

// OAuthUserInfo contains user information from OAuth provider
type OAuthUserInfo struct {
	Provider OAuthProvider
	ID       string
	Email    string
	Name     string
	Picture  string
}

// OAuthCallbackRequest represents the callback data from OAuth provider
type OAuthCallbackRequest struct {
	Code  string `form:"code" binding:"required"`
	State string `form:"state" binding:"required"`
}
