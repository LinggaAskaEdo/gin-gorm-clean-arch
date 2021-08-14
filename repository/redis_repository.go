package repository

import "github.com/LinggaAskaEdo/gin-gorm-clean-arch/lib"

// RedisRepository database structure
type RedisRepository struct {
	lib.Redis
	logger lib.Logger
}

// NewUserRepository creates a new user repository
func NewRedisRepository(redis lib.Redis, logger lib.Logger) RedisRepository {
	return RedisRepository{
		Redis:  redis,
		logger: logger,
	}
}
