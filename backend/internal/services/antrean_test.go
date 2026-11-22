package services

import (
	"fmt"
	"testing"
	"time"

	"klinik/internal/models"
	"klinik/internal/testsupport"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func seedAntrean(t *testing.T, n int) (*gorm.DB, models.Jadwal, string, []uint) {
	t.Helper()
	db := testsupport.SetupTestDB(t, "services")
	dokter := models.Dokter{NID: "D1", Nama: "dr. Senin", Keahlian: "Umum"}
	require.NoError(t, db.Create(&dokter).Error)
	jadwal := models.Jadwal{DokterID: dokter.ID, Hari: 1, JamMulai: "08:00", JamSelesai: "12:00", Kuota: 10}
	require.NoError(t, db.Create(&jadwal).Error)

	senin := nextWeekday(t, time.Monday, 0)
	var ids []uint
	for i := 0; i < n; i++ {
		res := models.Reservasi{
			Kode: "RSV-A" + senin + "-" + fmt.Sprint(i), PasienID: uint(i + 1), JadwalID: jadwal.ID,
			TanggalKunjungan: senin, Keluhan: "cek", Status: models.StatusMenunggu, NomorAntrian: i + 1,
		}
		require.NoError(t, db.Create(&res).Error)
		ids = append(ids, res.ID)
	}
	return db, jadwal, senin, ids
}

func TestPanggilBerikutnyaDanSnapshot(t *testing.T) {
	db, jadwal, senin, ids := seedAntrean(t, 3)
	svc := &AntreanService{}

	// panggil pertama → nomor 1 dipanggil
	r, err := svc.PanggilBerikutnya(db, jadwal.DokterID, senin)
	require.NoError(t, err)
	assert.Equal(t, ids[0], r.ID)
	assert.Equal(t, models.StatusDipanggil, r.Status)

	snap, err := svc.Snapshot(db, jadwal.DokterID, senin)
	require.NoError(t, err)
	assert.Equal(t, uint(1), snap.SedangMelayani)
	assert.Equal(t, 2, snap.SisaMenunggu)

	// selesai nomor 1 → snapshot tetap melayani 1 (nomor tertinggi dipanggil/selesai)
	require.NoError(t, db.Model(&models.Reservasi{}).Where("id = ?", ids[0]).Update("status", models.StatusSelesai).Error)
	snap, _ = svc.Snapshot(db, jadwal.DokterID, senin)
	assert.Equal(t, uint(1), snap.SedangMelayani)

	// panggil berikutnya → nomor 2
	r2, err := svc.PanggilBerikutnya(db, jadwal.DokterID, senin)
	require.NoError(t, err)
	assert.Equal(t, ids[1], r2.ID)
	snap, _ = svc.Snapshot(db, jadwal.DokterID, senin)
	assert.Equal(t, uint(2), snap.SedangMelayani)
	assert.Equal(t, 1, snap.SisaMenunggu)

	// antrean habis → error
	_, _, _, _ = db, jadwal, senin, ids
	require.NoError(t, db.Model(&models.Reservasi{}).Where("id IN ?", ids).Update("status", models.StatusSelesai).Error)
	_, err = svc.PanggilBerikutnya(db, jadwal.DokterID, senin)
	assert.Error(t, err)
}
