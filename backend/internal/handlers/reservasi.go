package handlers

import (
	"net/http"
	"strconv"
	"time"

	"klinik/internal/database"
	"klinik/internal/middleware"
	"klinik/internal/models"
	"klinik/internal/services"

	"github.com/gin-gonic/gin"
)

// hookAntreanChange dipasang layer WebSocket (Task antrean) untuk broadcast perubahan antrean.
var hookAntreanChange func(dokterID uint, tanggal string)

// onCreatedReservasi dipasang ke BookingService: broadcast setelah commit.
var onCreatedReservasi func(r *models.Reservasi)

func pasienDariClaims(c *gin.Context) (*models.Pasien, *middleware.Claims, bool) {
	claims := c.MustGet("claims").(*middleware.Claims)
	var pasien models.Pasien
	if err := database.DB.Where("user_id = ?", claims.UserID).First(&pasien).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profil pasien tidak ditemukan"})
		return nil, nil, false
	}
	return &pasien, claims, true
}

func ReservasiSaya(c *gin.Context) {
	pasien, _, ok := pasienDariClaims(c)
	if !ok {
		return
	}

	var list []models.Reservasi
	database.DB.Where("pasien_id = ?", pasien.ID).
		Order("tanggal_kunjungan DESC, nomor_antrian ASC").Limit(100).Find(&list)

	var jadwalIDs []uint
	for _, r := range list {
		jadwalIDs = append(jadwalIDs, r.JadwalID)
	}
	var jadwals []models.Jadwal
	if len(jadwalIDs) > 0 {
		database.DB.Where("id IN ?", jadwalIDs).Find(&jadwals)
	}
	var dokterIDs []uint
	for _, j := range jadwals {
		dokterIDs = append(dokterIDs, j.DokterID)
	}
	var dokters []models.Dokter
	if len(dokterIDs) > 0 {
		database.DB.Where("id IN ?", dokterIDs).Find(&dokters)
	}
	jadwalMap := map[uint]models.Jadwal{}
	for _, j := range jadwals {
		jadwalMap[j.ID] = j
	}
	dokterMap := map[uint]models.Dokter{}
	for _, d := range dokters {
		dokterMap[d.ID] = d
	}

	data := make([]gin.H, 0, len(list))
	for _, r := range list {
		j := jadwalMap[r.JadwalID]
		d := dokterMap[j.DokterID]
		data = append(data, gin.H{
			"id": r.ID, "kode": r.Kode, "tanggal_kunjungan": r.TanggalKunjungan,
			"keluhan": r.Keluhan, "status": r.Status, "nomor_antrian": r.NomorAntrian,
			"catatan_admin": r.CatatanAdmin,
			"dokter":        gin.H{"id": d.ID, "nama": d.Nama, "keahlian": d.Keahlian},
			"jadwal":        gin.H{"hari": j.Hari, "jam_mulai": j.JamMulai, "jam_selesai": j.JamSelesai},
		})
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func BatalReservasi(c *gin.Context) {
	pasien, _, ok := pasienDariClaims(c)
	if !ok {
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reservasi tidak ditemukan"})
		return
	}

	var r models.Reservasi
	if err := database.DB.First(&r, id).Error; err != nil || r.PasienID != pasien.ID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reservasi tidak ditemukan"})
		return
	}
	if r.Status != models.StatusMenunggu {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Hanya reservasi berstatus menunggu yang bisa dibatalkan"})
		return
	}
	hariIni := time.Now().In(services.WIB).Format("2006-01-02")
	if r.TanggalKunjungan <= hariIni {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Reservasi hari-H tidak bisa dibatalkan lewat aplikasi; hubungi klinik"})
		return
	}

	r.Status = models.StatusBatal
	database.DB.Save(&r)
	if hookAntreanChange != nil {
		var j models.Jadwal
		if database.DB.First(&j, r.JadwalID).Error == nil {
			hookAntreanChange(j.DokterID, r.TanggalKunjungan)
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "Reservasi dibatalkan", "data": r})
}

type reservasiInput struct {
	JadwalID uint   `json:"jadwal_id" binding:"required"`
	Tanggal  string `json:"tanggal" binding:"required"`
	Keluhan  string `json:"keluhan" binding:"required,max=500"`
}

func BuatReservasi(c *gin.Context) {
	pasien, _, ok := pasienDariClaims(c)
	if !ok {
		return
	}

	var in reservasiInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Data tidak valid", "detail": err.Error()})
		return
	}

	svc := &services.BookingService{OnCreated: onCreatedReservasi}
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
