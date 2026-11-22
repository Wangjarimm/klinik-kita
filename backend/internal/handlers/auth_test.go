package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"klinik/internal/database"
	"klinik/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// setupTestDB memakai database klinikgo_test di MySQL compose (127.0.0.1:3307),
// tanpa seeder: setiap test mengisi datanya sendiri.
func setupTestDB(t *testing.T) {
	t.Helper()
	root, err := gorm.Open(mysql.Open("root:klinik@tcp(127.0.0.1:3307)/?parseTime=true"), &gorm.Config{
		Logger: logger.Default(logger.LogLevel(1)),
	})
	require.NoError(t, err)
	require.NoError(t, root.Exec("CREATE DATABASE IF NOT EXISTS klinikgo_test").Error)

	db, err := gorm.Open(mysql.Open("root:klinik@tcp(127.0.0.1:3307)/klinikgo_test?parseTime=true"), &gorm.Config{
		Logger: logger.Default(logger.LogLevel(1)),
	})
	require.NoError(t, err)
	require.NoError(t, db.Migrator().DropTable(
		&models.Reservasi{}, &models.Penutupan{}, &models.Jadwal{},
		&models.Pasien{}, &models.Dokter{}, &models.User{},
	).Error)
	require.NoError(t, db.AutoMigrate(
		&models.User{}, &models.Pasien{}, &models.Dokter{},
		&models.Jadwal{}, &models.Penutupan{}, &models.Reservasi{},
	).Error)
	database.DB = db
}

func newTestRouter(t *testing.T) *gin.Engine {
	t.Helper()
	setupTestDB(t)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	RegisterRoutes(r, "test-secret")
	return r
}

func doJSON(r *gin.Engine, method, path string, body any, cookies ...*http.Cookie) (*httptest.ResponseRecorder, map[string]any) {
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	for _, c := range cookies {
		req.AddCookie(c)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var out map[string]any
	if w.Body.Len() > 0 {
		_ = json.Unmarshal(w.Body.Bytes(), &out)
	}
	return w, out
}

func loginAs(t *testing.T, r *gin.Engine, username, password string) *http.Cookie {
	t.Helper()
	w, _ := doJSON(r, "POST", "/api/v1/auth/login", map[string]string{"username": username, "password": password})
	if w.Code != http.StatusOK {
		t.Fatalf("login %s gagal: %d %s", username, w.Code, w.Body.String())
	}
	for _, c := range w.Result().Cookies() {
		if c.Name == "session" {
			return c
		}
	}
	t.Fatal("cookie session tidak ditemukan")
	return nil
}

func TestRegisterCreatesUserAndPasien(t *testing.T) {
	r := newTestRouter(t)
	w, out := doJSON(r, "POST", "/api/v1/auth/register", map[string]string{
		"username": "budi", "email": "budi@mail.com", "password": "rahasia123",
	})
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	assert.Equal(t, "pasien", out["role"])

	var user models.User
	require.NoError(t, database.DB.Where("username = ?", "budi").First(&user).Error)
	assert.Equal(t, "pasien", user.Role)
	var pasien models.Pasien
	require.NoError(t, database.DB.Where("user_id = ?", user.ID).First(&pasien).Error)
}

func TestRegisterDuplicateUsername409(t *testing.T) {
	r := newTestRouter(t)
	body := map[string]string{"username": "dupe", "email": "dupe@mail.com", "password": "rahasia123"}
	doJSON(r, "POST", "/api/v1/auth/register", body)
	w, _ := doJSON(r, "POST", "/api/v1/auth/register", map[string]string{"username": "dupe", "email": "dupe2@mail.com", "password": "rahasia123"})
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestLoginWrongPassword401(t *testing.T) {
	r := newTestRouter(t)
	doJSON(r, "POST", "/api/v1/auth/register", map[string]string{"username": "sari", "email": "sari@mail.com", "password": "benar123"})
	w, _ := doJSON(r, "POST", "/api/v1/auth/login", map[string]string{"username": "sari", "password": "salah"})
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestMeWithoutCookie401(t *testing.T) {
	r := newTestRouter(t)
	w, _ := doJSON(r, "GET", "/api/v1/auth/me", nil)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestMeWithCookieReturnsRole(t *testing.T) {
	r := newTestRouter(t)
	doJSON(r, "POST", "/api/v1/auth/register", map[string]string{"username": "budi", "email": "budi@mail.com", "password": "rahasia123"})
	cookie := loginAs(t, r, "budi", "rahasia123")
	w, out := doJSON(r, "GET", "/api/v1/auth/me", nil, cookie)
	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	assert.Equal(t, "budi", out["username"])
	assert.Equal(t, "pasien", out["role"])
}

func TestRoleRequiredRejectsWrongRole(t *testing.T) {
	r := newTestRouter(t)
	doJSON(r, "POST", "/api/v1/auth/register", map[string]string{"username": "budi", "email": "budi@mail.com", "password": "rahasia123"})
	cookie := loginAs(t, r, "budi", "rahasia123")
	w, _ := doJSON(r, "GET", "/api/v1/admin/dokter", nil, cookie)
	assert.Equal(t, http.StatusForbidden, w.Code)
}
