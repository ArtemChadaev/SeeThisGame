package rest

import (
	"errors"
	"net/http"
	"path/filepath"

	"github.com/ArtemChadaev/SeeThisGame/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func getUserID(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return 0, domain.ErrInvalidToken // Константа ошибки из domain
	}

	idInt, ok := id.(int)
	if !ok {
		return 0, domain.ErrInvalidToken
	}

	return idInt, nil
}

func (h *Handler) getMySettings(c *gin.Context) {
	userId, err := getUserID(c)
	if err != nil {
		handleError(c, err)
		return
	}

	settings, err := h.services.UserSettingsService.GetByUserID(userId)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, settings)
}

func (h *Handler) setNameIcon(c *gin.Context) {
	userId, err := getUserID(c)
	if err != nil {
		handleError(c, err)
		return
	}

	newName := c.PostForm("name")
	if newName == "" {
		handleError(c, domain.NewInvalidRequestError(errors.New("name field is empty")))
		return
	}

	iconUrl := ""
	file, err := c.FormFile("icon")
	if err == nil {
		ext := filepath.Ext(file.Filename)
		uniqueFilename := uuid.New().String() + ext
		savePath := filepath.Join("static", "icons", uniqueFilename)

		if err := c.SaveUploadedFile(file, savePath); err != nil {
			handleError(c, domain.ErrFailedSaveImg)
			return
		}
		iconUrl = "/static/icons/" + uniqueFilename
	} else if !errors.Is(err, http.ErrMissingFile) {
		handleError(c, domain.NewInternalServerError(err))
		return
	}

	if err := h.services.UserSettingsService.UpdateInfo(userId, newName, iconUrl); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Профиль успешно обновлен",
		"iconUrl": iconUrl,
	})
}

// Исправлено: ресивер Handler вместо http2.Handler
func (h *Handler) dayCoin(c *gin.Context) {
	userId, err := getUserID(c)
	if err != nil {
		handleError(c, err)
		return
	}

	if err = h.services.UserSettingsService.GetGrantDailyReward(userId); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Награда получена"})
}
