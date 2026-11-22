package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"klinik/internal/database"
	"klinik/internal/models"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Admin mengubah status reservasi → pasien menerima pesan di WS notifikasi.
func TestNotifikasiPerUser(t *testing.T) {
	r, cookiePasien, res, _ := seedReservasiPasien(t, 3, nextWeekdayStr(t, time.Monday))
	buatPetugas(t)
	cookieAdmin := loginAs(t, r, "adm", "petugas123")

	srv := httptest.NewServer(r)
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws/notifikasi"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, http.Header{"Cookie": {
		cookiePasien.Name + "=" + cookiePasien.Value,
	}})
	require.NoError(t, err)
	defer conn.Close()

	w, _ := doJSON(r, "PATCH", fmt.Sprintf("/api/v1/admin/reservasi/%d/status", res.ID),
		map[string]string{"status": models.StatusSelesai, "catatan": "flu ringan"}, cookieAdmin)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var msg map[string]any
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	require.NoError(t, conn.ReadJSON(&msg))
	assert.Equal(t, "notifikasi", msg["tipe"])
	assert.Contains(t, msg["pesan"], res.Kode)
	assert.Contains(t, msg["pesan"], models.StatusSelesai)
}

// M-4: pasien juga menerima notifikasi saat antreannya dipanggil (bukan hanya saat status diubah).
func TestNotifikasiSaatPanggilBerikutnya(t *testing.T) {
	r, cookiePasien, res, _ := seedReservasiPasien(t, 3, nextWeekdayStr(t, time.Monday))
	buatPetugas(t)
	cookieAdmin := loginAs(t, r, "adm", "petugas123")

	srv := httptest.NewServer(r)
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws/notifikasi"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, http.Header{"Cookie": {
		cookiePasien.Name + "=" + cookiePasien.Value,
	}})
	require.NoError(t, err)
	defer conn.Close()

	var jadwal models.Jadwal
	require.NoError(t, database.DB.First(&jadwal, res.JadwalID).Error)

	w, _ := doJSON(r, "POST", "/api/v1/admin/antrean/panggil",
		map[string]any{"dokter_id": jadwal.DokterID, "tanggal": res.TanggalKunjungan}, cookieAdmin)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var msg map[string]any
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	require.NoError(t, conn.ReadJSON(&msg))
	assert.Equal(t, "notifikasi", msg["tipe"])
	assert.Contains(t, msg["pesan"], fmt.Sprintf("%02d", res.NomorAntrian))
}

func TestStatistikDanRiwayat(t *testing.T) {
	r, _, res, pasien := seedReservasiPasien(t, 3, nextWeekdayStr(t, time.Monday))
	buatPetugas(t)
	cookieAdmin := loginAs(t, r, "adm", "petugas123")

	// ubah status → selesai + catatan, lalu muncul di riwayat
	w, _ := doJSON(r, "PATCH", fmt.Sprintf("/api/v1/admin/reservasi/%d/status", res.ID),
		map[string]string{"status": models.StatusSelesai, "catatan": "flu ringan, istirahat"}, cookieAdmin)
	require.Equal(t, http.StatusOK, w.Code)

	w, out := doJSON(r, "GET", fmt.Sprintf("/api/v1/admin/pasien/%d/riwayat", pasien.ID), nil, cookieAdmin)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	riwayat := out["data"].([]any)
	require.Len(t, riwayat, 1)
	item := riwayat[0].(map[string]any)
	assert.Equal(t, "flu ringan, istirahat", item["catatan_admin"])
	assert.Equal(t, models.StatusSelesai, item["status"])

	// statistik: angka dari DB nyata (hari ini 0 reservasi kecuali hari kunjungan besok)
	w, out = doJSON(r, "GET", "/api/v1/admin/statistik", nil, cookieAdmin)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	require.NotNil(t, out["hari_ini"], "body: %s", w.Body.String())
	hariIni := out["hari_ini"].(map[string]any)
	assert.Equal(t, float64(0), hariIni[models.StatusMenunggu], "tidak ada reservasi hari ini di fixture")
	dokterTop := out["dokter_terpopuler"].([]any)
	require.Len(t, dokterTop, 1)
	assert.Equal(t, "dr. Senin", dokterTop[0].(map[string]any)["nama"])
	perHari := out["per_hari"].([]any)
	require.Len(t, perHari, 30)
}
