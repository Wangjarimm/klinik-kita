package services

import (
	"fmt"
	"strings"
	"time"

	"klinik/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func kompakt(tanggal string) string {
	return strings.ReplaceAll(tanggal, "-", "")
}

type InputBooking struct {
	JadwalID uint
	Tanggal  string // "YYYY-MM-DD"
	Keluhan  string
}

// BookingErr dibedakan per Kind agar handler memetakan 409/422 dengan pesan spec.
type BookingErr struct {
	Kind  string // tanggal_lampau | lewat_batas | tidak_praktik | tutup | penuh | duplikat | tidak_ditemukan | format_tanggal
	Pesan string
}

func (e *BookingErr) Error() string { return e.Pesan }

type BookingService struct {
	// OnCreated dipanggil setelah commit; dipasang oleh layer WebSocket (broadcast antrean).
	OnCreated func(*models.Reservasi)
}

func (s *BookingService) Create(db *gorm.DB, pasienID uint, in InputBooking) (*models.Reservasi, error) {
	tgl, err := time.Parse("2006-01-02", in.Tanggal)
	if err != nil {
		return nil, &BookingErr{Kind: "format_tanggal", Pesan: "Format tanggal harus YYYY-MM-DD"}
	}
	tgl = time.Date(tgl.Year(), tgl.Month(), tgl.Day(), 0, 0, 0, 0, WIB)
	hariIni := time.Date(time.Now().In(WIB).Year(), time.Now().In(WIB).Month(), time.Now().In(WIB).Day(), 0, 0, 0, 0, WIB)
	switch {
	case tgl.Before(hariIni):
		return nil, &BookingErr{Kind: "tanggal_lampau", Pesan: "Tanggal sudah lewat"}
	case tgl.After(hariIni.AddDate(0, 0, BatasReservasiHari)):
		return nil, &BookingErr{Kind: "lewat_batas", Pesan: PesanDiLuarRentang}
	}

	var reservasi *models.Reservasi
	err = db.Transaction(func(tx *gorm.DB) error {
		// Kunci baris jadwal: booking paralel untuk jadwal yang sama diserialisasi di sini.
		var jadwal models.Jadwal
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&jadwal, in.JadwalID).Error; err != nil {
			return &BookingErr{Kind: "tidak_ditemukan", Pesan: "Jadwal tidak ditemukan"}
		}
		if jadwal.Hari != int(tgl.Weekday()) {
			return &BookingErr{Kind: "tidak_praktik", Pesan: "Dokter tidak praktik pada tanggal tersebut"}
		}

		var penutupan models.Penutupan
		if err := tx.Where("(dokter_id = ? OR dokter_id IS NULL) AND tanggal = ?", jadwal.DokterID, in.Tanggal).
			First(&penutupan).Error; err == nil {
			return &BookingErr{Kind: "tutup", Pesan: "Tanggal tersebut tutup: " + penutupan.Alasan}
		}

		var hitung int64
		tx.Model(&models.Reservasi{}).
			Joins("JOIN jadwals j ON j.id = reservasis.jadwal_id").
			Where("reservasis.pasien_id = ? AND j.dokter_id = ? AND reservasis.tanggal_kunjungan = ? AND reservasis.status <> ?",
				pasienID, jadwal.DokterID, in.Tanggal, models.StatusBatal).
			Count(&hitung)
		if hitung > 0 {
			return &BookingErr{Kind: "duplikat", Pesan: "Anda sudah punya reservasi aktif pada tanggal tersebut"}
		}

		tx.Model(&models.Reservasi{}).
			Joins("JOIN jadwals j ON j.id = reservasis.jadwal_id").
			Where("j.dokter_id = ? AND reservasis.tanggal_kunjungan = ? AND reservasis.status <> ?",
				jadwal.DokterID, in.Tanggal, models.StatusBatal).
			Count(&hitung)
		if int(hitung) >= jadwal.Kuota {
			return &BookingErr{Kind: "penuh", Pesan: PesanKuotaPenuh}
		}

		r := &models.Reservasi{
			PasienID:         pasienID,
			JadwalID:         jadwal.ID,
			TanggalKunjungan: in.Tanggal,
			Keluhan:          in.Keluhan,
			Status:           models.StatusMenunggu,
			NomorAntrian:     int(hitung) + 1,
		}
		// Kode memakai ID global (unik selamanya): nomor antrean per dokter bisa
		// berulang setelah batal, jadi tidak boleh jadi dasar keunikan kode.
		if err := tx.Create(r).Error; err != nil {
			return err
		}
		r.Kode = fmt.Sprintf("RSV-%s-%d", kompakt(in.Tanggal), r.ID)
		if err := tx.Model(r).Update("kode", r.Kode).Error; err != nil {
			return err
		}
		reservasi = r
		return nil
	})
	if err != nil {
		return nil, err
	}
	if s.OnCreated != nil {
		s.OnCreated(reservasi)
	}
	return reservasi, nil
}
