package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ArtemChadaev/SeeThisGame/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ErrorResponse struct {
	ErrorField       string `json:"error"`
	ErrorDescription string `json:"error_description,omitempty"`
}

func handleError(c *gin.Context, err error) {
	var appErr *domain.AppError // Используем AppError из domain

	if errors.As(err, &appErr) {
		logMessage := appErr.Message
		if appErr.Err != nil {
			logMessage = fmt.Sprintf("%s: %v", appErr.Message, appErr.Err)
		}
		logrus.Error(logMessage)

		c.AbortWithStatusJSON(appErr.HTTPStatus, ErrorResponse{
			ErrorField:       appErr.Code,
			ErrorDescription: appErr.Message,
		})
	} else {
		logrus.Errorf("unexpected error: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{
			ErrorField:       "internal_server_error",
			ErrorDescription: "An internal server error occurred.",
		})
	}
}
