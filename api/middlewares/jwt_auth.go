package middlewares

import (
	"net/http"

	"github.com/LinggaAskaEdo/gin-gorm-clean-arch/lib"
	"github.com/LinggaAskaEdo/gin-gorm-clean-arch/services"
	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware middleware for jwt authentication
type JWTAuthMiddleware struct {
	service services.JWTAuthService
	logger  lib.Logger
}

// NewJWTAuthMiddleware creates new jwt auth middleware
func NewJWTAuthMiddleware(logger lib.Logger, service services.JWTAuthService) JWTAuthMiddleware {
	return JWTAuthMiddleware{
		service: service,
		logger:  logger,
	}
}

// Setup sets up jwt auth middleware
func (m JWTAuthMiddleware) Setup() {}

// Handler handles middleware functionality
func (m JWTAuthMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken, err := m.service.ExtractToken(c.Request.Header.Get("Authorization"))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": http.StatusUnauthorized,
				"error":  "You are not authorized",
			})
			c.Abort()
			return
		}

		authorized, err := m.service.AuthorizeToken(authToken)
		if authorized {
			// TODO: add logic check uuid is exist in redis, if exist next, if not exist abort
			c.Next()
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		m.logger.Error(err)
		c.Abort()
	}
}
