package services

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"klinik/internal/models"
	"klinik/internal/testsupport"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func seedBookingFixture(t *testing.T, kuota int) (*gorm.DB, models.Jadwal, string) {
	t.Helper()
	db := testsupport.SetupTestDB(t, "services")
	dokter := models.Dokter{NID: "D1", Nama: "dr. Senin", Keahlian: "Umum"}
	require.NoError(t, db.Create(&dokter).Error)
	jadwal := models.Jadwal{DokterID: dokter.ID, Hari: 1, JamMulai: "08:00", JamSelesai: "12:00", Kuota: kuota}
	require.NoError(t, db.Create(&jadwal).Error)
	return db, jadwal, nextWeekday(t, time.Monday, 0)
}

func cobaBooking(t *testing.T, db *gorm.DB, pasienID uint, jadwalID uint, tanggal string) error {
	t.Helper()
	svc := &BookingService{}
	_, err := svc.Create(db, pasienID, InputBooking{JadwalID: jadwalID, Tanggal: tanggal, Keluhan: "demam 2 hari"})
	return err
}

func kindOf(t *testing.T, err error) string {
	t.Helper()
	be, ok := err.(*BookingErr)
	require.True(t, ok, "bukan BookingErr: %v", err)
	return be.Kind
}

// Review Focus 1: balapan kuota — tepat satu yang menang dari 8 booking paralel.
func TestBookingConcurrencyKuotaSatu(t *testing.T) {
	db, jadwal, senin := seedBookingFixture(t, 1)

	const N = 8
	results := make([]error, N)
	var wg sync.WaitGroup
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			results[i] = cobaBooking(t, db, uint(i+1), jadwal.ID, senin)
		}(i)
	}
	wg.Wait()

	sukses, penuh := 0, 0
	for _, err := range results {
		if err == nil {
			sukses++
		} else if kindOf(t, err) == "penuh" {
			penuh++
		}
	}
	assert.Equal(t, 1, sukses, "harus tepat satu booking sukses")
	assert.Equal(t, N-1, penuh, "sisanya harus ditolak kuota penuh")

	var r models.Reservasi
	require.NoError(t, db.Where("jadwal_id = ?", jadwal.ID).First(&r).Error)
	assert.Equal(t, 1, r.NomorAntrian)
}

func TestBookingValidasiTanggal(t *testing.T) {
	db, jadwal, senin := seedBookingFixture(t, 3)

	// tanggal lampau
	lampau := time.Now().In(WIB).AddDate(0, 0, -1).Format("2006-01-02")
	assert.Equal(t, "tanggal_lampau", kindOf(t, cobaBooking(t, db, 1, jadwal.ID, lampau)))

	// di luar 30 hari
	jauh := time.Now().In(WIB).AddDate(0, 0, 40).Format("2006-01-02")
	assert.Equal(t, "lewat_batas", kindOf(t, cobaBooking(t, db, 1, jadwal.ID, jauh)))

	// hari tidak praktik (Selasa)
	selasa := nextWeekday(t, time.Tuesday, 0)
	assert.Equal(t, "tidak_praktik", kindOf(t, cobaBooking(t, db, 1, jadwal.ID, selasa)))

	// tutup oleh admin
	require.NoError(t, db.Create(&models.Penutupan{
		DokterID: &jadwal.DokterID, Tanggal: senin, Alasan: "Dokter cuti",
	}).Error)
	err := cobaBooking(t, db, 1, jadwal.ID, senin)
	assert.Equal(t, "tutup", kindOf(t, err))
	assert.Contains(t, err.(*BookingErr).Pesan, "Dokter cuti")
}

func TestBookingDuplikatDanNomorAntrian(t *testing.T) {
	db, jadwal, senin := seedBookingFixture(t, 3)

	require.NoError(t, cobaBooking(t, db, 1, jadwal.ID, senin))
	// pasien yang sama, dokter+tanggal sama, masih aktif → duplikat
	assert.Equal(t, "duplikat", kindOf(t, cobaBooking(t, db, 1, jadwal.ID, senin)))

	// pasien lain → nomor antrian 2, kode format RSV-tanggal-002
	require.NoError(t, cobaBooking(t, db, 2, jadwal.ID, senin))
	var r models.Reservasi
	require.NoError(t, db.Where("pasien_id = ? AND jadwal_id = ?", 2, jadwal.ID).First(&r).Error)
	assert.Equal(t, 2, r.NomorAntrian)
	assert.Equal(t, fmt.Sprintf("RSV-%s-002", senin), r.Kode)
}
