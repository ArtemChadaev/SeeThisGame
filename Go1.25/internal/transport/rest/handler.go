package rest

import (
	"github.com/ArtemChadaev/SeeThisGame/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	services *service.Service
	redis    *redis.Client
}

func NewHandler(services *service.Service, redis *redis.Client) *Handler {
	return &Handler{
		services: services,
		redis:    redis,
	}
}

// InitRoutes настраивает маршруты приложения
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	// Группа авторизации с ограничением по IP
	auth := router.Group("/auth", h.authRateLimiter)
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/refresh", h.updateToken)
	}

	// Группа API с проверкой токена и лимитом запросов
	api := router.Group("/api", h.userIdentify, h.rateLimiter)
	{
		settings := api.Group("/settings")
		{
			settings.GET("/", h.getMySettings)
			settings.PUT("/", h.setNameIcon)
			settings.POST("/dayCoin", h.dayCoin)
			// settings.POST("/subscript", h.subscribe) // Добавь хендлер, когда будет готов
		}
	}

	return router
}
