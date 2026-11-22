"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import toast from "react-hot-toast";
import Swal from "sweetalert2";

import { apiFetch } from "@/lib/api";
import type { Dokter, Reservasi } from "@/lib/tipe";
import { Kartu } from "@/components/Kartu";
import { Tombol } from "@/components/Tombol";
import { BadgeStatus } from "@/components/BadgeStatus";
import { Select, Textarea } from "@/components/Input";
import { EmptyState, ErrorState, Memuat } from "@/components/EmptyState";
import { jadwalTeks, LABEL_STATUS, tanggalIndo } from "@/lib/format";
import { useState } from "react";
import { useAntrean } from "@/lib/useAntrean";

const SEMUA_STATUS = ["", "menunggu", "dikonfirmasi", "dipanggil", "selesai", "batal"];

// Manajemen reservasi hari praktik: filter, cek-in, panggil berikutnya (live), selesai, batal.
export default function AdminReservasi() {
  const queryClient = useQueryClient();
  const [tanggal, setTanggal] = useState(() => new Date(Date.now() + 7 * 3600_000).toISOString().slice(0, 10));
  const [dokterID, setDokterID] = useState("");
  const [status, setStatus] = useState("");

  const { data: dokterData } = useQuery({
    queryKey: ["admin-dokter"],
    queryFn: () => apiFetch<{ data: Dokter[] }>("/admin/dokter"),
  });

  const { data, isLoading, isError, error, refetch } = useQuery({
    queryKey: ["admin-reservasi", tanggal, dokterID, status],
    queryFn: () =>
      apiFetch<{ data: Reservasi[] }>(
        `/admin/reservasi?tanggal=${tanggal}&dokter_id=${dokterID}&status=${status}`,
      ),
  });

  const antrean = useAntrean(Number(dokterID), tanggal, dokterID !== "");

  const ubahStatus = useMutation({
    mutationFn: ({ id, status, catatan }: { id: number; status: string; catatan?: string }) =>
      apiFetch(`/admin/reservasi/${id}/status`, {
        method: "PATCH",
        body: JSON.stringify({ status, catatan: catatan ?? "" }),
      }),
    onSuccess: () => {
      toast.success("Status diperbarui");
      queryClient.invalidateQueries({ queryKey: ["admin-reservasi"] });
    },
    onError: (e: Error) => toast.error(e.message),
  });

  const panggil = useMutation({
    mutationFn: () =>
      apiFetch<{ data: Reservasi }>("/admin/antrean/panggil", {
        method: "POST",
        body: JSON.stringify({ dokter_id: Number(dokterID), tanggal }),
      }),
    onSuccess: (hasil) => {
      toast.success(`Nomor ${String(hasil.data.nomor_antrian).padStart(2, "0")} dipanggil`);
      queryClient.invalidateQueries({ queryKey: ["admin-reservasi"] });
    },
    onError: (e: Error) => toast.error(e.message),
  });

  const tandaiSelesai = (r: Reservasi) => {
    let catatan = "";
    Swal.fire({
      title: `Selesaikan nomor ${String(r.nomor_antrian).padStart(2, "0")}?`,
      input: "textarea",
      inputPlaceholder: "Catatan singkat (opsional, masuk ke riwayat pasien)",
      showCancelButton: true,
      confirmButtonText: "Tandai Selesai",
      cancelButtonText: "Batal",
    }).then((hasil) => {
      if (hasil.isConfirmed) {
        catatan = (hasil.value as string) || "";
        ubahStatus.mutate({ id: r.id, status: "selesai", catatan });
      }
    });
  };

  const batalkan = (r: Reservasi) => {
    Swal.fire({
      title: `Batalkan reservasi ${r.kode}?`,
      input: "text",
      inputPlaceholder: "Alasan pembatalan",
      showCancelButton: true,
      confirmButtonText: "Batalkan",
      cancelButtonText: "Batal",
      confirmButtonColor: "#b91c1c",
    }).then((hasil) => {
      if (hasil.isConfirmed && hasil.value) {
        ubahStatus.mutate({ id: r.id, status: "batal", catatan: hasil.value as string });
      } else if (hasil.isConfirmed) {
        toast.error("Alasan wajib diisi saat membatalkan");
      }
    });
  };

  return (
    <>
      <h1 className="text-2xl font-bold">Reservasi & Antrean</h1>

      <div className="mt-4 grid gap-3 sm:grid-cols-3">
        <Select label="Tanggal" value={tanggal} onChange={(e) => setTanggal(e.target.value)}>
          {[...Array(14)].map((_, i) => {
            const t = new Date(Date.now() + 7 * 3600_000);
            t.setDate(t.getDate() + i);
            const iso = t.toISOString().slice(0, 10);
            return (
              <option key={iso} value={iso}>
                {tanggalIndo(iso)}
              </option>
            );
          })}
        </Select>
        <Select label="Dokter" value={dokterID} onChange={(e) => setDokterID(e.target.value)}>
          <option value="">Semua dokter</option>
          {dokterData?.data.map((d) => (
            <option key={d.id} value={d.id}>
              {d.nama}
            </option>
          ))}
        </Select>
        <Select label="Status" value={status} onChange={(e) => setStatus(e.target.value)}>
          {SEMUA_STATUS.map((s) => (
            <option key={s} value={s}>
              {s === "" ? "Semua status" : LABEL_STATUS[s]}
            </option>
          ))}
        </Select>
      </div>

      {dokterID && (
        <Kartu className="mt-4 flex flex-wrap items-center justify-between gap-4 p-5">
          <div className="flex items-center gap-8">
            <div>
              <p className="text-xs text-[var(--color-teks-sec)]">Sedang melayani</p>
              <p className="text-3xl font-bold tabular-nums text-[var(--color-aksi)]">
                {antrean.snap && antrean.snap.sedang_melayani > 0
                  ? String(antrean.snap.sedang_melayani).padStart(2, "0")
                  : "belum mulai"}
              </p>
            </div>
            <div>
              <p className="text-xs text-[var(--color-teks-sec)]">Sisa menunggu</p>
              <p className="text-3xl font-bold tabular-nums">{antrean.snap?.sisa_menunggu ?? "-"}</p>
            </div>
          </div>
          <Tombol onClick={() => panggil.mutate()} disabled={panggil.isPending || !dokterID}>
            {panggil.isPending ? "Memanggil..." : "Panggil Berikutnya"}
          </Tombol>
        </Kartu>
      )}

      <div className="mt-6 space-y-3">
        {isLoading && <Memuat proses="Memuat reservasi..." />}
        {isError && <ErrorState pesan={(error as Error).message} onCobaLagi={() => refetch()} />}
        {!isLoading && !isError && (data?.data.length ?? 0) === 0 && (
          <EmptyState
            judul="Tidak ada reservasi pada filter ini"
            pesan="Ubah tanggal, dokter, atau status. Reservasi muncul begitu pasien memesan lewat kalender."
          />
        )}
        {data?.data.map((r) => (
          <Kartu key={r.id} className="flex flex-wrap items-start justify-between gap-3 p-4">
            <div>
              <p className="text-sm font-semibold tabular-nums">
                {String(r.nomor_antrian).padStart(2, "0")} · {r.kode}
              </p>
              <p className="text-sm">
                {r.dokter?.nama ?? "-"} · {jadwalTeks(r.jadwal ?? { hari: 0, jam_mulai: "", jam_selesai: "" })}
              </p>
              <p className="text-xs text-[var(--color-teks-sec)]">Keluhan: {r.keluhan}</p>
              {r.catatan_admin && (
                <p className="text-xs text-[var(--color-teks-sec)]">Catatan: {r.catatan_admin}</p>
              )}
            </div>
            <div className="flex flex-col items-end gap-2">
              <BadgeStatus status={r.status} />
              <div className="flex flex-wrap justify-end gap-1">
                {r.status === "menunggu" && (
                  <Tombol
                    varian="outline"
                    onClick={() => ubahStatus.mutate({ id: r.id, status: "dikonfirmasi" })}
                  >
                    Cek-in
                  </Tombol>
                )}
                {(r.status === "dikonfirmasi" || r.status === "menunggu") && (
                  <Tombol onClick={() => ubahStatus.mutate({ id: r.id, status: "dipanggil" })}>
                    Panggil
                  </Tombol>
                )}
                {(r.status === "dipanggil" || r.status === "dikonfirmasi") && (
                  <Tombol varian="outline" onClick={() => tandaiSelesai(r)}>
                    Selesai
                  </Tombol>
                )}
                {r.status !== "selesai" && r.status !== "batal" && (
                  <Tombol varian="bahaya" onClick={() => batalkan(r)}>
                    Batal
                  </Tombol>
                )}
              </div>
            </div>
          </Kartu>
        ))}
      </div>
    </>
  );
}
