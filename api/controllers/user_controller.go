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
	service services.UserService
	logger  lib.Logger
}

// NewUserController creates new user controller
func NewUserController(userService services.UserService, logger lib.Logger) UserController {
	return UserController{
		service: userService,
		logger:  logger,
	}
}

// GetOneUser gets one user
func (u UserController) GetOneUser(c *gin.Context) {
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
	users, err := u.service.GetAllUser()
	if err != nil {
		u.logger.Error(err)

	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

// SaveUser saves the user
func (u UserController) SaveUser(c *gin.Context) {
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
