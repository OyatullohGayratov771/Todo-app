package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Http struct {
		Host string
		Port string
	}

	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}
}

var AppConfig Config

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed to load .env file")
	}

	AppConfig = Config{
		Http: struct {
			Host string
			Port string
		}{
			Host: os.Getenv("HTTP_HOST"),
			Port: os.Getenv("HTTP_PORT"),
		},
		Database: struct {
			Host     string
			Port     string
			User     string
			Password string
			Name     string
		}{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
		},
	}
}
