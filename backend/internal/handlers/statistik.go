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

// StatistikAdmin: seluruh angka berasal dari query database, tanpa nilai dummy.
func StatistikAdmin(c *gin.Context) {
	hariIni := time.Now().In(services.WIB).Format("2006-01-02")
	awal := time.Now().In(services.WIB).AddDate(0, 0, -29).Format("2006-01-02")

	var rowsStatus []struct {
		Status string
		Total  int
	}
	database.DB.Model(&models.Reservasi{}).
		Select("status, COUNT(*) AS total").
		Where("tanggal_kunjungan = ?", hariIni).
		Group("status").Scan(&rowsStatus)
	hariIniMap := map[string]int{
		models.StatusMenunggu: 0, models.StatusDikonfirmasi: 0,
		models.StatusDipanggil: 0, models.StatusSelesai: 0, models.StatusBatal: 0,
	}
	for _, r := range rowsStatus {
		hariIniMap[r.Status] = r.Total
	}

	var rowsHari []struct {
		Tanggal string
		Total   int
	}
	database.DB.Model(&models.Reservasi{}).
		Select("tanggal_kunjungan AS tanggal, COUNT(*) AS total").
		Where("tanggal_kunjungan BETWEEN ? AND ? AND status <> ?", awal, hariIni, models.StatusBatal).
		Group("tanggal_kunjungan").Scan(&rowsHari)
	totalPerTanggal := map[string]int{}
	for _, r := range rowsHari {
		totalPerTanggal[r.Tanggal] = r.Total
	}
	perHari := make([]gin.H, 0, 30)
	for i := 29; i >= 0; i-- {
		tgl := time.Now().In(services.WIB).AddDate(0, 0, -i).Format("2006-01-02")
		perHari = append(perHari, gin.H{"tanggal": tgl, "total": totalPerTanggal[tgl]})
	}

	var dokterRows []map[string]interface{}
	database.DB.Table("dokters d").
		Select("d.nama AS nama, COUNT(r.id) AS total").
		Joins("JOIN jadwals j ON j.dokter_id = d.id").
		Joins("JOIN reservasis r ON r.jadwal_id = j.id AND r.status <> ?", models.StatusBatal).
		Group("d.id, d.nama").Order("total DESC").Limit(5).Scan(&dokterRows)
	dokterTop := make([]gin.H, 0, len(dokterRows))
	for _, row := range dokterRows {
		dokterTop = append(dokterTop, gin.H{"nama": row["nama"], "total": row["total"]})
	}

	c.JSON(http.StatusOK, gin.H{
		"hari_ini":          hariIniMap,
		"per_hari":          perHari,
		"dokter_terpopuler": dokterTop,
	})
}

// RiwayatPasien: reservasi lampau + keluhan + catatan admin (riwayat kunjungan ringan).
func RiwayatPasien(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pasien tidak ditemukan"})
		return
	}
	var pasien models.Pasien
	if err := database.DB.First(&pasien, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pasien tidak ditemukan"})
		return
	}

	var list []models.Reservasi
	database.DB.Where("pasien_id = ?", pasien.ID).
		Order("tanggal_kunjungan DESC, nomor_antrian ASC").Find(&list)

	data := make([]gin.H, 0, len(list))
	for _, r := range list {
		data = append(data, gin.H{
			"id": r.ID, "kode": r.Kode, "tanggal_kunjungan": r.TanggalKunjungan,
			"keluhan": r.Keluhan, "status": r.Status, "nomor_antrian": r.NomorAntrian,
			"catatan_admin": r.CatatanAdmin,
		})
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

// kirimNotifikasiStatus: pesan ke pemilik reservasi saat admin mengubah status.
func kirimNotifikasiStatus(r *models.Reservasi) {
	var pasien models.Pasien
	if err := database.DB.First(&pasien, r.PasienID).Error; err != nil || pasien.UserID == nil {
		return
	}
	pesan := "Reservasi " + r.Kode + " status: " + r.Status
	if r.Status == models.StatusBatal && r.CatatanAdmin != "" {
		pesan += " (" + r.CatatanAdmin + ")"
	}
	hubWS.NotifyUser(*pasien.UserID, ws.NotifikasiMsg{
		Judul: "Status Reservasi",
		Pesan: pesan,
		Tipe:  "info",
	})
}
