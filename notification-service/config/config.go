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

	Kafka struct {
		Host  string
		Port  string
		Group string
		Topic string
	}
}

var AppConfig Config

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed to load .env file")
	}
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, relying on environment variables")
	}

	AppConfig = Config{
		Http: struct {
			Host string
			Port string
		}{
			Host: os.Getenv("HTTP_HOST"),
			Port: os.Getenv("HTTP_PORT"),
		},
		Kafka: struct {
			Host string
			Port string
			Group string
			Topic string
		}{
			Host:  os.Getenv("KAFKA_HOST"),
			Port:  os.Getenv("KAFKA_PORT"),
			Group: os.Getenv("KAFKA_GROUP"),
			Topic: os.Getenv("KAFKA_TOPIC"),
		},
	}

	log.Printf("Loaded config: %+v\n", AppConfig)
}
