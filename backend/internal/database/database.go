package database

import (
	"log"

	"klinik/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(dsn string) error {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	if err := db.AutoMigrate(
		&models.User{}, &models.Pasien{}, &models.Dokter{},
		&models.Jadwal{}, &models.Penutupan{}, &models.Reservasi{},
	); err != nil {
		return err
	}
	DB = db
	Seed()
	return nil
}

// Seed idempoten: akun petugas default + dokter & jadwal contoh (bukan data nyata).
func Seed() {
	var count int64
	DB.Model(&models.User{}).Count(&count)
	if count > 0 {
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte("petugas123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	petugas := models.User{Username: "petugas", Email: "petugas@klinik.local", Password: string(hash), Role: "petugas"}
	if err := DB.Create(&petugas).Error; err != nil {
		log.Fatal(err)
	}

	dokter := []models.Dokter{
		{NID: "D001", Nama: "dr. Ani Wulandari", Keahlian: "Anak", NoHP: "081200000001"},
		{NID: "D002", Nama: "dr. Budi Santoso", Keahlian: "Poli Umum", NoHP: "081200000002"},
	}
	if err := DB.Create(&dokter).Error; err != nil {
		log.Fatal(err)
	}
	var jadwal []models.Jadwal
	for _, d := range dokter {
		for hari := 1; hari <= 5; hari++ {
			jadwal = append(jadwal, models.Jadwal{DokterID: d.ID, Hari: hari, JamMulai: "08:00", JamSelesai: "12:00", Kuota: 3})
		}
	}
	if err := DB.Create(&jadwal).Error; err != nil {
		log.Fatal(err)
	}
}
