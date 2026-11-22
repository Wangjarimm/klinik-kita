# Klinik Kita

Sistem reservasi dokter online: pasien booking lewat kalender berkuota, staff mengelola
jadwal dan antrian live. Proyek portfolio.

- `backend/` Go (Gin + GORM + MySQL 8) + WebSocket antrian + invoice PDF
- `frontend/` Next.js (App Router, TypeScript, Tailwind CSS v4)
- `DESIGN.md` arah visual

## Menjalankan

```bash
cp .env.example .env
docker compose up --build
```

- Frontend: http://localhost:3000
- API: http://localhost:8080/api/v1/health
- MySQL host: 127.0.0.1:3307 (database `klinikgo`)

## Pengembangan lokal

- Backend: `cd backend && go test ./...` (butuh MySQL dari compose berjalan)
- Frontend: `cd frontend && npm install && npm run dev`
