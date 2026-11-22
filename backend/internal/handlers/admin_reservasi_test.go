package handlers

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdminListReservasi(t *testing.T) {
	r, _, res, _ := seedReservasiPasien(t, 3, nextWeekdayStr(t, time.Monday))
	buatPetugas(t)
	cookieAdmin := loginAs(t, r, "adm", "petugas123")

	// filter tanggal kunjungan
	w, out := doJSON(r, "GET", fmt.Sprintf("/api/v1/admin/reservasi?tanggal=%s", res.TanggalKunjungan), nil, cookieAdmin)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	list := out["data"].([]any)
	require.Len(t, list, 1)
	item := list[0].(map[string]any)
	assert.Equal(t, res.Kode, item["kode"])
	assert.Contains(t, item, "dokter")
	assert.Contains(t, item, "jadwal")

	// filter tanggal lain → kosong
	w, out = doJSON(r, "GET", "/api/v1/admin/reservasi?tanggal=2030-01-01", nil, cookieAdmin)
	require.Equal(t, http.StatusOK, w.Code)
	assert.Len(t, out["data"].([]any), 0)
}
