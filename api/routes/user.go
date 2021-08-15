package routes

import (
	"github.com/LinggaAskaEdo/gin-gorm-clean-arch/api/controllers"
	"github.com/LinggaAskaEdo/gin-gorm-clean-arch/api/middlewares"
	"github.com/LinggaAskaEdo/gin-gorm-clean-arch/lib"
)

// UserRoutes struct
type UserRoutes struct {
	logger         lib.Logger
	handler        lib.RequestHandler
	userController controllers.UserController
	authMiddleware middlewares.JWTAuthMiddleware
}

// Setup user routes
func (s UserRoutes) Setup() {
	s.logger.Info("Setting up routes")
	api := s.handler.Gin.Group("/api").Use(s.authMiddleware.Handler())
	{
		api.POST("/logout", s.userController.Logout)
		api.GET("/user/:id", s.userController.GetOneUser)
		api.GET("/user", s.userController.GetUser)
		api.POST("/user", s.userController.SaveUser)
	}
}

// NewUserRoutes creates new user controller
func NewUserRoutes(
	logger lib.Logger,
	handler lib.RequestHandler,
	userController controllers.UserController,
	authMiddleware middlewares.JWTAuthMiddleware,
) UserRoutes {
	return UserRoutes{
		handler:        handler,
		logger:         logger,
		userController: userController,
		authMiddleware: authMiddleware,
	}
}
