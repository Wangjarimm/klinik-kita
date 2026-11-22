package handlers

import (
	"klinik/internal/middleware"
	"klinik/internal/models"

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

	// Dokter publik: dipakai landing page dan pencarian pasien.
	dokter := r.Group("/api/v1/dokter")
	dokter.GET("", ListDokter)
	dokter.GET("/:id", DetailDokter)
	dokter.GET("/:id/jadwal", ListJadwalDokter)
	dokter.GET("/:id/ketersediaan", KetersediaanDokter)

	// Reservasi pasien (butuh role pasien).
	pasienAPI := r.Group("/api/v1", middleware.AuthRequired(jwtSecret), middleware.RoleRequired(models.RolePasien))
	pasienAPI.POST("/reservasi", BuatReservasi)

	admin := r.Group("/api/v1/admin", middleware.AuthRequired(jwtSecret), middleware.RoleRequired(models.RolePetugas))
	admin.POST("/dokter", AdminCreateDokter)
	admin.GET("/dokter", AdminListDokter)
	admin.PUT("/dokter/:id", AdminUpdateDokter)
	admin.DELETE("/dokter/:id", AdminDeleteDokter)
	admin.POST("/dokter/:id/jadwal", AdminCreateJadwal)
	admin.PUT("/jadwal/:id", AdminUpdateJadwal)
	admin.DELETE("/jadwal/:id", AdminDeleteJadwal)
	admin.POST("/penutupan", AdminCreatePenutupan)
	admin.GET("/penutupan", AdminListPenutupan)
	admin.DELETE("/penutupan/:id", AdminDeletePenutupan)
	admin.POST("/pasien", AdminCreatePasien)
	admin.GET("/pasien", AdminListPasien)
	admin.PUT("/pasien/:id", AdminUpdatePasien)
	admin.DELETE("/pasien/:id", AdminDeletePasien)
}
