package main

import (
	"log"
	"xyz-multifinance/config"
	"xyz-multifinance/database"
	"xyz-multifinance/logger"
	"xyz-multifinance/routing"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load Config
	cfg := config.Load()

	// Init Logger
	logger.Setup()

	// Init Database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Init Router
	r := gin.Default()

	// Register Routes
	routing.SetupRoutes(r, db, cfg)

	// Run Server
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
