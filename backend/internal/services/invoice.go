package services

import (
	"bytes"
	"fmt"
	"time"

	"klinik/internal/models"

	"github.com/go-pdf/fpdf"
	"gorm.io/gorm"
)

var namaHari = []string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}

var warnaTeal = struct{ R, G, B int }{15, 118, 110} // #0F766E sesuai DESIGN.md

type DataInvoice struct {
	Reservasi models.Reservasi
	Pasien    models.Pasien
	Dokter    models.Dokter
	Jadwal    models.Jadwal
}

func MuatInvoice(db *gorm.DB, reservasiID uint) (*DataInvoice, error) {
	var d DataInvoice
	if err := db.First(&d.Reservasi, reservasiID).Error; err != nil {
		return nil, err
	}
	if err := db.First(&d.Pasien, d.Reservasi.PasienID).Error; err != nil {
		return nil, err
	}
	if err := db.First(&d.Jadwal, d.Reservasi.JadwalID).Error; err != nil {
		return nil, err
	}
	if err := db.First(&d.Dokter, d.Jadwal.DokterID).Error; err != nil {
		return nil, err
	}
	return &d, nil
}

// RenderInvoice: PDF A4. shortcut: memakai core font Helvetica karena copy UI
// Indonesia murni Latin; ganti ke font TTF (DejaVu) jika butuh karakter non-Latin.
func RenderInvoice(d *DataInvoice) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 18, 20)
	pdf.SetAutoPageBreak(true, 18)
	pdf.AddPage()

	// Header wordmark
	pdf.SetFont("Helvetica", "B", 20)
	pdf.SetTextColor(warnaTeal.R, warnaTeal.G, warnaTeal.B)
	pdf.Cell(0, 12, "Klinik Kita")
	pdf.Ln(8)
	pdf.SetFont("Helvetica", "", 12)
	pdf.SetTextColor(90, 87, 82)
	pdf.Cell(0, 7, "Bukti Reservasi")
	pdf.Ln(12)

	// Kode + nomor antrian
	pdf.SetFont("Helvetica", "", 11)
	pdf.SetTextColor(28, 25, 23)
	pdf.Cell(95, 8, "Kode Reservasi: "+d.Reservasi.Kode)
	pdf.Cell(0, 8, "Status: "+d.Reservasi.Status)
	pdf.Ln(14)

	pdf.SetFillColor(warnaTeal.R, warnaTeal.G, warnaTeal.B)
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(0, 18, fmt.Sprintf("NOMOR ANTRIAN: %02d", d.Reservasi.NomorAntrian),
		"", 1, "C", true, 0, "")
	pdf.Ln(8)

	// Data
	barisan := [][2]string{
		{"Nama Pasien", d.Pasien.NamaLengkap},
		{"NIK", d.Pasien.NIK},
		{"Dokter", fmt.Sprintf("%s (Keahlian: %s)", d.Dokter.Nama, d.Dokter.Keahlian)},
		{"Jadwal", fmt.Sprintf("%s, %s-%s WIB", namaHari[d.Jadwal.Hari], d.Jadwal.JamMulai, d.Jadwal.JamSelesai)},
		{"Tanggal Kunjungan", d.Reservasi.TanggalKunjungan},
		{"Keluhan", d.Reservasi.Keluhan},
		{"Tanggal Cetak", time.Now().In(WIB).Format("02-01-2006 15:04 WIB")},
	}
	pdf.SetFont("Helvetica", "", 11)
	for _, b := range barisan {
		pdf.SetTextColor(90, 87, 82)
		pdf.CellFormat(45, 8, b[0], "", 0, "", false, 0, "")
		pdf.SetTextColor(28, 25, 23)
		pdf.MultiCell(0, 8, b[1], "", "L", false)
	}
	pdf.Ln(6)
	pdf.SetFont("Helvetica", "", 9)
	pdf.SetTextColor(90, 87, 82)
	pdf.MultiCell(0, 5, "Tunjukkan bukti ini kepada petugas klinik saat kedatangan.\nReservasi hari-H hanya bisa dibatalkan melalui petugas klinik.", "", "L", false)

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
