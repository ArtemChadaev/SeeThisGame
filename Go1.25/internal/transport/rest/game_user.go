package rest

import (
	"net/http"
	"github.com/ArtemChadaev/SeeThisGame/internal/domain"
	"github.com/gin-gonic/gin"
)

// ШАГ 1: Создание персонажа (с проверкой лимита)
func (h *Handler) createCharacter(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		handleError(c, err)
		return
	}

	var input struct {
		Nickname string `json:"nickname" binding:"required"`
	}
	if err := c.BindJSON(&input); err != nil {
		handleError(c, domain.NewInvalidRequestError(err))
		return
	}

	// Вызываем сервис (в нем зашита логика проверки 5 персонажей)
	charID, err := h.services.InitialCreateCharacter(c.Request.Context(), userId, input.Nickname)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": charID})
}

// Список всех персонажей пользователя
func (h *Handler) getMyCharacters(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		handleError(c, err)
		return
	}

    // Предполагается, что в сервисе есть метод получения списка
	characters, err := h.services.GetMyCharacters(c.Request.Context(), userId)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, characters)
}

// Получение полного профиля (только если персонаж выбран)
func (h *Handler) getGameProfile(c *gin.Context) {
	charID, err := h.getGameUserId(c)
	if err != nil {
		handleError(c, err)
		return
	}

	profile, err := h.services.GetFullProfile(c.Request.Context(), charID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, profile)
}

// ШАГ 2: Обновление настроек UI
func (h *Handler) updateGameSettings(c *gin.Context) {
	charID, err := h.getGameUserId(c)
	if err != nil {
		handleError(c, err)
		return
	}

	var input domain.GameUserSettings
	if err := c.BindJSON(&input); err != nil {
		handleError(c, domain.NewInvalidRequestError(err))
		return
	}

	if err := h.services.UpdateSettings(c.Request.Context(), charID, input); err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

// ШАГ 3: Обновление мира (сеттинг, тон, иерархия)
func (h *Handler) updateGameWorld(c *gin.Context) {
	charID, err := h.getGameUserId(c)
	if err != nil {
		handleError(c, err)
		return
	}

	var input domain.WorldState
	if err := c.BindJSON(&input); err != nil {
		handleError(c, domain.NewInvalidRequestError(err))
		return
	}

	if err := h.services.UpdateWorldState(c.Request.Context(), charID, input); err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}