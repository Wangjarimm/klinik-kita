package handlers

import (
	"net/http"
	"strconv"
	"time"

	"klinik/internal/database"
	"klinik/internal/services"

	"github.com/gin-gonic/gin"
)

func KetersediaanDokter(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dokter tidak ditemukan"})
		return
	}

	var bulan time.Time
	if b := c.Query("bulan"); b != "" {
		parsed, err := time.Parse("2006-01", b)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Parameter bulan harus format YYYY-MM"})
			return
		}
		bulan = parsed
	} else {
		bulan = time.Now().In(services.WIB)
	}

	hari, err := services.KetersediaanBulan(database.DB, uint(id), bulan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kesalahan server"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"hari": hari})
}
