package api

import (
	"net/http"

	"github.com/ArtemChadaev/SeeThisGame/internal/domain"
	"github.com/ArtemChadaev/SeeThisGame/internal/transport/rest"
	"github.com/gin-gonic/gin"
)

func (h *rest.Handler) signUp(c *gin.Context) {
	var input domain.rest

	if err := c.BindJSON(&input); err != nil {
		handleError(c, domain.rest.NewInvalidRequestError(err))
		return
	}

	_, err := h.services.CreateUser(input)
	if err != nil {
		handleError(c, err)
		return
	}

	tokens, err := h.services.GenerateTokens(input.Email, input.Password)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, tokens)
}

func (h *rest.Handler) signIn(c *gin.Context) {
	var input domain.rest

	if err := c.BindJSON(&input); err != nil {
		handleError(c, domain.rest.NewInvalidRequestError(err))
		return
	}

	tokens, err := h.services.GenerateTokens(input.Email, input.Password)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *rest.Handler) updateToken(c *gin.Context) {
	var input domain.rest

	if err := c.BindJSON(&input); err != nil {
		handleError(c, domain.rest.NewInvalidRequestError(err))
		return
	}

	tokens, err := h.services.GetAccessToken(input.RefreshToken)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, tokens)
}
