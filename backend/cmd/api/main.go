package main

import (
	"log"

	"klinik/internal/config"
	"klinik/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	r := gin.Default()

	api := r.Group("/api/v1")
	api.GET("/health", handlers.Health)

	log.Printf("api listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
