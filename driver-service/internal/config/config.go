package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	MongoURI string
	DBName   string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("INFO: .env file not read, checking environment variables...")
	}

	return &Config{
		Port:     getEnv("PORT", "8080"),
		MongoURI: getEnv("MONGO_URI", "mongodb://localhost:27017"),
		DBName:   getEnv("DB_NAME", "bitaksi_driver_db"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
