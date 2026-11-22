package handlers

import (
	"klinik/internal/database"
	"klinik/internal/middleware"
	"klinik/internal/models"

	"github.com/gin-gonic/gin"
)

var jwtSecret []byte

// RegisterRoutes memasang seluruh endpoint /api/v1. secret dipakai JWT.
func RegisterRoutes(r *gin.Engine, secret string) {
	jwtSecret = []byte(secret)

	r.GET("/api/v1/health", Health)

	// Hubungkan hook broadcast antrean (booking baru, pembatalan).
	onCreatedReservasi = func(res *models.Reservasi) {
		var j models.Jadwal
		if database.DB.First(&j, res.JadwalID).Error == nil {
			broadcastAntrean(j.DokterID, res.TanggalKunjungan)
		}
	}
	hookAntreanChange = broadcastAntrean

	// Antrean live: WS publik read-only; endpoint pengubahan di belakang role.
	r.GET("/ws/antrean", WSAntrean)
	r.GET("/ws/notifikasi", middleware.AuthRequired(jwtSecret), WSNotifikasi)

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
	pasienAPI.GET("/reservasi/saya", ReservasiSaya)
	pasienAPI.PATCH("/reservasi/:id/batal", BatalReservasi)

	// Invoice: pemilik reservasi (pasien) atau petugas.
	authed := r.Group("/api/v1", middleware.AuthRequired(jwtSecret))
	authed.GET("/reservasi/invoice/:id", InvoiceReservasi)

	admin := r.Group("/api/v1/admin", middleware.AuthRequired(jwtSecret), middleware.RoleRequired(models.RolePetugas))
	admin.POST("/dokter", AdminCreateDokter)
	admin.GET("/dokter", AdminListDokter)
	admin.PUT("/dokter/:id", AdminUpdateDokter)
	admin.DELETE("/dokter/:id", AdminDeleteDokter)
	admin.POST("/dokter/:id/jadwal", AdminCreateJadwal)
	admin.PUT("/jadwal/:id", AdminUpdateJadwal)
	admin.DELETE("/jadwal/:id", AdminDeleteJadwal)
	admin.POST("/antrean/panggil", PanggilBerikutnya)
	admin.PATCH("/reservasi/:id/status", UbahStatusReservasi)
	admin.POST("/penutupan", AdminCreatePenutupan)
	admin.GET("/penutupan", AdminListPenutupan)
	admin.DELETE("/penutupan/:id", AdminDeletePenutupan)
	admin.POST("/pasien", AdminCreatePasien)
	admin.GET("/pasien", AdminListPasien)
	admin.GET("/statistik", StatistikAdmin)
	admin.GET("/reservasi", AdminListReservasi)
	admin.GET("/pasien/:id/riwayat", RiwayatPasien)
	admin.PUT("/pasien/:id", AdminUpdatePasien)
	admin.DELETE("/pasien/:id", AdminDeletePasien)
}
