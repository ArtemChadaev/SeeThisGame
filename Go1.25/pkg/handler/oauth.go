package handler

import (
	"net/http"

	"github.com/ArtemChadaev/go"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initiateOAuth(c *gin.Context) {
	provider := c.Param("provider")
	if provider != "google" && provider != "github" {
		handleError(c, rest.NewInvalidRequestError(nil))
		return
	}

	url, err := h.services.OAuth.GetAuthURL(provider)
	if err != nil {
		handleError(c, err)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *Handler) oauthCallback(c *gin.Context) {
	provider := c.Param("provider")
	if provider != "google" && provider != "github" {
		handleError(c, rest.NewInvalidRequestError(nil))
		return
	}

	var input rest.OAuthCallbackRequest
	if err := c.BindQuery(&input); err != nil {
		handleError(c, rest.NewInvalidRequestError(err))
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
