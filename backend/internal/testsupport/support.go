package testsupport

import (
	"testing"

	"klinik/internal/database"
	"klinik/internal/models"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupTestDB memakai klinikgo_test di MySQL compose (127.0.0.1:3307), tanpa seeder.
func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	root, err := gorm.Open(mysql.Open("root:klinik@tcp(127.0.0.1:3307)/?parseTime=true"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)
	require.NoError(t, root.Exec("CREATE DATABASE IF NOT EXISTS klinikgo_test").Error)

	db, err := gorm.Open(mysql.Open("root:klinik@tcp(127.0.0.1:3307)/klinikgo_test?parseTime=true"), &gorm.Config{
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
