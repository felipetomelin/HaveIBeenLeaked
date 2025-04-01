package config

import (
	"github.com/lpernett/godotenv"
	"os"
)

type Config struct {
	PublicHost string
	Port       string
	DbUser     string
	DbPassword string
	DbName     string
	DbAddr     string
}

var Env = appSettings()

func appSettings() Config {
	err := godotenv.Load()
	if err != nil {
		return Config{}
	}
	return Config{
		PublicHost: "http://localhost",
		Port:       "8080",
		DbUser:     "root",
		DbPassword: "password",
		DbAddr:     "localhost:3306",
		DbName:     "",
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
