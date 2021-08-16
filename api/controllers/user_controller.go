package controllers

import (
	"net/http"
	"strconv"

	"github.com/LinggaAskaEdo/gin-gorm-clean-arch/constants"
	"github.com/LinggaAskaEdo/gin-gorm-clean-arch/lib"
	dto "github.com/LinggaAskaEdo/gin-gorm-clean-arch/models/dto"
	entity "github.com/LinggaAskaEdo/gin-gorm-clean-arch/models/entity"
	"github.com/LinggaAskaEdo/gin-gorm-clean-arch/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserController data type
type UserController struct {
	authService services.JWTAuthService
	service     services.UserService
	logger      lib.Logger
}

// NewUserController creates new user controller
func NewUserController(authService services.JWTAuthService, userService services.UserService, logger lib.Logger) UserController {
	return UserController{
		authService: authService,
		service:     userService,
		logger:      logger,
	}
}

// Logout user
func (u UserController) Logout(c *gin.Context) {
	u.logger.Info("Logout route called")

	authToken, _ := u.authService.ExtractToken(c.Request.Header.Get("Authorization"))

	au, err := u.authService.ExtractTokenMetadata(authToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  err.Error()})
		return
	}

	deleted, delErr := u.service.DeleteToken(au.AccessUUID, au.RefreshUUID)
	if delErr != nil || deleted == 0 { //if any goes wrong
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Successfully logged out"})
}

// GetOneUser gets one user
func (u UserController) GetOneUser(c *gin.Context) {
	u.logger.Info("GetOneUser route called")

	paramID := c.Param("id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	user, err := u.service.GetOneUser(uint(id))

	if err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

// GetUser gets the user
func (u UserController) GetUser(c *gin.Context) {
	u.logger.Info("GetUser route called")

	users, err := u.service.GetAllUser()
	if err != nil {
		u.logger.Error(err)

	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

// SaveUser saves the user
func (u UserController) SaveUser(c *gin.Context) {
	u.logger.Info("SaveUser route called")

	request := dto.Request{}
	trxHandle := c.MustGet(constants.DBTransaction).(*gorm.DB)

	if err := c.ShouldBindJSON(&request); err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	user := entity.User{
		Name:     request.Name,
		Email:    request.Email,
		Age:      request.Age,
		Birthday: request.Birthday}

	result, err := u.service.WithTrx(trxHandle).CreateUser(user)
	if err != nil {
		u.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}
