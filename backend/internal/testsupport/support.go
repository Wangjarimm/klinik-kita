package testsupport

import (
	"os"
	"testing"

	"klinik/internal/database"
	"klinik/internal/models"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestRootDSN mengembalikan DSN root tanpa nama database (untuk CREATE DATABASE).
func TestRootDSN() string {
	if v := os.Getenv("TEST_DSN"); v != "" {
		return v
	}
	return "root:klinik@tcp(127.0.0.1:3307)/"
}

// SetupTestDB memakai klinikgo_test di MySQL compose, tanpa seeder.
func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	root, err := gorm.Open(mysql.Open(TestRootDSN()+"?parseTime=true"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)
	require.NoError(t, root.Exec("CREATE DATABASE IF NOT EXISTS klinikgo_test").Error)

	db, err := gorm.Open(mysql.Open(TestRootDSN()+"klinikgo_test?parseTime=true"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)
	require.NoError(t, db.Migrator().DropTable(
		&models.Reservasi{}, &models.Penutupan{}, &models.Jadwal{},
		&models.Pasien{}, &models.Dokter{}, &models.User{},
	))
	require.NoError(t, db.AutoMigrate(
		&models.User{}, &models.Pasien{}, &models.Dokter{},
		&models.Jadwal{}, &models.Penutupan{}, &models.Reservasi{},
	))
	database.DB = db
	return db
}
