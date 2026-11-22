package handlers

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"klinik/internal/database"
	"klinik/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoicePemilik200(t *testing.T) {
	r, cookie, res, _ := seedReservasiPasien(t, 3, nextWeekdayStr(t, time.Monday))

	w := doRaw(r, "GET", fmt.Sprintf("/api/v1/reservasi/invoice/%d", res.ID), nil, cookie)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	assert.Contains(t, w.Header().Get("Content-Type"), "application/pdf")
	assert.Contains(t, w.Header().Get("Content-Disposition"), "attachment")
	assert.Contains(t, w.Header().Get("Content-Disposition"), fmt.Sprintf("invoice-%s.pdf", res.Kode))
	assert.True(t, len(w.Body.Bytes()) > 4 && string(w.Body.Bytes()[:4]) == "%PDF", "body bukan PDF")
}

func TestInvoiceNonPemilik404(t *testing.T) {
	r, _, res, _ := seedReservasiPasien(t, 3, nextWeekdayStr(t, time.Monday))
	cookieLain := buatPasienLogin(t, r, "lain")

	w := doRaw(r, "GET", fmt.Sprintf("/api/v1/reservasi/invoice/%d", res.ID), nil, cookieLain)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestInvoicePetugasBoleh(t *testing.T) {
	r := newTestRouter(t)
	buatPetugas(t)
	cookieAdmin := loginAs(t, r, "adm", "petugas123")

	// reservasi dibuat lewat endpoint pasien agar alur lengkap
	cookiePasien := buatPasienLogin(t, r, "pasien_pdf")
	dokter := models.Dokter{NID: "D9", Nama: "dr. PDF", Keahlian: "Umum"}
	require.NoError(t, database.DB.Create(&dokter).Error)
	jadwal := models.Jadwal{DokterID: dokter.ID, Hari: 1, JamMulai: "08:00", JamSelesai: "12:00", Kuota: 3}
	require.NoError(t, database.DB.Create(&jadwal).Error)
	w, out := doJSON(r, "POST", "/api/v1/reservasi", map[string]any{
		"jadwal_id": jadwal.ID, "tanggal": nextWeekdayStr(t, time.Monday), "keluhan": "kontrol",
	}, cookiePasien)
	require.Equal(t, http.StatusCreated, w.Code, w.Body.String())
	resID := uint(out["data"].(map[string]any)["id"].(float64))

	w = doRaw(r, "GET", fmt.Sprintf("/api/v1/reservasi/invoice/%d", resID), nil, cookieAdmin)
	assert.Equal(t, http.StatusOK, w.Code)
}
