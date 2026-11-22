package services

import (
	"fmt"
	"time"

	"klinik/internal/models"

	"gorm.io/gorm"
)

var WIB = time.FixedZone("WIB", 7*3600)

const (
	KetTersedia      = "tersedia"
	KetPenuh         = "penuh"
	KetTutup         = "tutup"
	KetTidakPraktik  = "tidak_praktik"
	PesanKuotaPenuh  = "Kuota tanggal ini penuh, silakan reservasi di waktu lain atau minggu depan."
	PesanDiLuarRentang = "Tanggal di luar rentang reservasi (maksimal 30 hari ke depan)."
	BatasReservasiHari = 30
)

type HariStatus struct {
	Tanggal string `json:"tanggal"`
	Status  string `json:"status"`
	Sisa    int    `json:"sisa"`
	Alasan  string `json:"alasan"`
}

// KetersediaanBulan: status per tanggal untuk kalender pasien (aturan §8 spec).
// Tanggal lampau dan di luar 30 hari dianggap "tutup" dengan alasan, supaya UI
// bisa menonaktifkannya tanpa pesan yang menyesatkan.
func KetersediaanBulan(db *gorm.DB, dokterID uint, bulan time.Time) ([]HariStatus, error) {
	now := time.Now().In(WIB)
	hariIni := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, WIB)
	batas := hariIni.AddDate(0, 0, BatasReservasiHari)
	awal := time.Date(bulan.Year(), bulan.Month(), 1, 0, 0, 0, 0, WIB)
	akhir := awal.AddDate(0, 1, -1)

	var jadwals []models.Jadwal
	if err := db.Where("dokter_id = ?", dokterID).Find(&jadwals).Error; err != nil {
		return nil, err
	}
	jadwalPerHari := map[int]models.Jadwal{}
	for _, j := range jadwals {
		jadwalPerHari[j.Hari] = j
	}

	var penutupan []models.Penutupan
	if err := db.Where("(dokter_id = ? OR dokter_id IS NULL) AND tanggal BETWEEN ? AND ?",
		dokterID, awal.Format("2006-01-02"), akhir.Format("2006-01-02")).Find(&penutupan).Error; err != nil {
		return nil, err
	}
	tutupPada := map[string]models.Penutupan{}
	for _, p := range penutupan {
		tutupPada[p.Tanggal] = p
	}

	var rows []struct {
		Tanggal string
		Total   int
	}
	if err := db.Table("reservasis r").
		Select("r.tanggal_kunjungan AS tanggal, COUNT(*) AS total").
		Joins("JOIN jadwals j ON j.id = r.jadwal_id").
		Where("j.dokter_id = ? AND r.status <> ? AND r.tanggal_kunjungan BETWEEN ? AND ?",
			dokterID, models.StatusBatal, awal.Format("2006-01-02"), akhir.Format("2006-01-02")).
		Group("r.tanggal_kunjungan").Scan(&rows).Error; err != nil {
		return nil, err
	}
	terpakai := map[string]int{}
	for _, r := range rows {
		terpakai[r.Tanggal] = r.Total
	}

	var out []HariStatus
	for d := awal; !d.After(akhir); d = d.AddDate(0, 0, 1) {
		tgl := d.Format("2006-01-02")
		hs := HariStatus{Tanggal: tgl}
		switch {
		case d.Before(hariIni):
			hs.Status = KetTutup
			hs.Alasan = "Tanggal sudah lewat"
		case d.After(batas):
			hs.Status = KetTutup
			hs.Alasan = PesanDiLuarRentang
		case tutupPada[tgl].Alasan != "":
			hs.Status = KetTutup
			hs.Alasan = tutupPada[tgl].Alasan
		case jadwalPerHari[int(d.Weekday())].ID == 0:
			hs.Status = KetTidakPraktik
		default:
			sisa := jadwalPerHari[int(d.Weekday())].Kuota - terpakai[tgl]
			if sisa <= 0 {
				hs.Status = KetPenuh
				hs.Alasan = PesanKuotaPenuh
			} else {
				hs.Status = KetTersedia
				hs.Sisa = sisa
			}
		}
		out = append(out, hs)
	}
	return out, nil
}

// Ambil satu jadwal praktik untuk dokter pada tanggal tertentu (dipakai booking).
func JadwalUntukTanggal(db *gorm.DB, dokterID uint, tanggal string) (*models.Jadwal, error) {
	t, err := time.Parse("2006-01-02", tanggal)
	if err != nil {
		return nil, fmt.Errorf("format tanggal tidak valid")
	}
	var j models.Jadwal
	if err := db.Where("dokter_id = ? AND hari = ?", dokterID, int(t.Weekday())).First(&j).Error; err != nil {
		return nil, err
	}
	return &j, nil
}
