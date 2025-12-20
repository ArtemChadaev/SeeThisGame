package rest

import (
	"net/http"

	"github.com/ArtemChadaev/SeeThisGame/internal/domain"
	"github.com/gin-gonic/gin"
)

func (h *Handler) signUp(c *gin.Context) {
	var input domain.User // Заменили rest на User

	if err := c.BindJSON(&input); err != nil {
		handleError(c, domain.NewInvalidRequestError(err))
		return
	}

	_, err := h.services.AuthorizationService.CreateUser(input)
	if err != nil {
		handleError(c, err)
		return
	}

	tokens, err := h.services.AuthorizationService.GenerateTokens(input.Email, input.Password)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, tokens)
}

func (h *Handler) signIn(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.BindJSON(&input); err != nil {
		handleError(c, domain.NewInvalidRequestError(err))
		return
	}

	tokens, err := h.services.AuthorizationService.GenerateTokens(input.Email, input.Password)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *Handler) updateToken(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}

	if err := c.BindJSON(&input); err != nil {
		handleError(c, domain.NewInvalidRequestError(err))
		return
	}

	tokens, err := h.services.AuthorizationService.GetAccessToken(input.RefreshToken)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, tokens)
}
