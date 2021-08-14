package lib

import (
	"github.com/go-redis/redis"
)

// Redis modal
type Redis struct {
	*redis.Client
}

// NewRedis creates a new redis instance
func NewRedis(env Env, zapLogger Logger) Redis {
	redisURL := env.RedisUrl
	redisPass := env.RedisPassword
	redisDb := env.RedisDB

	client := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: redisPass,
		DB:       redisDb,
	})

	_, err := client.Ping().Result()
	if err != nil {
		zapLogger.Info("Url: ", redisURL)
		zapLogger.Panic("Can't connect to Redis")
	}

	zapLogger.Info("Redis connection established")

	return Redis{
		Client: client,
	}
}
