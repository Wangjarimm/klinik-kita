"use client";

import Link from "next/link";
import { useQuery } from "@tanstack/react-query";
import { Search } from "lucide-react";
import { useState } from "react";

import { apiFetch } from "@/lib/api";
import type { Dokter } from "@/lib/tipe";
import { Kartu } from "@/components/Kartu";
import { EmptyState, ErrorState, Memuat } from "@/components/EmptyState";
import { Tombol } from "@/components/Tombol";

// Katalog dokter + pencarian nama/keahlian (param q ke API).
export default function DaftarDokter() {
  const [q, setQ] = useState("");
  const { data, isLoading, isError, error, refetch } = useQuery({
    queryKey: ["dokter", q],
    queryFn: () => apiFetch<{ data: Dokter[] }>(`/dokter?q=${encodeURIComponent(q)}`),
  });

  return (
    <>
      <h1 className="text-2xl font-bold">Pilih Dokter</h1>
      <div className="mt-4 flex items-center gap-2 rounded-[var(--radius-k)] border border-[var(--color-garis)] bg-white px-3">
        <Search size={18} aria-hidden className="text-[var(--color-teks-sec)]" />
        <input
          type="search"
          aria-label="Cari dokter"
          placeholder="Cari nama atau keahlian, mis. Anak"
          value={q}
          onChange={(e) => setQ(e.target.value)}
          className="min-h-[44px] w-full bg-transparent text-sm outline-none"
        />
      </div>

      <div className="mt-6">
        {isLoading && <Memuat proses="Memuat daftar dokter..." />}
        {isError && <ErrorState pesan={(error as Error).message} onCobaLagi={() => refetch()} />}
        {!isLoading && !isError && (data?.data.length ?? 0) === 0 && (
          <EmptyState
            judul={q ? "Tidak ada dokter yang cocok" : "Belum ada dokter terdaftar"}
            pesan={
              q
                ? `Pencarian "${q}" tidak menemukan dokter. Coba kata lain, mis. nama atau keahlian.`
                : "Petugas klinik belum menambahkan dokter. Cek kembali nanti."
            }
          />
        )}
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {data?.data.map((d) => (
            <Kartu key={d.id} className="flex flex-col gap-4 p-5">
              <div className="flex items-center gap-3">
                <span
                  aria-hidden
                  className="flex h-11 w-11 items-center justify-center rounded-full bg-[var(--color-aksi-soft)] text-sm font-bold text-[var(--color-aksi)]"
                >
                  {d.nama.replace("dr. ", "").slice(0, 2).toUpperCase()}
                </span>
                <div>
                  <h2 className="text-sm font-semibold">{d.nama}</h2>
                  <p className="text-xs text-[var(--color-teks-sec)]">{d.keahlian}</p>
                </div>
              </div>
              <Link href={`/dokter/${d.id}`}>
                <Tombol varian="outline" className="w-full">
                  Lihat Jadwal
                </Tombol>
              </Link>
            </Kartu>
          ))}
        </div>
      </div>
    </>
  );
}
