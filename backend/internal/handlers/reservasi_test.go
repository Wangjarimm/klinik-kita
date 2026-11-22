package handlers

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"klinik/internal/database"
	"klinik/internal/models"
	"klinik/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func hash(p string) string {
	b, _ := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	return string(b)
}

func seedReservasiPasien(t *testing.T, kuota int, tanggal string) (*gin.Engine, *http.Cookie, models.Reservasi, models.Pasien) {
	t.Helper()
	r := newTestRouter(t)
	dokter := models.Dokter{NID: "D1", Nama: "dr. Senin", Keahlian: "Umum"}
	require.NoError(t, database.DB.Create(&dokter).Error)
	jadwal := models.Jadwal{DokterID: dokter.ID, Hari: int(hariTanggal(tanggal)), JamMulai: "08:00", JamSelesai: "12:00", Kuota: kuota}
	require.NoError(t, database.DB.Create(&jadwal).Error)

	user := models.User{Username: "budi", Email: "budi@mail.com", Password: hash("rahasia123"), Role: models.RolePasien}
	require.NoError(t, database.DB.Create(&user).Error)
	pasien := models.Pasien{UserID: &user.ID, NamaLengkap: "Budi"}
	require.NoError(t, database.DB.Create(&pasien).Error)

	svc := &services.BookingService{}
	res, err := svc.Create(database.DB, pasien.ID, services.InputBooking{
		JadwalID: jadwal.ID, Tanggal: tanggal, Keluhan: "demam",
	})
	require.NoError(t, err)
	cookie := loginAs(t, r, "budi", "rahasia123")
	return r, cookie, *res, pasien
}

func hariTanggal(tanggal string) time.Weekday {
	t, _ := time.Parse("2006-01-02", tanggal)
	return t.Weekday()
}

func TestReservasiSaya(t *testing.T) {
	r, cookie, res, _ := seedReservasiPasien(t, 3, nextWeekdayStr(t, time.Monday))

	w, out := doJSON(r, "GET", "/api/v1/reservasi/saya", nil, cookie)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	list := out["data"].([]any)
	require.Len(t, list, 1)
	first := list[0].(map[string]any)
	assert.Equal(t, res.Kode, first["kode"])
	assert.Contains(t, first, "dokter")
	assert.Contains(t, first, "jadwal")
}

func TestBatalSampaiH1(t *testing.T) {
	r, cookie, res, _ := seedReservasiPasien(t, 3, nextWeekdayStr(t, time.Monday))

	w, out := doJSON(r, "PATCH", fmt.Sprintf("/api/v1/reservasi/%d/batal", res.ID), nil, cookie)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	assert.Equal(t, "batal", out["data"].(map[string]any)["status"])

	// kuota kembali: sisa = kuota penuh lagi
	var jadwal models.Jadwal
	require.NoError(t, database.DB.First(&jadwal, res.JadwalID).Error)
	bulan, _ := time.Parse("2006-01", res.TanggalKunjungan[:7])
	hari, err := services.KetersediaanBulan(database.DB, jadwal.DokterID, bulan)
	require.NoError(t, err)
	for _, h := range hari {
		if h.Tanggal == res.TanggalKunjungan {
			assert.Equal(t, "tersedia", h.Status)
			assert.Equal(t, jadwal.Kuota, h.Sisa)
		}
	}
}

func TestBatalHariHTolak(t *testing.T) {
	hariIni := time.Now().In(services.WIB).Format("2006-01-02")
	r, cookie, res, _ := seedReservasiPasien(t, 3, hariIni)

	w, _ := doJSON(r, "PATCH", fmt.Sprintf("/api/v1/reservasi/%d/batal", res.ID), nil, cookie)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestBatalBukanMiliknya404(t *testing.T) {
	r, _, res, _ := seedReservasiPasien(t, 3, nextWeekdayStr(t, time.Monday))
	cookieLain := buatPasienLogin(t, r, "pasien_lain")

	w, _ := doJSON(r, "PATCH", fmt.Sprintf("/api/v1/reservasi/%d/batal", res.ID), nil, cookieLain)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func nextWeekdayStr(t *testing.T, want time.Weekday) string {
	t.Helper()
	now := time.Now().In(services.WIB)
	add := (int(want) - int(now.Weekday()) + 7) % 7
	if add == 0 {
		add = 7
	}
	return now.AddDate(0, 0, add).Format("2006-01-02")
}
