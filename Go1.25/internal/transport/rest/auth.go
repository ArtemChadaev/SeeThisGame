package rest

import (
	"net/http"

	"github.com/ArtemChadaev/SeeThisGame/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ---  Методы логина для сайта, не менялись ---

func (h *Handler) signUp(c *gin.Context) {
	var input domain.User
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

// --- Выбор персонажа ---

// selectCharacter вызывается, когда юзер уже залогинен на сайте и создает токены для конкретного персонажа в игре
func (h *Handler) selectCharacter(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		handleError(c, err)
		return
	}

	var input struct {
		CharacterID uuid.UUID `json:"characterId" binding:"required"`
	}

	if err := c.BindJSON(&input); err != nil {
		handleError(c, domain.NewInvalidRequestError(err))
		return
	}

	// Генерируем новые токены, в которых уже будет зашит GameUserID
	tokens, err := h.services.AuthorizationService.GenerateGameTokens(userId, input.CharacterID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// --- Обновлено: Универсальный Refresh ---

func (h *Handler) updateToken(c *gin.Context) {
	var input struct {
		RefreshToken string    `json:"refreshToken" binding:"required"`
		GameUserID   uuid.UUID `json:"gameUserId"` // Это поле опционально
	}

	if err := c.BindJSON(&input); err != nil {
		handleError(c, domain.NewInvalidRequestError(err))
		return
	}

	// Если в JSON пришел gameUserId, он передастся в сервис. 
	// Если не пришел — передастся uuid.Nil (нулевой UUID), и токен будет "сайтовым".
	tokens, err := h.services.AuthorizationService.GetAccessToken(input.RefreshToken, input.GameUserID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, tokens)
}

// getUserId извлекает ID пользователя из контекста Gin.
func (h *Handler) getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get("userId")
	if !ok {
		// Если ID нет, значит middleware не сработал или токен пустой
		return 0, domain.ErrInvalidToken 
	}

	idInt, ok := id.(int)
	if !ok {
		return 0, domain.ErrInvalidToken
	}

	return idInt, nil
}

// getGameUserId извлекает UUID персонажа из контекста.
// Возвращает uuid.Nil и ошибку, если персонаж не выбран.
func (h *Handler) getGameUserId(c *gin.Context) (uuid.UUID, error) {
	id, ok := c.Get("gameUserId")
	if !ok {
		return uuid.Nil, domain.ErrCharacterNotSelected
	}

	idUUID, ok := id.(uuid.UUID)
	if !ok {
		return uuid.Nil, domain.ErrInvalidToken
	}

	return idUUID, nil
}