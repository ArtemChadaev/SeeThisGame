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

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	// Группа авторизации (публичная: регистрация, вход, рефреш)
	auth := router.Group("/auth", h.authRateLimiter)
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/refresh", h.updateToken)
	}

	// Группа API (защищенная: требует accessToken)
	api := router.Group("/api", h.userIdentify, h.rateLimiter)
	{
		// 1. Управление сессией персонажа
		// Здесь мы получаем "игровой" токен с UUID внутри
		api.POST("/auth/select-character", h.selectCharacter)

		// 2. Системные настройки аккаунта
		settings := api.Group("/settings")
		{
			settings.GET("/", h.getMySettings)
			settings.PUT("/", h.setNameIcon)
			settings.POST("/dayCoin", h.dayCoin)
		}

		// 3. Игровой процесс и персонажи
		game := api.Group("/game")
		{
			// Создание и список
			game.POST("/create", h.createCharacter)    // ШАГ 1: Создать (пустой + лимит 5)
			game.GET("/list", h.getMyCharacters)       // Список всех персонажей юзера

			// Работа с конкретным персонажем (требуют, чтобы персонаж был выбран)
			game.GET("/profile", h.getGameProfile)     // Получить Nickname, Level, World и т.д.
			game.PUT("/settings", h.updateGameSettings) // ШАГ 2: Установка/смена UI настроек
			game.PUT("/world", h.updateGameWorld)       // ШАГ 3: Установка/смена Мира
		}
	}

	return router
}