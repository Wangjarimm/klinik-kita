"use client";

import Link from "next/link";
import { useParams } from "next/navigation";
import { useQuery } from "@tanstack/react-query";
import { useState } from "react";
import toast from "react-hot-toast";

import { apiFetch } from "@/lib/api";
import type { Dokter, HariStatus, Jadwal } from "@/lib/tipe";
import { Kartu } from "@/components/Kartu";
import { Tombol } from "@/components/Tombol";
import { RelWaktu } from "@/components/RelWaktu";
import { Kalender } from "@/components/Kalender";
import { ErrorState, Memuat } from "@/components/EmptyState";
import { jadwalTeks, NAMA_HARI } from "@/lib/format";

type Detail = { data: { dokter: Dokter; jadwal: Jadwal[] } };

// Detail dokter: jadwal praktik (motif rel waktu) + kalender ketersediaan bulanan.
export default function DetailDokter() {
  const { id } = useParams<{ id: string }>();
  const [bulan, setBulan] = useState(() => {
    const wib = new Date(Date.now() + 7 * 3600_000);
    return wib.toISOString().slice(0, 7);
  });
  const [terpilih, setTerpilih] = useState<string | undefined>();

  const { data, isLoading, isError, error, refetch } = useQuery({
    queryKey: ["dokter-detail", id],
    queryFn: () => apiFetch<Detail>(`/dokter/${id}`),
    enabled: !!id,
  });

  const pilih = (tanggal: string, status: HariStatus["status"], alasan: string) => {
    if (status === "tersedia") {
      setTerpilih(tanggal);
      return;
    }
    if (status === "penuh") {
      toast(alasan || "Kuota tanggal ini penuh, silakan reservasi di waktu lain atau minggu depan.");
      return;
    }
    if (status === "tutup") {
      toast(`Tanggal tersebut tutup: ${alasan || "ditutup klinik"}`);
      return;
    }
    toast("Dokter tidak praktik pada hari itu. Silakan pilih hari praktik lain.");
  };

  if (isLoading) return <Memuat proses="Memuat detail dokter..." />;
  if (isError)
    return <ErrorState pesan={(error as Error).message} onCobaLagi={() => refetch()} />;
  if (!data) return null;
  const { dokter, jadwal } = data.data;

  return (
    <>
      <Link href="/dokter" className="text-sm font-semibold text-[var(--color-aksi)]">
        ← Kembali ke daftar dokter
      </Link>
      <h1 className="mt-3 text-2xl font-bold">{dokter.nama}</h1>
      <p className="text-sm text-[var(--color-teks-sec)]">Keahlian: {dokter.keahlian}</p>

      <div className="mt-6 grid gap-6 lg:grid-cols-2">
        <Kartu className="p-5">
          <h2 className="text-lg font-semibold">Jadwal Praktik</h2>
          {jadwal.length === 0 ? (
            <p className="mt-3 text-sm text-[var(--color-teks-sec)]">
              Dokter ini belum punya jadwal praktik. Hubungi klinik untuk informasi terbaru.
            </p>
          ) : (
            <ul className="mt-4 space-y-4">
              {jadwal.map((j) => (
                <li key={j.id}>
                  <p className="text-sm font-semibold">
                    {NAMA_HARI[j.hari]} <span className="font-normal text-[var(--color-teks-sec)]">({jadwalTeks(j)})</span>
                  </p>
                  <div className="mt-2">
                    <RelWaktu jamMulai={j.jam_mulai} jamSelesai={j.jam_selesai} />
                  </div>
                  <p className="mt-1 text-xs text-[var(--color-teks-sec)]">
                    Kuota {j.kuota} pasien per hari praktik
                  </p>
                </li>
              ))}
            </ul>
          )}
        </Kartu>

        <Kartu className="p-5">
          <h2 className="text-lg font-semibold">Pilih Tanggal</h2>
          <div className="mt-4">
            <Kalender dokterID={Number(id)} bulan={bulan} onBulan={setBulan} pilih={pilih} terpilih={terpilih} />
          </div>
          {terpilih && (
            <div className="mt-4 rounded-[var(--radius-k)] bg-[var(--color-aksi-soft)] p-3 text-sm">
              Tanggal <span className="font-semibold">{terpilih}</span> siap dipesan.
              <Link href={`/dokter/${id}/booking?tanggal=${terpilih}`}>
                <Tombol className="mt-3 w-full">Lanjut Buat Reservasi</Tombol>
              </Link>
            </div>
          )}
        </Kartu>
      </div>
    </>
  );
}
