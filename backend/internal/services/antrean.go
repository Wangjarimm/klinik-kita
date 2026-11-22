package services

import (
	"strconv"

	"klinik/internal/models"
	"klinik/internal/ws"

	"gorm.io/gorm"
)

type AntreanService struct {
	Hub *ws.Hub
}

// Snapshot menurunkan kondisi antrean dari reservasi (tanpa tabel antrean):
// sedang melayani = nomor tertinggi berstatus dipanggil/selesai,
// sisa menunggu = jumlah berstatus menunggu/dikonfirmasi.
func (s *AntreanService) Snapshot(db *gorm.DB, dokterID uint, tanggal string) (ws.Snapshot, error) {
	snap := ws.Snapshot{DokterID: dokterID, Tanggal: tanggal}

	var m *uint
	if err := db.Model(&models.Reservasi{}).
		Joins("JOIN jadwals j ON j.id = reservasis.jadwal_id").
		Where("j.dokter_id = ? AND reservasis.tanggal_kunjungan = ? AND reservasis.status IN ?",
			dokterID, tanggal, []string{models.StatusDipanggil, models.StatusSelesai}).
		Select("MAX(reservasis.nomor_antrian)").Scan(&m).Error; err != nil {
		return snap, err
	}
	if m != nil {
		snap.SedangMelayani = *m
	}

	var sisa int64
	if err := db.Model(&models.Reservasi{}).
		Joins("JOIN jadwals j ON j.id = reservasis.jadwal_id").
		Where("j.dokter_id = ? AND reservasis.tanggal_kunjungan = ? AND reservasis.status IN ?",
			dokterID, tanggal, []string{models.StatusMenunggu, models.StatusDikonfirmasi}).
			Count(&sisa).Error; err != nil {
		return snap, err
	}
	snap.SisaMenunggu = int(sisa)
	return snap, nil
}

// PanggilBerikutnya: reservasi paling depan (menunggu/dikonfirmasi) → dipanggil, lalu broadcast.
func (s *AntreanService) PanggilBerikutnya(db *gorm.DB, dokterID uint, tanggal string) (*models.Reservasi, error) {
	var r models.Reservasi
	err := db.Joins("JOIN jadwals j ON j.id = reservasis.jadwal_id").
		Where("j.dokter_id = ? AND reservasis.tanggal_kunjungan = ? AND reservasis.status IN ?",
			dokterID, tanggal, []string{models.StatusMenunggu, models.StatusDikonfirmasi}).
		Order("reservasis.nomor_antrian ASC").First(&r).Error
	if err != nil {
		return nil, &BookingErr{Kind: "kosong", Pesan: "Tidak ada antrean yang menunggu"}
	}
	r.Status = models.StatusDipanggil
	if err := db.Save(&r).Error; err != nil {
		return nil, err
	}
	s.broadcast(db, dokterID, tanggal)
	return &r, nil
}

func (s *AntreanService) broadcast(db *gorm.DB, dokterID uint, tanggal string) {
	if s.Hub == nil {
		return
	}
	snap, err := s.Snapshot(db, dokterID, tanggal)
	if err != nil {
		return
	}
	s.Hub.Broadcast(RoomAntrean(dokterID, tanggal), snap)
}

// BroadcastSnapshot mengirim kondisi terbaru ke semua penonton room antrean.
func (s *AntreanService) BroadcastSnapshot(db *gorm.DB, dokterID uint, tanggal string) {
	s.broadcast(db, dokterID, tanggal)
}

func RoomAntrean(dokterID uint, tanggal string) string {
	return "antrean:" + strconv.FormatUint(uint64(dokterID), 10) + ":" + tanggal
}
