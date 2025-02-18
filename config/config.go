package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	MongoURI  string
	MongoDB   string
	JWTSecret string
	ZAPAPIURL string
	ZAPAPIKey string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	return Config{
		Port:      os.Getenv("PORT"),
		MongoURI:  os.Getenv("MONGO_URI"),
		MongoDB:   os.Getenv("MONGO_DB"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		ZAPAPIURL: os.Getenv("ZAP_API_URL"),
		ZAPAPIKey: os.Getenv("ZAP_API_KEY"),
	}
}
