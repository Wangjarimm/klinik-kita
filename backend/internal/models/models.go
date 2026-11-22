package models

import "time"

// Data contoh di seeder (dokter, jadwal) diberi label jelas sebagai contoh, bukan data nyata.
// Tanggal disimpan sebagai string "YYYY-MM-DD" dan jam "HH:MM" supaya bebas zona waktu.

type User struct {
	ID       uint   `gorm:"primaryKey" json:"user_id"`
	Username string `gorm:"size:64;uniqueIndex" json:"username"`
	Email    string `gorm:"size:128;uniqueIndex" json:"email"`
	Password string `json:"-"`
	Role     string `gorm:"size:16;index" json:"role"` // pasien | petugas
}

type Pasien struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	UserID       *uint  `gorm:"index" json:"user_id"`
	NamaLengkap  string `gorm:"size:128" json:"nama_lengkap"`
	NIK          string `gorm:"column:nik;size:16;index" json:"nik"`
	JenisKelamin string `gorm:"size:16" json:"jenis_kelamin"`
	TempatLahir  string `gorm:"size:64" json:"tempat_lahir"`
	TanggalLahir string `gorm:"size:10" json:"tanggal_lahir"`
	Alamat       string `gorm:"size:255" json:"alamat"`
	NoHP         string `gorm:"size:20" json:"no_hp"`
}

type Dokter struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	NID      string `gorm:"column:nid;size:16;uniqueIndex" json:"nid"`
	Nama     string `gorm:"size:128" json:"nama"`
	Keahlian string `gorm:"size:128" json:"keahlian"`
	NoHP     string `gorm:"size:20" json:"no_hp"`
}

type Jadwal struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	DokterID   uint   `gorm:"index" json:"dokter_id"`
	Hari       int    `json:"hari"` // 0=Minggu .. 6=Sabtu
	JamMulai   string `gorm:"size:5" json:"jam_mulai"`
	JamSelesai string `gorm:"size:5" json:"jam_selesai"`
	Kuota      int    `gorm:"default:3" json:"kuota"`
}

type Penutupan struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	DokterID *uint  `gorm:"index" json:"dokter_id"` // null = seluruh klinik
	Tanggal  string `gorm:"size:10;index" json:"tanggal"`
	Alasan   string `gorm:"size:255" json:"alasan"`
}

type Reservasi struct {
	ID              uint   `gorm:"primaryKey" json:"id"`
	Kode            string `gorm:"size:24;uniqueIndex" json:"kode"`
	PasienID        uint   `gorm:"index" json:"pasien_id"`
	JadwalID        uint   `gorm:"index" json:"jadwal_id"`
	TanggalKunjungan string `gorm:"size:10;index" json:"tanggal_kunjungan"`
	Keluhan         string `gorm:"type:text" json:"keluhan"`
	Status          string `gorm:"size:16;index" json:"status"`
	NomorAntrian    int    `json:"nomor_antrian"`
	CatatanAdmin    string    `gorm:"type:text" json:"catatan_admin"`
	CreatedAt       time.Time `json:"created_at"`
}

const (
	StatusMenunggu     = "menunggu"
	StatusDikonfirmasi = "dikonfirmasi"
	StatusDipanggil    = "dipanggil"
	StatusSelesai      = "selesai"
	StatusBatal        = "batal"
)
