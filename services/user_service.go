package services

import (
	"github.com/LinggaAskaEdo/gin-gorm-clean-arch/lib"
	entity "github.com/LinggaAskaEdo/gin-gorm-clean-arch/models/entity"
	"github.com/LinggaAskaEdo/gin-gorm-clean-arch/repository"
	"gorm.io/gorm"
)

// UserService service layer
type UserService struct {
	logger     lib.Logger
	repository repository.UserRepository
}

// NewUserService creates a new userservice
func NewUserService(logger lib.Logger, repository repository.UserRepository) UserService {
	return UserService{
		logger:     logger,
		repository: repository,
	}
}

// WithTrx delegates transaction to repository database
func (s UserService) WithTrx(trxHandle *gorm.DB) UserService {
	s.repository = s.repository.WithTrx(trxHandle)
	return s
}

// GetOneUser gets one user
func (s UserService) GetOneUser(id uint) (user entity.User, err error) {
	return user, s.repository.Find(&user, id).Error
}

func (s UserService) GetUserByEmailAndPassword(email string, password string) (user entity.User, err error) {
	return user, s.repository.First(&user, "email = ? AND password = ?", email, password).Error
}

// GetAllUser get all the user
func (s UserService) GetAllUser() (users []entity.User, err error) {
	return users, s.repository.Find(&users).Error
}

// CreateUser call to create the user
func (s UserService) CreateUser(user entity.User) (result entity.User, err error) {
	return result, s.repository.Create(&user).Error
}
