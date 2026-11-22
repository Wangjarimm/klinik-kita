"use client";

import Link from "next/link";
import { Suspense, useState } from "react";
import { useQuery } from "@tanstack/react-query";

import { apiFetch } from "@/lib/api";
import type { Dokter, HariStatus } from "@/lib/tipe";
import { Kalender } from "@/components/Kalender";
import { useMe } from "@/lib/auth";

type SlotPilihan = { tanggal: string; sisa: number; alasan: string; status: HariStatus["status"] };

// Landing publik: hero + demo kalender (data real dari API), cara kerja, dokter, invoice.
export default function Landing() {
  return (
    <Suspense fallback={<main className="min-h-screen" />}>
      <LandingKonten />
    </Suspense>
  );
}

function LandingKonten() {
  const { data: me } = useMe();
  const jalur = me ? (me.role === "pasien" ? "/dokter" : "/admin") : "/login";

  const { data: dokterData } = useQuery({
    queryKey: ["dokter", ""],
    queryFn: () => apiFetch<{ data: Dokter[] }>("/dokter"),
  });
  const dokter = dokterData?.data ?? [];
  const [pilihDokter, setPilihDokter] = useState(0);
  const dokterAktif = dokter[pilihDokter];

  const wib = new Date(Date.now() + 7 * 3600_000);
  const [bulan, setBulan] = useState(`${wib.getFullYear()}-${String(wib.getMonth() + 1).padStart(2, "0")}`);
  const [slot, setSlot] = useState<SlotPilihan | null>(null);

  const pilihTanggal = (tanggal: string, status: HariStatus["status"], alasan: string, sisa: number) => {
    setSlot({ tanggal, sisa, alasan, status });
  };

  const ctapasien = dokterAktif ? `/dokter/${dokterAktif.id}` : "/login";

  return (
    <>
      <header className="sticky top-0 z-10 border-b border-[var(--color-garis)] bg-[var(--color-bg)]">
        <div className="mx-auto flex min-h-[68px] max-w-5xl flex-wrap items-center gap-4 px-5">
          <Link href="/" className="mr-auto flex items-center gap-3 text-xl font-bold tracking-tight">
            <svg width="34" height="34" viewBox="0 0 40 40" aria-hidden="true">
              <rect width="40" height="40" rx="10" fill="var(--color-aksi)" />
              <path d="M20 13v14M13 20h14" stroke="#fff" strokeWidth="2.6" strokeLinecap="round" />
            </svg>
            Klinik Kita
          </Link>
          <a href="#cara" className="hidden text-sm font-medium text-[var(--color-teks-sec)] hover:text-[var(--color-teks)] md:block">
            Cara kerja
          </a>
          <a href="#dokter" className="hidden text-sm font-medium text-[var(--color-teks-sec)] hover:text-[var(--color-teks)] md:block">
            Dokter
          </a>
          <Link href={me ? jalur : "/login"}>
            <button className="btn-outline">Masuk</button>
          </Link>
          {!me && (
            <Link href="/register">
              <button className="btn-primary">Daftar</button>
            </Link>
          )}
        </div>
      </header>

      <main>
        {/* Hero + demo kalender */}
        <div className="mx-auto grid max-w-5xl gap-10 px-5 py-14 md:grid-cols-[1.05fr_.95fr] md:items-center md:py-20">
          <div>
            <h1 className="max-w-[15ch] text-4xl font-extrabold leading-[1.08] tracking-tight sm:text-5xl">
              Pesan jadwal dokter, dapatkan nomor antrian Anda langsung.
            </h1>
            <p className="mt-5 max-w-[46ch] text-lg text-[var(--color-teks-sec)]">
              Pilih dokter, pilih tanggal di kalender, dan dapatkan nomor antrian Anda. Sisa slot
              setiap hari terbatas dan ditampilkan apa adanya.
            </p>
            <div className="mt-8 flex flex-wrap items-center gap-3">
              <Link href={jalur}>
                <button className="btn-primary min-h-[52px] px-7 text-base">Buat reservasi</button>
              </Link>
              <a href="#cara">
                <button className="btn-outline min-h-[52px] px-7 text-base">Lihat cara kerjanya</button>
              </a>
            </div>
            <p className="mt-4 text-sm text-[var(--color-teks-sec)]">
              Sudah punya akun? Masuk dulu supaya riwayat reservasi Anda tersimpan.
            </p>
          </div>

          {/* Demo kalender: sisa slot REAL dari database */}
          <div className="rounded-[var(--radius-l)] border border-[var(--color-garis)] bg-[var(--color-surface)] p-5 shadow-[0_18px_40px_-22px_rgba(15,60,55,.35)]">
            <div className="mb-3 flex items-baseline justify-between gap-2">
              <strong>Coba pilih tanggal</strong>
              <span className="rounded-md bg-[var(--color-amber-soft)] px-2.5 py-0.5 text-xs font-semibold text-[var(--color-amber)]">
                Sisa slot asli dari database
              </span>
            </div>

            {dokter.length > 0 ? (
              <>
                <div className="mb-4 grid grid-cols-2 gap-2" role="group" aria-label="Pilih dokter">
                  {dokter.slice(0, 2).map((d, i) => (
                    <button
                      key={d.id}
                      onClick={() => { setPilihDokter(i); setSlot(null); }}
                      aria-pressed={pilihDokter === i}
                      className={`flex min-h-[56px] items-center gap-2.5 rounded-[var(--radius-k)] border px-3 py-2 text-left ${
                        pilihDokter === i
                          ? "border-[var(--color-aksi)] bg-[var(--color-aksi-soft)] outline outline-1 outline-[var(--color-aksi)]"
                          : "border-[var(--color-garis)] bg-[var(--color-bg)]"
                      }`}
                    >
                      <span
                        aria-hidden
                        className={`grid h-9 w-9 flex-none place-items-center rounded-full text-xs font-extrabold ${
                          pilihDokter === i ? "bg-[var(--color-aksi)] text-white" : "bg-[var(--color-aksi-soft)] text-[var(--color-aksi)]"
                        }`}
                      >
                        {d.nama.replace("dr. ", "").slice(0, 2).toUpperCase()}
                      </span>
                      <span>
                        <b className="block text-sm leading-tight">{d.nama}</b>
                        <small className="block text-xs text-[var(--color-teks-sec)]">{d.keahlian}</small>
                      </span>
                    </button>
                  ))}
                </div>

                <Kalender dokterID={dokterAktif.id} bulan={bulan} onBulan={(b) => { setBulan(b); setSlot(null); }} pilih={pilihTanggal} />

                <div aria-live="polite" className="mt-3">
                  {!slot && (
                    <div className="rounded-[var(--radius-m)] border border-[var(--color-garis)] bg-[var(--color-bg)] px-4 py-4 text-sm text-[var(--color-teks-sec)]">
                      <b className="block text-[var(--color-teks)]">Belum ada tanggal dipilih</b>
                      Pilih satu tanggal di kalender untuk melihat sisa slotnya.
                    </div>
                  )}
                  {slot && slot.status === "tersedia" && (
                    <div>
                      <div className="grid grid-cols-[1fr_auto] rounded-[var(--radius-m)] border border-[var(--color-garis)] bg-[var(--color-bg)]">
                        <div className="p-4">
                          <small className="block text-xs text-[var(--color-teks-sec)]">{slot.tanggal}</small>
                          <b>{dokterAktif.nama}</b>
                          <small className="block text-xs text-[var(--color-teks-sec)]">
                            Sisa slot: {slot.sisa}
                          </small>
                        </div>
                        <div className="grid min-w-[112px] place-items-center border-l-2 border-dashed border-[var(--color-garis)] p-4 text-center">
                          <div>
                            <small className="block text-xs text-[var(--color-teks-sec)]">Nomor antrian</small>
                            <div className="text-2xl font-extrabold tracking-tight text-[var(--color-amber)]">
                              ditentukan saat booking
                            </div>
                          </div>
                        </div>
                      </div>
                      <Link href={ctapasien} className="mt-3 block">
                        <button className="btn-primary w-full">Lanjut buat reservasi untuk tanggal ini</button>
                      </Link>
                    </div>
                  )}
                  {slot && slot.status === "penuh" && (
                    <div className="rounded-[var(--radius-m)] border border-[var(--color-garis)] bg-[var(--color-bg)] px-4 py-4 text-sm text-[var(--color-teks-sec)]">
                      <b className="block text-[var(--color-teks)]">Slot {dokterAktif.nama} penuh pada {slot.tanggal}</b>
                      {slot.alasan || "Pilih tanggal lain, atau ganti dokter."}
                    </div>
                  )}
                  {slot && (slot.status === "tutup" || slot.status === "tidak_praktik") && (
                    <div className="rounded-[var(--radius-m)] border border-[var(--color-garis)] bg-[var(--color-bg)] px-4 py-4 text-sm text-[var(--color-teks-sec)]">
                      <b className="block text-[var(--color-teks)]">Tanggal {slot.tanggal} tidak bisa dipesan</b>
                      {slot.alasan || "Dokter tidak praktik pada hari itu."}
                    </div>
                  )}
                </div>
              </>
            ) : (
              <p className="py-6 text-sm text-[var(--color-teks-sec)]">
                Daftar dokter sedang kosong. Sisa slot muncul di sini setelah petugas menambahkan
                dokter dan jadwal.
              </p>
            )}
            <p className="mt-3 text-xs text-[var(--color-teks-sec)]">
              Nomor antrian final ditentukan saat reservasi dibuat dan tercetak di invoice.
            </p>
          </div>
        </div>

        {/* Cara kerja: alur asli 3 langkah */}
        <section id="cara" className="border-t border-[var(--color-garis)] py-16">
          <div className="mx-auto grid max-w-5xl gap-8 px-5 md:grid-cols-[.8fr_1.2fr]">
            <div>
              <h2 className="text-3xl font-extrabold leading-tight tracking-tight">
                Dari rumah ke ruang periksa, tanpa antre menebak-nebak
              </h2>
              <p className="mt-3 max-w-[36ch] text-[var(--color-teks-sec)]">
                Reservasi dilakukan online. Selebihnya cukup datang langsung.
              </p>
            </div>
            <ol className="m-0 list-none border-t border-[var(--color-garis)] p-0">
              {[
                ["Reservasi online", "Pilih dokter dan tanggal di kalender. Slot yang tersisa untuk hari itu ditampilkan apa adanya, termasuk saat sudah penuh."],
                ["Datang langsung", "Nomor antrian sudah Anda pegang sejak reservasi, jadi Anda tinggal datang ke klinik."],
                ["Tunjukkan invoice", "Perlihatkan invoice reservasi Anda di klinik sebagai bukti."],
              ].map(([judul, isi], i) => (
                <li key={i} className="grid grid-cols-[auto_1fr] gap-5 border-b border-[var(--color-garis)] py-6">
                  <span aria-hidden className="w-[2ch] text-2xl font-extrabold text-[var(--color-aksi)]">{i + 1}</span>
                  <div>
                    <h3 className="text-lg font-bold">{judul}</h3>
                    <p className="mt-1 max-w-[52ch] text-[var(--color-teks-sec)]">{isi}</p>
                  </div>
                </li>
              ))}
            </ol>
          </div>
        </section>

        {/* Dokter real dari database */}
        <section id="dokter" className="border-t border-[var(--color-garis)] bg-[var(--color-aksi-soft)] py-16">
          <div className="mx-auto max-w-5xl px-5">
            <h2 className="text-3xl font-extrabold tracking-tight">Dokter kami</h2>
            <div className="mt-8 grid gap-5 sm:grid-cols-2">
              {dokter.map((d) => (
                <article
                  key={d.id}
                  className={`flex flex-col gap-4 rounded-[var(--radius-l)] border border-[var(--color-garis)] bg-[var(--color-surface)] p-7 ${
                    d.id === dokterAktif?.id ? "sm:p-9" : ""
                  }`}
                >
                  <span
                    aria-hidden
                    className="grid h-16 w-16 place-items-center rounded-full bg-[var(--color-aksi-soft)] text-xl font-extrabold text-[var(--color-aksi)]"
                  >
                    {d.nama.replace("dr. ", "").slice(0, 2).toUpperCase()}
                  </span>
                  <div>
                    <h3 className="text-xl font-extrabold">{d.nama}</h3>
                    <p className="text-[var(--color-teks-sec)]">{d.keahlian}</p>
                  </div>
                  <Link href={`/dokter/${d.id}`} className="mt-auto">
                    <button className={d.id === dokterAktif?.id ? "btn-primary" : "btn-outline"}>
                      Reservasi dengan {d.nama.replace("dr. ", "")}
                    </button>
                  </Link>
                </article>
              ))}
            </div>
          </div>
        </section>

        {/* Invoice */}
        <section className="border-t border-[var(--color-garis)] py-16">
          <div className="mx-auto grid max-w-5xl items-center gap-6 px-5 md:grid-cols-[1fr_auto]">
            <div>
              <h2 className="text-3xl font-extrabold tracking-tight">Reservasi selesai, invoice jadi tiket Anda</h2>
              <p className="mt-3 max-w-[54ch] text-[var(--color-teks-sec)]">
                Setelah reservasi, Anda mendapat invoice dan nomor antrian. Bawa keduanya saat datang
                ke klinik.
              </p>
            </div>
            <Link href={jalur}>
              <button className="btn-primary min-h-[52px] whitespace-nowrap px-7 text-base">
                Pilih dokter dan tanggal
              </button>
            </Link>
          </div>
        </section>
      </main>

      <footer className="border-t border-[var(--color-garis)] py-8 text-sm text-[var(--color-teks-sec)]">
        <div className="mx-auto flex max-w-5xl flex-wrap justify-between gap-4 px-5">
          <span>Klinik Kita. Reservasi online, datang langsung, tunjukkan invoice.</span>
          <a href="#cara" className="hover:text-[var(--color-teks)]">Kembali ke atas</a>
        </div>
      </footer>
    </>
  );
}
