package rest

import (
	"net/http"

	"github.com/ArtemChadaev/SeeThisGame/internal/domain"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initiateOAuth(c *gin.Context) {
	provider := c.Param("provider")
	if provider != "google" && provider != "github" {
		handleError(c, domain.NewInvalidRequestError(nil))
		return
	}

	url, err := h.services.OAuthService.GetAuthURL(provider)
	if err != nil {
		handleError(c, err)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *Handler) oauthCallback(c *gin.Context) {
	provider := c.Param("provider")
	if provider != "google" && provider != "github" {
		handleError(c, domain.NewInvalidRequestError(nil))
		return
	}

	// Код и стейт приходят в URL (query)
	code := c.Query("code")
	if code == "" {
		handleError(c, domain.NewInvalidRequestError(nil))
		return
	}

	// TODO: Проверка параметра state для защиты от CSRF

	tokens, err := h.services.OAuthService.HandleCallback(provider, code)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, tokens)
}
