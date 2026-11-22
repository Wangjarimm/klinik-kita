package services

import (
	"fmt"
	"testing"
	"time"

	"klinik/internal/models"
	"klinik/internal/testsupport"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// nextWeekday mengembalikan tanggal "YYYY-MM-DD" untuk weekday berikutnya (minggu depan).
func nextWeekday(t *testing.T, want time.Weekday, offsetWeek int) string {
	t.Helper()
	now := time.Now().In(time.FixedZone("WIB", 7*3600))
	add := (int(want) - int(now.Weekday()) + 7) % 7
	if add == 0 {
		add = 7
	}
	return now.AddDate(0, 0, add+offsetWeek*7).Format("2006-01-02")
}

func TestKetersediaanBulan(t *testing.T) {
	db := testsupport.SetupTestDB(t, "services")

	dokter := models.Dokter{NID: "D1", Nama: "dr. Senin", Keahlian: "Umum"}
	require.NoError(t, db.Create(&dokter).Error)
	// praktik hanya hari Senin, kuota 2
	jadwal := models.Jadwal{DokterID: dokter.ID, Hari: 1, JamMulai: "08:00", JamSelesai: "12:00", Kuota: 2}
	require.NoError(t, db.Create(&jadwal).Error)

	senin := nextWeekday(t, time.Monday, 0)

	// 2 reservasi aktif di hari Senin itu → penuh
	for i := 0; i < 2; i++ {
		require.NoError(t, db.Create(&models.Reservasi{
			Kode: "RSV-T" + senin[:10] + "-" + fmt.Sprint(i), PasienID: uint(i + 1), JadwalID: jadwal.ID,
			TanggalKunjungan: senin, Keluhan: "tes", Status: models.StatusMenunggu, NomorAntrian: i + 1,
		}).Error)
	}

	bulan, err := time.Parse("2006-01", senin[:7])
	require.NoError(t, err)
	hari, err := KetersediaanBulan(db, dokter.ID, bulan)
	require.NoError(t, err)

	byDate := map[string]HariStatus{}
	for _, h := range hari {
		byDate[h.Tanggal] = h
	}

	assert.Equal(t, "penuh", byDate[senin].Status)
	assert.Equal(t, 0, byDate[senin].Sisa)
	assert.Contains(t, byDate[senin].Alasan, "penuh")

	// hari Selasa: tidak praktik
	selasa := nextWeekday(t, time.Tuesday, 0)
	assert.Equal(t, "tidak_praktik", byDate[selasa].Status)

	// tanggal tutup oleh admin (dokter ini): tutup + alasan
	libur := nextWeekday(t, time.Monday, 1) // Senin minggu berikutnya
	require.NoError(t, db.Create(&models.Penutupan{
		DokterID: &dokter.ID, Tanggal: libur, Alasan: "Dokter cuti",
	}).Error)
	hari, err = KetersediaanBulan(db, dokter.ID, bulan)
	require.NoError(t, err)
	byDate = map[string]HariStatus{}
	for _, h := range hari {
		byDate[h.Tanggal] = h
	}
	assert.Equal(t, "tutup", byDate[libur].Status)
	assert.Equal(t, "Dokter cuti", byDate[libur].Alasan)

	// reservasi batal tidak mengurangi kuota
	require.NoError(t, db.Model(&models.Reservasi{}).Where("jadwal_id = ?", jadwal.ID).
		Update("status", models.StatusBatal).Error)
	hari, _ = KetersediaanBulan(db, dokter.ID, bulan)
	for _, h := range hari {
		if h.Tanggal == senin {
			assert.Equal(t, "tersedia", h.Status)
			assert.Equal(t, 2, h.Sisa)
		}
	}
}
