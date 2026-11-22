"use client";

import Link from "next/link";
import { useQuery } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { Dokter } from "@/lib/tipe";
import { Kartu } from "@/components/Kartu";
import { Tombol } from "@/components/Tombol";
import { EmptyState, ErrorState, Memuat } from "@/components/EmptyState";
import { useMe } from "@/lib/auth";

// Landing publik: apa yang dilakukan sistem, dokter yang tersedia, satu CTA jelas.
export default function Landing() {
  const { data: me } = useMe();
  const { data, isLoading, isError, error, refetch } = useQuery({
    queryKey: ["dokter", ""],
    queryFn: () => apiFetch<{ data: Dokter[] }>("/dokter"),
  });

  const jalurCTA = me ? (me.role === "pasien" ? "/dokter" : "/admin") : "/login";

  return (
    <>
      <header>
        <nav className="mx-auto flex max-w-5xl items-center justify-between px-4 py-4">
          <Link href="/" className="text-lg font-bold text-[var(--color-teks)]">
            Klinik Kita
          </Link>
          {me ? (
            <Tombol varian="outline" onClick={() => (window.location.href = jalurCTA)}>
              Buka dashboard
            </Tombol>
          ) : (
            <div className="flex gap-2">
              <Link href="/login">
                <Tombol varian="outline">Masuk</Tombol>
              </Link>
              <Link href="/register">
                <Tombol>Daftar</Tombol>
              </Link>
            </div>
          )}
        </nav>
      </header>

      <main className="mx-auto max-w-5xl px-4 pb-16">
        <section className="py-14">
          <h1 className="max-w-xl text-3xl font-bold leading-tight">
            Reservasi dokter, nomor antrian Anda sudah menunggu di tangan.
          </h1>
          <p className="mt-3 max-w-xl text-[var(--color-teks-sec)]">
            Pilih dokter, pilih tanggal di kalender, dan dapatkan nomor antrian Anda. Sisa slot
            setiap hari terbatas dan ditampilkan apa adanya.
          </p>
          <div className="mt-6">
            <Link href={jalurCTA}>
              <Tombol className="px-6">Buat Reservasi</Tombol>
            </Link>
          </div>
        </section>

        <section aria-labelledby="judul-dokter" className="mt-4">
          <h2 id="judul-dokter" className="mb-4 text-xl font-semibold">
            Dokter kami
          </h2>
          {isLoading && <Memuat proses="Memuat daftar dokter..." />}
          {isError && <ErrorState pesan={(error as Error).message} onCobaLagi={() => refetch()} />}
          {!isLoading && !isError && (data?.data.length ?? 0) === 0 && (
            <EmptyState
              judul="Belum ada dokter terdaftar"
              pesan="Data dokter diisi oleh petugas klinik. Cek kembali setelah jadwal diperbarui."
            />
          )}
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {data?.data.map((d) => (
              <Kartu key={d.id} className="p-5">
                <div className="flex items-center gap-3">
                  <span
                    aria-hidden
                    className="flex h-11 w-11 items-center justify-center rounded-full bg-[var(--color-aksi-soft)] text-sm font-bold text-[var(--color-aksi)]"
                  >
                    {d.nama.replace("dr. ", "").slice(0, 2).toUpperCase()}
                  </span>
                  <div>
                    <h3 className="text-sm font-semibold">{d.nama}</h3>
                    <p className="text-xs text-[var(--color-teks-sec)]">{d.keahlian}</p>
                  </div>
                </div>
              </Kartu>
            ))}
          </div>
        </section>
      </main>

      <footer className="border-t border-[var(--color-garis)] py-6 text-center text-xs text-[var(--color-teks-sec)]">
        Klinik Kita. Reservasi online, datang langsung, tunjukkan invoice.
      </footer>
    </>
  );
}
