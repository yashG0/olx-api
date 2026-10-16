package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	Env  string
}

func MustLoad() Config {
	godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		panic("Port is Required")
	}
	env := os.Getenv("ENV")
	if env == "" {
		panic("Environment is Required")
	}
	return Config{Port: port, Env: env}
}
