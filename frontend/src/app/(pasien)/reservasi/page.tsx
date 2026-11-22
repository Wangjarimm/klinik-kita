"use client";

import Link from "next/link";
import { Suspense } from "react";
import { useSearchParams } from "next/navigation";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import toast from "react-hot-toast";
import Swal from "sweetalert2";

import { API_URL, apiFetch } from "@/lib/api";
import type { Reservasi } from "@/lib/tipe";
import { Kartu } from "@/components/Kartu";
import { BadgeStatus } from "@/components/BadgeStatus";
import { Tombol } from "@/components/Tombol";
import { RelWaktu } from "@/components/RelWaktu";
import { EmptyState, ErrorState, Memuat } from "@/components/EmptyState";
import { jadwalTeks, tanggalIndo } from "@/lib/format";
import { useAntrean } from "@/lib/useAntrean";

// Daftar reservasi milik pasien: batal (modal konfirmasi), invoice, antrean live.
export default function ReservasiSaya() {
  return (
    <Suspense fallback={<Memuat proses="Memuat reservasi..." />}>
      <ReservasiKonten />
    </Suspense>
  );
}

function ReservasiKonten() {
  const params = useSearchParams();
  const suksesKode = params.get("sukses");
  const queryClient = useQueryClient();

  const { data, isLoading, isError, error, refetch } = useQuery({
    queryKey: ["reservasi-saya"],
    queryFn: () => apiFetch<{ data: Reservasi[] }>("/reservasi/saya"),
  });

  const batal = useMutation({
    mutationFn: (id: number) =>
      apiFetch<{ message: string }>(`/reservasi/${id}/batal`, { method: "PATCH" }),
    onSuccess: () => {
      toast.success("Reservasi dibatalkan");
      queryClient.invalidateQueries({ queryKey: ["reservasi-saya"] });
    },
    onError: (e: Error) => toast.error(e.message),
  });

  const konfirmasiBatal = (r: Reservasi) => {
    Swal.fire({
      title: "Batalkan reservasi ini?",
      text: `Kode ${r.kode}, kunjungan ${r.tanggal_kunjungan}. Tindakan ini tidak bisa diulang.`,
      icon: "warning",
      showCancelButton: true,
      confirmButtonText: "Ya, batalkan",
      cancelButtonText: "Jangan batalkan",
      confirmButtonColor: "#b91c1c",
    }).then((hasil) => {
      if (hasil.isConfirmed) batal.mutate(r.id);
    });
  };

  // antrean live hanya untuk reservasi aktif tanggal terdekat: satu koneksi untuk yang pertama aktif
  const aktif = data?.data.find(
    (r) => r.status === "menunggu" || r.status === "dikonfirmasi" || r.status === "dipanggil",
  );
  const antrean = useAntrean(aktif?.dokter?.id ?? 0, aktif?.tanggal_kunjungan ?? "", !!aktif);

  return (
    <>
      <h1 className="text-2xl font-bold">Reservasi Saya</h1>

      {suksesKode && (
        <div className="mt-4 rounded-[var(--radius-m)] border border-[var(--color-aksi)] bg-[var(--color-aksi-soft)] p-4 text-sm">
          Reservasi <span className="font-semibold">{suksesKode}</span> berhasil dibuat. Nomor
          antrian ada di detail di bawah dan di invoice.
        </div>
      )}

      {aktif && antrean.snap && (
        <Kartu className="mt-4 p-5">
          <h2 className="text-lg font-semibold">Antrean hari kunjungan Anda</h2>
          {antrean.terputus && (
            <p className="mt-2 text-xs text-[var(--color-amber)]">
              Koneksi antrean terputus, mencoba ulang otomatis saat halaman dimuat ulang.
            </p>
          )}
          <div className="mt-3 flex items-center gap-6">
            <div>
              <p className="text-xs text-[var(--color-teks-sec)]">Nomor Anda</p>
              <p className="text-3xl font-bold tabular-nums">
                {String(aktif.nomor_antrian).padStart(2, "0")}
              </p>
            </div>
            <div>
              <p className="text-xs text-[var(--color-teks-sec)]">Sedang melayani</p>
              <p className="text-3xl font-bold tabular-nums text-[var(--color-aksi)]">
                {antrean.snap.sedang_melayani > 0
                  ? String(antrean.snap.sedang_melayani).padStart(2, "0")
                  : "belum mulai"}
              </p>
            </div>
            <div>
              <p className="text-xs text-[var(--color-teks-sec)]">Sisa menunggu</p>
              <p className="text-3xl font-bold tabular-nums">{antrean.snap.sisa_menunggu}</p>
            </div>
          </div>
        </Kartu>
      )}

      <div className="mt-6 space-y-4">
        {isLoading && <Memuat proses="Memuat daftar reservasi..." />}
        {isError && <ErrorState pesan={(error as Error).message} onCobaLagi={() => refetch()} />}
        {!isLoading && !isError && (data?.data.length ?? 0) === 0 && (
          <EmptyState
            judul="Belum ada reservasi"
            pesan="Reservasi pertama Anda dibuat dari daftar dokter: pilih dokter, pilih tanggal, isi keluhan."
            aksi={
              <Link href="/dokter">
                <Tombol>Buat Reservasi</Tombol>
              </Link>
            }
          />
        )}
        {data?.data.map((r) => (
          <Kartu key={r.id} className="p-5">
            <div className="flex flex-wrap items-start justify-between gap-3">
              <div>
                <p className="font-semibold tabular-nums">{r.kode}</p>
                <p className="text-sm text-[var(--color-teks-sec)]">
                  {r.dokter?.nama} ({r.dokter?.keahlian}) · {jadwalTeks(r.jadwal ?? { hari: 0, jam_mulai: "", jam_selesai: "" })}
                </p>
                <p className="text-sm">{tanggalIndo(r.tanggal_kunjungan)}</p>
              </div>
              <div className="flex flex-col items-end gap-2">
                <BadgeStatus status={r.status} />
                <p className="text-sm">
                  Nomor antrian{" "}
                  <span className="font-bold tabular-nums">
                    {String(r.nomor_antrian).padStart(2, "0")}
                  </span>
                </p>
              </div>
            </div>
            {r.jadwal && (
              <div className="mt-3">
                <RelWaktu jamMulai={r.jadwal.jam_mulai} jamSelesai={r.jadwal.jam_selesai} />
              </div>
            )}
            <p className="mt-3 text-sm text-[var(--color-teks-sec)]">Keluhan: {r.keluhan}</p>
            {r.catatan_admin && (
              <p className="mt-1 text-sm">Catatan klinik: {r.catatan_admin}</p>
            )}
            <div className="mt-4 flex flex-wrap gap-2">
              <a href={`${API_URL}/api/v1/reservasi/invoice/${r.id}`} target="_blank" rel="noreferrer">
                <Tombol varian="outline">Download Invoice</Tombol>
              </a>
              {r.status === "menunggu" && (
                <Tombol varian="bahaya" onClick={() => konfirmasiBatal(r)} disabled={batal.isPending}>
                  Batalkan Reservasi
                </Tombol>
              )}
            </div>
          </Kartu>
        ))}
      </div>
    </>
  );
}
