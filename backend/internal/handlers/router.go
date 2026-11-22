package handlers

import (
	"klinik/internal/middleware"

	"github.com/gin-gonic/gin"
)

var jwtSecret []byte

// RegisterRoutes memasang seluruh endpoint /api/v1. secret dipakai JWT.
func RegisterRoutes(r *gin.Engine, secret string) {
	jwtSecret = []byte(secret)
	auth := r.Group("/api/v1/auth")
	auth.POST("/register", Register)
	auth.POST("/login", Login)
	auth.POST("/logout", Logout)
	auth.GET("/me", middleware.AuthRequired(jwtSecret), Me)

	// Terproteksi petugas; handler CRUD asli datang di task berikutnya.
	admin := r.Group("/api/v1/admin", middleware.AuthRequired(jwtSecret), middleware.RoleRequired("petugas"))
	admin.GET("/dokter", func(c *gin.Context) { c.JSON(200, gin.H{}) })
}
