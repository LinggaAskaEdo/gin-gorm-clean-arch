package lib

import (
	"log"

	"github.com/spf13/viper"
)

// Env has environment stored
type Env struct {
	ServerPort    string `mapstructure:"SERVER_PORT"`
	Environment   string `mapstructure:"ENV"`
	LogOutput     string `mapstructure:"LOG_OUTPUT"`
	DBUsername    string `mapstructure:"DB_USER"`
	DBPassword    string `mapstructure:"DB_PASS"`
	DBHost        string `mapstructure:"DB_HOST"`
	DBPort        string `mapstructure:"DB_PORT"`
	DBName        string `mapstructure:"DB_NAME"`
	DBMaxOpenConn int    `mapstructure:"DB_MAX_OPEN_CONN"`
	DBMaxIdleConn int    `mapstructure:"DB_MAX_IDLE_CONN"`
	DBMaxLifeTime int    `mapstructure:"DB_MAX_LIFE_TIME"`
	DBMaxIdleTime int    `mapstructure:"DB_MAX_IDLE_TIME"`
	RedisUrl      string `mapstructure:"REDIS_URL"`
	RedisPassword string `mapstructure:"REDIS_PASS"`
	RedisDB       int    `mapstructure:"REDIS_DB"`
	AccessSecret  string `mapstructure:"ACCESS_SECRET"`
	RefreshSecret string `mapstructure:"REFRESH_SECRET"`
}

// NewEnv creates a new environment
func NewEnv() Env {
	env := Env{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("☠️ cannot read configuration")
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("☠️ environment can't be loaded: ", err)
	}

	log.Printf("%+v \n", env)
	return env
}
