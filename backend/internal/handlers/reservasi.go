package handlers

import (
	"net/http"

	"klinik/internal/database"
	"klinik/internal/middleware"
	"klinik/internal/models"
	"klinik/internal/services"

	"github.com/gin-gonic/gin"
)

type reservasiInput struct {
	JadwalID uint   `json:"jadwal_id" binding:"required"`
	Tanggal  string `json:"tanggal" binding:"required"`
	Keluhan  string `json:"keluhan" binding:"required,max=500"`
}

func BuatReservasi(c *gin.Context) {
	claims := c.MustGet("claims").(*middleware.Claims)

	var pasien models.Pasien
	if err := database.DB.Where("user_id = ?", claims.UserID).First(&pasien).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profil pasien tidak ditemukan"})
		return
	}

	var in reservasiInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Data tidak valid", "detail": err.Error()})
		return
	}

	svc := &services.BookingService{}
	r, err := svc.Create(database.DB, pasien.ID, services.InputBooking{
		JadwalID: in.JadwalID, Tanggal: in.Tanggal, Keluhan: in.Keluhan,
	})
	if err != nil {
		be, ok := err.(*services.BookingErr)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Kesalahan server"})
			return
		}
		status := http.StatusConflict
		switch be.Kind {
		case "tanggal_lampau", "lewat_batas", "tidak_praktik", "format_tanggal":
			status = http.StatusUnprocessableEntity
		}
		c.JSON(status, gin.H{"error": be.Pesan, "kind": be.Kind})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": r})
}
