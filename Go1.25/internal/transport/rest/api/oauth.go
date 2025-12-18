package api

import (
	"net/http"

	"github.com/ArtemChadaev/SeeThisGame/internal/domain"
	"github.com/ArtemChadaev/SeeThisGame/internal/transport/rest"
	"github.com/gin-gonic/gin"
)

func (h *rest.Handler) initiateOAuth(c *gin.Context) {
	provider := c.Param("provider")
	if provider != "google" && provider != "github" {
		handleError(c, domain.rest.NewInvalidRequestError(nil))
		return
	}

	url, err := h.services.OAuth.GetAuthURL(provider)
	if err != nil {
		handleError(c, err)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *rest.Handler) oauthCallback(c *gin.Context) {
	provider := c.Param("provider")
	if provider != "google" && provider != "github" {
		handleError(c, domain.rest.NewInvalidRequestError(nil))
		return
	}

	var input domain.rest
	if err := c.BindQuery(&input); err != nil {
		handleError(c, domain.rest.NewInvalidRequestError(err))
		return
	}

	// TODO: Verify state parameter to prevent CSRF

	tokens, err := h.services.OAuth.HandleCallback(provider, input.Code)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, tokens)
}
