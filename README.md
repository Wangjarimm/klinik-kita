# Klinik Kita

Sistem reservasi dokter online: pasien booking lewat kalender berkuota, petugas mengelola
jadwal dan antrian live. Proyek portfolio, dibangun ulang dari proyek lama (Go raw mux +
HTML statis) menjadi monorepo modern.

- `backend/` Go 1.22 (Gin + GORM + MySQL 8), WebSocket antrian, invoice PDF
- `frontend/` Next.js 15 (App Router, TypeScript, Tailwind CSS v4)
- `DESIGN.md` arah visual (teal + stone, ENERGY 2 / RHYTHM 2 / MOTION 1)

## Menjalankan

```bash
cp .env.example .env
docker compose up --build
```

- Frontend: http://localhost:3000
- API: http://localhost:8080/api/v1/health
- MySQL host: 127.0.0.1:3307 (database `klinikgo`, dibuat + di-seed otomatis)

## Akun awal (seeder, data contoh)

- Petugas: `petugas` / `petugas123`
- Dokter contoh: dr. Ani Wulandari (Anak), dr. Budi Santoso (Poli Umum), jadwal Senin-Jumat 08:00-12:00, kuota 3/hari
- Pasien: daftar sendiri lewat halaman Daftar

## Alur utama

1. Pasien: daftar → masuk → pilih dokter → pilih tanggal di kalender (tanggal penuh/tutup
   dikunci dengan pesan alasan) → isi keluhan → reservasi jadi, nomor antrian tercetak di
   invoice yang bisa di-download.
2. Petugas: halaman Reservasi & Antrean → Cek-in → **Panggil Berikutnya** → pasien dan
   layar display ter-update real-time (WebSocket) → Selesai + catatan → masuk riwayat pasien.
3. Layar antrian klinik (TV): `http://localhost:3000/display?dokter=1&tanggal=YYYY-MM-DD`
   (tanpa login, nomor antrean bukan data sensitif).

## Pengembangan lokal

- Backend: `cd backend && go test ./...` (butuh MySQL dari compose berjalan; race test:
  `TEST_DSN="root:klinik@tcp(host.docker.internal:3307)/" go test -race ./internal/services/`
  dari dalam container golang:1.26)
- Frontend: `cd frontend && npm install && npm run dev`

## Catatan keamanan

- JWT di httpOnly cookie, bcrypt, middleware role per endpoint, cek kepemilikan reservasi
  di server.
- `JWT_SECRET` dari `.env`, jangan pakai nilai contoh di luar pengembangan.
