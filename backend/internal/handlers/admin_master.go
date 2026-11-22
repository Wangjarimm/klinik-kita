package handlers

import (
	"net/http"

	"klinik/internal/database"
	"klinik/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ---- Dokter (admin) ----

type dokterInput struct {
	NID      string `json:"nid" binding:"required"`
	Nama     string `json:"nama" binding:"required"`
	Keahlian string `json:"keahlian"`
	NoHP     string `json:"no_hp"`
}

func AdminCreateDokter(c *gin.Context) {
	var in dokterInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Data tidak valid", "detail": err.Error()})
		return
	}
	d := models.Dokter{NID: in.NID, Nama: in.Nama, Keahlian: in.Keahlian, NoHP: in.NoHP}
	if err := database.DB.Create(&d).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "NID sudah dipakai"})
		return
	}
	c.JSON(http.StatusCreated, d)
}

func AdminListDokter(c *gin.Context) {
	var data []models.Dokter
	database.DB.Order("nama").Find(&data)
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func AdminUpdateDokter(c *gin.Context) {
	var d models.Dokter
	if err := database.DB.First(&d, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dokter tidak ditemukan"})
		return
	}
	var in dokterInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Data tidak valid", "detail": err.Error()})
		return
	}
	d.NID, d.Nama, d.Keahlian, d.NoHP = in.NID, in.Nama, in.Keahlian, in.NoHP
	database.DB.Save(&d)
	c.JSON(http.StatusOK, d)
}

func AdminDeleteDokter(c *gin.Context) {
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.Jadwal{}, "dokter_id = ?", c.Param("id")).Error; err != nil {
			return err
		}
		return tx.Delete(&models.Dokter{}, c.Param("id")).Error
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kesalahan server"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Dokter dihapus"})
}

// ---- Jadwal (admin) ----

type jadwalInput struct {
	Hari       int    `json:"hari" binding:"gte=0,lte=6"`
	JamMulai   string `json:"jam_mulai" binding:"required"`
	JamSelesai string `json:"jam_selesai" binding:"required"`
	Kuota      int    `json:"kuota" binding:"gte=1"`
}

func AdminCreateJadwal(c *gin.Context) {
	var d models.Dokter
	if err := database.DB.First(&d, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dokter tidak ditemukan"})
		return
	}
	var in jadwalInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Data tidak valid", "detail": err.Error()})
		return
	}
	if in.Kuota == 0 {
		in.Kuota = 3
	}
	j := models.Jadwal{DokterID: d.ID, Hari: in.Hari, JamMulai: in.JamMulai, JamSelesai: in.JamSelesai, Kuota: in.Kuota}
	if err := database.DB.Create(&j).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kesalahan server"})
		return
	}
	c.JSON(http.StatusCreated, j)
}

func AdminUpdateJadwal(c *gin.Context) {
	var j models.Jadwal
	if err := database.DB.First(&j, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Jadwal tidak ditemukan"})
		return
	}
	var in jadwalInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Data tidak valid", "detail": err.Error()})
		return
	}
	j.Hari, j.JamMulai, j.JamSelesai, j.Kuota = in.Hari, in.JamMulai, in.JamSelesai, in.Kuota
	database.DB.Save(&j)
	c.JSON(http.StatusOK, j)
}

func AdminDeleteJadwal(c *gin.Context) {
	database.DB.Delete(&models.Jadwal{}, c.Param("id"))
	c.JSON(http.StatusOK, gin.H{"message": "Jadwal dihapus"})
}

// ---- Penutupan (admin) ----

type penutupanInput struct {
	DokterID *uint  `json:"dokter_id"`
	Tanggal  string `json:"tanggal" binding:"required"`
	Alasan   string `json:"alasan" binding:"required"`
}

func AdminCreatePenutupan(c *gin.Context) {
	var in penutupanInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Data tidak valid", "detail": err.Error()})
		return
	}
	p := models.Penutupan{DokterID: in.DokterID, Tanggal: in.Tanggal, Alasan: in.Alasan}
	if err := database.DB.Create(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kesalahan server"})
		return
	}
	c.JSON(http.StatusCreated, p)
}

