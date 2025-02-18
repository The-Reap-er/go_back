package main

import (
	"go_back/api"
	"go_back/config"
	"go_back/database"
	"go_back/log"
)

func main() {
	// Initialize Logger
	log.InitLogger()

	// Load configuration
	cfg := config.LoadConfig()

	// Connect to Database
	database.ConnectDB(cfg)

	// Setup Router
	r := api.SetupRouter(cfg)

	// Start Server
	if err := r.Run(":" + cfg.Port); err != nil {
		log.ErrorLogger.Println("Failed to run server:", err)
	}
}
