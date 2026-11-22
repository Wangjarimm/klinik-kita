# DESIGN.md — Arah Visual "Klinik Kita"

Produk: sistem reservasi dokter online. Audiens: pasien (mobile-first) dan admin klinik (desktop).
Tema: light only (keputusan pemilik, 2026-11-22). Data demo boleh, data dummy tidak.

## Kepribadian

"Tenang & tepercaya". Sistem kesehatan, bukan panggung: whitespace lega, hierarki jelas,
gerakan minim. Dial: **ENERGY 2 / RHYTHM 2 / MOTION 1**.

## Palet (2 inti + 1 aksen, maksimal, R-29)

- **Core 1 — Teal** `#0F766E`: warna aksi (tombol utama, link aktif, ring fokus).
  Alasan: asosiasi medis yang tidak kaku, kontras AA di atas netral hangat.
- **Core 2 — Netral hangat (stone)**: bg `#FAFAF9`, surface `#FFFFFF`, teks `#1C1917` /
  sekunder `#57534E`. Alasan: hangat sedikit supaya tidak steril.
- **Aksen — Amber** `#B45309`: hanya untuk status `menunggu` dan perhatian. Alasan: satu
  aksen di momen status nyata, bukan dekorasi.
- Status reservasi (penanda keadaan nyata, R-31): `menunggu` amber, `dikonfirmasi` teal
  outline, `dipanggil` teal solid, `selesai` stone muted, `batal` red `#B91C1C`.
- Gradien dekoratif: tidak dipakai. Bayangan: hanya elemen mengambang (dropdown, modal).
- Radius: skala 8 / 12 / 16 px; tidak ada pill-everything (R-11).

## Tipografi

**Plus Jakarta Sans**, satu keluarga untuk semua UI. Alasan: didesain foundry Indonesia
(Tokotype), karakter ramah-tenang, lisensi terbuka. Skala: display 28/36 semibold, judul
20/28 semibold, body 15/24, label 13/500. Kode reservasi & nomor antrian: `tabular-nums`.
Tanpa monospace-estetik, tanpa uppercase wide-tracking (R-06).

## Motif identitas

**"Rel waktu praktik"**: garis horizontal berlabel jam praktik (jam_mulai → jam_selesai)
yang digambar konsisten di kartu jadwal dokter, detail reservasi, invoice, dan layar display.
Alasan: memvisualkan data jadwal nyata, satu gestur yang membuat tampilan terasa milik
produk ini.

## Ikon

Lucide, dipilih dengan alasan: bobot garis tipisnya cocok dengan kepribadian tenang dan
glyph medisnya lengkap (calendar, clock, stethoscope, users, printer, bell). Hanya glyph
yang relevan konten; dilarang sparkle/star/robot/orb (R-04).

## Motion (MOTION 1)

Transisi hover/focus/state ≤ 200ms. Tanpa loop tanpa pemicu (R-19). Semua animasi
menghormati `prefers-reduced-motion: reduce`.

## Copy (bahasa Indonesia)

Spesifik dan jujur. CTA bernama aksi nyata: "Buat Reservasi", "Panggil Berikutnya",
"Download Invoice", "Batalkan Reservasi" (R-15). Tanpa em dash (R-02), tanpa buzzword
(R-16). Angka hanya dari database (R-17). Setiap layar data punya state **loading / empty /
error** yang menyebut penyebab dan tindakan berikutnya (R-27).

## Aksesibilitas

Kontras WCAG AA untuk semua pasangan warna (R-25); indikator fokus keyboard terlihat
(ring teal 2px, R-32); target sentuh ≥ 44px (R-03); elemen native dulu, custom hanya kalau
perlu, semua kontrol punya nama aksesibel (R-39).

## Implementasi frontend

Tailwind CSS v4, token di atas sebagai CSS variables; komponen kecil buatan sendiri
(Tombol, Kartu, Badge status, Input, Kalender). Tanpa library komponen besar yang memaksa
tampilan; ikon via `lucide-react`; kalender dibangun sendiri (bulan berjalan + status per
tanggal), bukan library kalender.
