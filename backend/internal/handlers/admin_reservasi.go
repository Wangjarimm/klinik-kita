package handlers

import (
	"net/http"

	"klinik/internal/database"
	"klinik/internal/models"

	"github.com/gin-gonic/gin"
)

// AdminListReservasi: daftar reservasi per tanggal (opsional dokter/status), join dokter+jadwal.
func AdminListReservasi(c *gin.Context) {
	tanggal := c.Query("tanggal")
	if tanggal == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Parameter tanggal wajib"})
		return
	}

	db := database.DB.Model(&models.Reservasi{}).
		Joins("JOIN jadwals j ON j.id = reservasis.jadwal_id").
		Where("reservasis.tanggal_kunjungan = ?", tanggal)
	if d := c.Query("dokter_id"); d != "" {
		db = db.Where("j.dokter_id = ?", d)
	}
	if s := c.Query("status"); s != "" {
		db = db.Where("reservasis.status = ?", s)
	}

	var list []models.Reservasi
	if err := db.Order("reservasis.nomor_antrian ASC").Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kesalahan server"})
		return
	}

	// kumpulkan dokter + jadwal untuk nama
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

	var pasienIDs []uint
	for _, r := range list {
		pasienIDs = append(pasienIDs, r.PasienID)
	}
	var pasiens []models.Pasien
	if len(pasienIDs) > 0 {
		database.DB.Where("id IN ?", pasienIDs).Find(&pasiens)
	}
	pasienMap := map[uint]models.Pasien{}
	for _, p := range pasiens {
		pasienMap[p.ID] = p
	}

	data := make([]gin.H, 0, len(list))
	for _, r := range list {
		j := jadwalMap[r.JadwalID]
		d := dokterMap[j.DokterID]
		p := pasienMap[r.PasienID]
		data = append(data, gin.H{
			"id": r.ID, "kode": r.Kode, "tanggal_kunjungan": r.TanggalKunjungan,
			"keluhan": r.Keluhan, "status": r.Status, "nomor_antrian": r.NomorAntrian,
			"catatan_admin": r.CatatanAdmin,
			"pasien":        gin.H{"id": p.ID, "nama_lengkap": p.NamaLengkap},
			"dokter":        gin.H{"id": d.ID, "nama": d.Nama, "keahlian": d.Keahlian},
			"jadwal":        gin.H{"hari": j.Hari, "jam_mulai": j.JamMulai, "jam_selesai": j.JamSelesai},
		})
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}
