package handlers

import (
	"net/http"
	"strconv"
	"time"

	"klinik/internal/database"
	"klinik/internal/models"
	"klinik/internal/services"
	"klinik/internal/ws"

	"github.com/gin-gonic/gin"
)

var (
	hubWS     = ws.NewHub()
	antreanSvc = &services.AntreanService{Hub: hubWS}
)

// WSAntrean publik (read-only): halaman display klinik tidak ber-login, nomor
// antrean bukan data sensitif. Auth tetap dipakai untuk endpoint yang mengubah.
func WSAntrean(c *gin.Context) {
	dokterID, err := strconv.ParseUint(c.Query("dokter_id"), 10, 32)
	if err != nil || c.Query("tanggal") == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Parameter dokter_id dan tanggal wajib"})
		return
	}
	tanggal := c.Query("tanggal")
	if _, err := time.Parse("2006-01-02", tanggal); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Format tanggal harus YYYY-MM-DD"})
		return
	}

	awal, err := antreanSvc.Snapshot(database.DB, uint(dokterID), tanggal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kesalahan server"})
		return
	}
	hubWS.ServeAntrean(c.Writer, c.Request, services.RoomAntrean(uint(dokterID), tanggal), awal)
}

type panggilInput struct {
	DokterID uint   `json:"dokter_id" binding:"required"`
	Tanggal  string `json:"tanggal" binding:"required"`
}

func PanggilBerikutnya(c *gin.Context) {
	var in panggilInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Data tidak valid", "detail": err.Error()})
		return
	}

	r, err := antreanSvc.PanggilBerikutnya(database.DB, in.DokterID, in.Tanggal)
	if err != nil {
		if be, ok := err.(*services.BookingErr); ok && be.Kind == "kosong" {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": be.Pesan})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kesalahan server"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": r})
}

type ubahStatusInput struct {
	Status  string `json:"status" binding:"required,oneof=dikonfirmasi dipanggil selesai batal"`
	Catatan string `json:"catatan"`
}

func UbahStatusReservasi(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reservasi tidak ditemukan"})
		return
	}
	var in ubahStatusInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Data tidak valid", "detail": err.Error()})
		return
	}

	var r models.Reservasi
	if err := database.DB.First(&r, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reservasi tidak ditemukan"})
		return
	}
	r.Status = in.Status
	if in.Catatan != "" {
		r.CatatanAdmin = in.Catatan
	}
	if err := database.DB.Save(&r).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kesalahan server"})
		return
	}

	var j models.Jadwal
	if database.DB.First(&j, r.JadwalID).Error == nil {
		antreanSvc.BroadcastSnapshot(database.DB, j.DokterID, r.TanggalKunjungan)
	}
	c.JSON(http.StatusOK, gin.H{"data": r})
}

func broadcastAntrean(dokterID uint, tanggal string) {
	antreanSvc.BroadcastSnapshot(database.DB, dokterID, tanggal)
}
