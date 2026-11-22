package handlers

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"klinik/internal/database"
	"klinik/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func buatPetugas(t *testing.T) {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte("petugas123"), bcrypt.DefaultCost)
	require.NoError(t, err)
	require.NoError(t, database.DB.Create(&models.User{
		Username: "adm", Email: "adm@klinik.local", Password: string(hash), Role: models.RolePetugas,
	}).Error)
}

func buatPasienLogin(t *testing.T, r *gin.Engine, username string) *http.Cookie {
	t.Helper()
	hash, _ := bcrypt.GenerateFromPassword([]byte("rahasia123"), bcrypt.DefaultCost)
	user := models.User{Username: username, Email: username + "@mail.com", Password: string(hash), Role: models.RolePasien}
	require.NoError(t, database.DB.Create(&user).Error)
	require.NoError(t, database.DB.Create(&models.Pasien{UserID: &user.ID, NamaLengkap: username}).Error)
	return loginAs(t, r, username, "rahasia123")
}

func TestAdminCRUDDokterDanJadwal(t *testing.T) {
	r := newTestRouter(t)
	buatPetugas(t)
	cookie := loginAs(t, r, "adm", "petugas123")

	w, out := doJSON(r, "POST", "/api/v1/admin/dokter", map[string]any{
		"nid": "D100", "nama": "dr. Cinta Lestari", "keahlian": "Gigi", "no_hp": "081100000001",
	}, cookie)
	require.Equal(t, http.StatusCreated, w.Code, w.Body.String())
	dokterID := int(out["id"].(float64))

	for _, h := range []int{1, 2, 3} {
		w, _ = doJSON(r, "POST", fmt.Sprintf("/api/v1/admin/dokter/%d/jadwal", dokterID),
			map[string]any{"hari": h, "jam_mulai": "08:00", "jam_selesai": "12:00", "kuota": 3}, cookie)
		require.Equal(t, http.StatusCreated, w.Code, w.Body.String())
	}

	// publik (tanpa login): daftar dokter + pencarian keahlian + jadwalnya
	w, out = doJSON(r, "GET", "/api/v1/dokter?q=gigi", nil)
	require.Equal(t, http.StatusOK, w.Code)
	require.Len(t, out["data"].([]any), 1)

	w, out = doJSON(r, "GET", fmt.Sprintf("/api/v1/dokter/%d/jadwal", dokterID), nil)
	require.Equal(t, http.StatusOK, w.Code)
	require.Len(t, out["data"].([]any), 3)

	// pasien tidak boleh CRUD admin
	cookiePasien := buatPasienLogin(t, r, "pasien1")
	w, _ = doJSON(r, "POST", "/api/v1/admin/dokter", map[string]any{"nid": "X", "nama": "x"}, cookiePasien)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestAdminPenutupanCRUD(t *testing.T) {
	r := newTestRouter(t)
	buatPetugas(t)
	cookie := loginAs(t, r, "adm", "petugas123")

	w, out := doJSON(r, "POST", "/api/v1/admin/penutupan", map[string]any{
		"tanggal": "2026-12-25", "alasan": "Libur nasional",
	}, cookie) // tanpa dokter_id = seluruh klinik
	require.Equal(t, http.StatusCreated, w.Code, w.Body.String())
	id := int(out["id"].(float64))

	w, out = doJSON(r, "GET", "/api/v1/admin/penutupan", nil, cookie)
	require.Len(t, out["data"].([]any), 1)

	w, _ = doJSON(r, "DELETE", fmt.Sprintf("/api/v1/admin/penutupan/%d", id), nil, cookie)
	require.Equal(t, http.StatusOK, w.Code)
	w, out = doJSON(r, "GET", "/api/v1/admin/penutupan", nil, cookie)
	require.Len(t, out["data"].([]any), 0)
}

// I-2: dokter/pasien dengan reservasi aktif tidak boleh dihapus (meninggalkan reservasi yatim).
func TestDeleteDitolakSaatAdaReservasiAktif(t *testing.T) {
	r, _, res, pasien := seedReservasiPasien(t, 3, nextWeekdayStr(t, time.Monday))
	buatPetugas(t)
	cookieAdmin := loginAs(t, r, "adm", "petugas123")

	var jadwal models.Jadwal
	require.NoError(t, database.DB.First(&jadwal, res.JadwalID).Error)

	w, out := doJSON(r, "DELETE", fmt.Sprintf("/api/v1/admin/dokter/%d", jadwal.DokterID), nil, cookieAdmin)
	require.Equal(t, http.StatusConflict, w.Code, w.Body.String())
	assert.Contains(t, out["error"], "reservasi aktif")

	w, out = doJSON(r, "DELETE", fmt.Sprintf("/api/v1/admin/pasien/%d", pasien.ID), nil, cookieAdmin)
	require.Equal(t, http.StatusConflict, w.Code, w.Body.String())
	assert.Contains(t, out["error"], "reservasi aktif")

	// reservasi dibatalkan → hapus lolos
	require.NoError(t, database.DB.Model(&models.Reservasi{}).Where("id = ?", res.ID).
		Update("status", models.StatusBatal).Error)
	w, _ = doJSON(r, "DELETE", fmt.Sprintf("/api/v1/admin/pasien/%d", pasien.ID), nil, cookieAdmin)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAdminPasienWalkIn(t *testing.T) {
	r := newTestRouter(t)
	buatPetugas(t)
	cookie := loginAs(t, r, "adm", "petugas123")

	w, out := doJSON(r, "POST", "/api/v1/admin/pasien", map[string]any{
		"nama_lengkap": "Slamet Walk-in", "nik": "3201234567890001", "jenis_kelamin": "L",
		"tempat_lahir": "Bandung", "tanggal_lahir": "1990-05-01", "alamat": "Jl. Mawar 1", "no_hp": "0812",
	}, cookie)
	require.Equal(t, http.StatusCreated, w.Code, w.Body.String())
	id := int(out["id"].(float64))

	var pasien models.Pasien
	require.NoError(t, database.DB.First(&pasien, id).Error)
	assert.Nil(t, pasien.UserID, "pasien walk-in tidak punya akun")

	w, out = doJSON(r, "GET", "/api/v1/admin/pasien", nil, cookie)
	require.Len(t, out["data"].([]any), 1)
}
