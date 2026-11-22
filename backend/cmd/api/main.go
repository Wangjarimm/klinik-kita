package main

import (
	"log"

	"klinik/internal/config"
	"klinik/internal/database"
	"klinik/internal/handlers"
	"klinik/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	if err := database.Connect(cfg.DBDSN); err != nil {
		log.Fatal(err)
	}
	middleware.CookieSameSite = cfg.CookieSameSite

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), middleware.CORS(cfg.WebOrigin))
	_ = r.SetTrustedProxies(nil)
	handlers.RegisterRoutes(r, cfg.JWTSecret)

	log.Printf("api listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
