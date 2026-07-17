package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
	"sync"
)

const (
	localConfigPath  = "./config/env/.env"
	deployConfigPath = "/.env"
)

type Config struct {
	GigachatAiConfig GigachatAiConfig
	Postgres         Postgres
	HTTP             HTTP
	Redis            Redis
	RateLimiter      RateLimiter
}

type RateLimiter struct {
	TTL   int `env:"RATE_LIMIT_TTL"`
	Limit int `env:"RATE_LIMIT"`
}
type Redis struct {
	Host string `env:"REDIS_HOST"`
	Port string `env:"REDIS_PORT"`
	Pass string `env:"REDIS_PASS"`
}
type Postgres struct {
	Dsn string `env:"POSTGRES_DSN"`
}

type HTTP struct {
	Port string `env:"HTTP_PORT"`
}
type GigachatAiConfig struct {
	AuthKey string `env:"GIGACHAT_AUTH_KEY"`
	Scope   string `env:"GIGACHAT_SCOPE"`

	OAuthUrl string `env:"GIGACHAT_OAUTH_URL"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	deploy := os.Getenv("DEPLOY") == "true"
	once.Do(func() {
		if !deploy {
			err := godotenv.Load(localConfigPath)
			if err != nil {
				panic("Конфигурационный local .env файл не найден")
			}
		} else {
			err := godotenv.Load(deployConfigPath)
			if err != nil {
				panic("Конфигурационный deploy .env файл не найден")
			}
		}

		instance = &Config{}
		if err := cleanenv.ReadEnv(instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			log.Println(help)
			panic(err)
		}
	})

	return instance
}
