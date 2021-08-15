package controllers

import (
	"net/http"
	"time"

	"github.com/LinggaAskaEdo/gin-gorm-clean-arch/constants"
	"github.com/LinggaAskaEdo/gin-gorm-clean-arch/lib"
	dto "github.com/LinggaAskaEdo/gin-gorm-clean-arch/models/dto"
	entity "github.com/LinggaAskaEdo/gin-gorm-clean-arch/models/entity"
	"github.com/LinggaAskaEdo/gin-gorm-clean-arch/services"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
)

// JWTAuthController struct
type JWTAuthController struct {
	logger      lib.Logger
	service     services.JWTAuthService
	userService services.UserService
}

// NewJWTAuthController creates new controller
func NewJWTAuthController(logger lib.Logger, service services.JWTAuthService, userService services.UserService) JWTAuthController {
	return JWTAuthController{
		logger:      logger,
		service:     service,
		userService: userService,
	}
}

// Login signs in user
func (jwt JWTAuthController) Login(c *gin.Context) {
	jwt.logger.Info("Login route called")

	var req dto.Request

	if err := c.ShouldBindJSON(&req); err != nil {
		jwt.logger.Error(err.Error())
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  "Invalid json provided"})
		return
	}

	type LoginValidation struct {
		Email    string `validate:"required,email"`
		Password string `validate:"required"`
	}

	loginValidation := &LoginValidation{Email: req.Email, Password: req.Password}

	err := validator.New().Struct(loginValidation)
	if err != nil {
		jwt.logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error()})
		return
	}

	user, err := jwt.userService.GetUserByEmailAndPassword(req.Email, req.Password)
	if err != nil {
		jwt.logger.Error(err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  "Please provide valid login details"})
		return
	}

	token, err := jwt.service.CreateToken(user)
	if err != nil {
		jwt.logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error()})
		return
	}

	saveToken := jwt.service.StoreToken(user, *token)
	if saveToken != nil {
		jwt.logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       http.StatusOK,
		"message":      constants.RequestSuccess,
		"accessToken":  token.AccessToken,
		"refreshToken": token.RefreshToken,
	})
}

// Register registers user
func (jwt JWTAuthController) Register(c *gin.Context) {
	jwt.logger.Info("Register route called")

	var req dto.Request

	if err := c.ShouldBindJSON(&req); err != nil {
		jwt.logger.Error(err.Error())
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  "Invalid json provided"})
		return
	}

	type RegisterValidation struct {
		Name     string    `validate:"required,min=5,max=50"`
		Email    string    `validate:"required,email"`
		Password string    `validate:"required,min=5"`
		Age      uint8     `validate:"required,min=17,max=45"`
		Birthday time.Time `validate:"required"`
	}

	registerValidation := &RegisterValidation{Name: req.Name, Email: req.Email, Password: req.Password, Age: req.Age, Birthday: req.Birthday}

	err := validator.New().Struct(registerValidation)
	if err != nil {
		jwt.logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error()})
		return
	}

	user := entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Age:      req.Age,
		Birthday: req.Birthday}

	userResult, err := jwt.userService.CreateUser(user)
	if err != nil {
		jwt.logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": constants.RequestSuccess,
		"data":    userResult,
	})
}
