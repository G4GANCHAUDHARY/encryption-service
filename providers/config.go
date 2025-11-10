package providers

import (
	"flag"
	env "github.com/caarlos0/env/v9"
	godotenv "github.com/joho/godotenv"
)

type App struct {
	Name    string `env:"APP_NAME"`
	Version string `env:"APP_VERSION"`
}

type HttpConfig struct {
	Address string `env:"LISTEN_ADDRESS"`
}

type RedisConfig struct {
	Address  string `env:"REDIS_LISTEN_ADDRESS"`
	DbNumber int    `env:"REDIS_DB_NUMBER"`
}

type DbConfig struct {
	Host     string `env:"DB_HOST"`
	Port     int    `env:"DB_PORT"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	Name     string `env:"DB_NAME"`
}

type AppConfig struct {
	App         App
	DbConfig    DbConfig
	HttpConfig  HttpConfig
	RedisConfig RedisConfig
}

var (
	config     *AppConfig
	configPath string
)

func loadConfig() (AppConfig, error) {
	if err := godotenv.Load(configPath); err != nil {
		panic(err)
	}

	var cfg AppConfig
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	return cfg, nil
}

func GetConfig(path string) (AppConfig, error) {
	if config == nil {
		flag.StringVar(&configPath, "config", path, "Path to .env file")
		flag.Parse()

		cnf, err := loadConfig()
		if err != nil {
			panic(err)
		}
		config = &cnf
	}
	return *config, nil
}
