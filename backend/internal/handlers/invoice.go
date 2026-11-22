package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"klinik/internal/database"
	"klinik/internal/middleware"
	"klinik/internal/models"
	"klinik/internal/services"

	"github.com/gin-gonic/gin"
)

func InvoiceReservasi(c *gin.Context) {
	claims := c.MustGet("claims").(*middleware.Claims)

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reservasi tidak ditemukan"})
		return
	}

	data, err := services.MuatInvoice(database.DB, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reservasi tidak ditemukan"})
		return
	}

	// Petugas boleh lihat semua; pasien hanya miliknya.
	if claims.Role != models.RolePetugas {
		var pasien models.Pasien
		if err := database.DB.Where("user_id = ?", claims.UserID).First(&pasien).Error; err != nil ||
			pasien.ID != data.Reservasi.PasienID {
			c.JSON(http.StatusNotFound, gin.H{"error": "Reservasi tidak ditemukan"})
			return
		}
	}

	pdf, err := services.RenderInvoice(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat invoice"})
		return
	}
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="invoice-%s.pdf"`, data.Reservasi.Kode))
	c.Data(http.StatusOK, "application/pdf", pdf)
}