func AdminListPenutupan(c *gin.Context) {
	var data []models.Penutupan
	database.DB.Order("tanggal DESC").Find(&data)
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func AdminDeletePenutupan(c *gin.Context) {
	database.DB.Delete(&models.Penutupan{}, c.Param("id"))
	c.JSON(http.StatusOK, gin.H{"message": "Penutupan dihapus"})
}

// ---- Pasien (admin, termasuk walk-in tanpa akun) ----

type pasienInput struct {
	NamaLengkap  string `json:"nama_lengkap" binding:"required"`
	NIK          string `json:"nik" binding:"required,len=16"`
	JenisKelamin string `json:"jenis_kelamin" binding:"required,oneof=L P"`
	TempatLahir  string `json:"tempat_lahir"`
	TanggalLahir string `json:"tanggal_lahir"`
	Alamat       string `json:"alamat"`
	NoHP         string `json:"no_hp" binding:"required"`
}

func AdminCreatePasien(c *gin.Context) {
	var in pasienInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Data tidak valid", "detail": err.Error()})
		return
	}
	p := models.Pasien{
		NamaLengkap: in.NamaLengkap, NIK: in.NIK, JenisKelamin: in.JenisKelamin,
		TempatLahir: in.TempatLahir, TanggalLahir: in.TanggalLahir, Alamat: in.Alamat, NoHP: in.NoHP,
	}
	if err := database.DB.Create(&p).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "NIK sudah terdaftar"})
		return
	}
	c.JSON(http.StatusCreated, p)
}

func AdminListPasien(c *gin.Context) {
	var data []models.Pasien
	q := c.Query("q")
	db := database.DB.Model(&models.Pasien{})
	if q != "" {
		db = db.Where("nama_lengkap LIKE ? OR nik LIKE ?", "%"+q+"%", "%"+q+"%")
	}
	db.Order("nama_lengkap").Find(&data)
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func AdminUpdatePasien(c *gin.Context) {
	var p models.Pasien
	if err := database.DB.First(&p, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pasien tidak ditemukan"})
		return
	}
	var in pasienInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Data tidak valid", "detail": err.Error()})
		return
	}
	p.NamaLengkap, p.NIK, p.JenisKelamin = in.NamaLengkap, in.NIK, in.JenisKelamin
	p.TempatLahir, p.TanggalLahir, p.Alamat, p.NoHP = in.TempatLahir, in.TanggalLahir, in.Alamat, in.NoHP
	database.DB.Save(&p)
	c.JSON(http.StatusOK, p)
}

func AdminDeletePasien(c *gin.Context) {
	database.DB.Delete(&models.Pasien{}, c.Param("id"))
	c.JSON(http.StatusOK, gin.H{"message": "Pasien dihapus"})
}

// ---- Dokter (publik, untuk landing & pencarian pasien) ----

func ListDokter(c *gin.Context) {
	var data []models.Dokter
	q := c.Query("q")
	db := database.DB.Model(&models.Dokter{})
	if q != "" {
		db = db.Where("nama LIKE ? OR keahlian LIKE ?", "%"+q+"%", "%"+q+"%")
	}
	db.Order("nama").Find(&data)
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func DetailDokter(c *gin.Context) {
	var d models.Dokter
	if err := database.DB.First(&d, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dokter tidak ditemukan"})
		return
	}
	var jadwal []models.Jadwal
	database.DB.Where("dokter_id = ?", d.ID).Order("hari, jam_mulai").Find(&jadwal)
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"dokter": d, "jadwal": jadwal}})
}

func ListJadwalDokter(c *gin.Context) {
	var data []models.Jadwal
	database.DB.Where("dokter_id = ?", c.Param("id")).Order("hari, jam_mulai").Find(&data)
	c.JSON(http.StatusOK, gin.H{"data": data})
}
