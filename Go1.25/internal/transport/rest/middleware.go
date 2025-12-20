package rest

import (
	"context"
	"strings"
	"time"

	"github.com/ArtemChadaev/SeeThisGame/internal/domain"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"

	rateLimitPerMinute = 20
	rateWindow         = 1 * time.Minute

	authRateLimitPerMinute = 10
	authRateWindow         = 1 * time.Minute
)

// userIdentify — проверка валидности Access токена
func (h *Handler) userIdentify(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		handleError(c, domain.ErrInvalidToken) // Используем константу из domain
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		handleError(c, domain.ErrInvalidToken)
		return
	}

	// Вызываем сервис через интерфейс
	userId, err := h.services.AuthorizationService.ParseToken(headerParts[1])
	if err != nil {
		handleError(c, err)
		return
	}

	c.Set(userCtx, userId)
}

// rateLimiter — ограничение частоты запросов по токену через Redis
func (h *Handler) rateLimiter(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		c.Next()
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		c.Next()
		return
	}
	accessToken := headerParts[1]

	ctx := context.Background()
	key := "rate_limit:" + accessToken

	pipe := h.redis.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, rateWindow)
	_, err := pipe.Exec(ctx)

	if err != nil {
		c.Next() // Если Redis упал, не блокируем пользователя
		return
	}

	if incr.Val() > rateLimitPerMinute {
		handleError(c, domain.ErrTooManyRequestsByAccessToken)
		c.Abort()
		return
	}

	c.Next()
}

// authRateLimiter — ограничение запросов к /auth по IP адресу
func (h *Handler) authRateLimiter(c *gin.Context) {
	ip := c.ClientIP()
	key := "rate_limit_auth:" + ip
	ctx := context.Background()

	pipe := h.redis.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, authRateWindow)
	_, err := pipe.Exec(ctx)

	if err != nil {
		c.Next()
		return
	}

	if incr.Val() > authRateLimitPerMinute {
		handleError(c, domain.ErrTooManyRequestsByIp)
		c.Abort()
		return
	}

	c.Next()
}
