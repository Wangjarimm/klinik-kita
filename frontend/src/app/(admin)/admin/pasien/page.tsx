"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import toast from "react-hot-toast";
import Swal from "sweetalert2";
import { useState } from "react";

import { apiFetch } from "@/lib/api";
import type { Pasien, Reservasi } from "@/lib/tipe";
import { Kartu } from "@/components/Kartu";
import { Tombol } from "@/components/Tombol";
import { BadgeStatus } from "@/components/BadgeStatus";
import { Input, Select } from "@/components/Input";
import { EmptyState, ErrorState, Memuat } from "@/components/EmptyState";

type PasienRow = Pasien & { user_id: number | null };

// Kelola pasien (termasuk walk-in) + lihat riwayat kunjungan versi ringan.
export default function AdminPasien() {
  const queryClient = useQueryClient();
  const [q, setQ] = useState("");
  const { data, isLoading, isError, error, refetch } = useQuery({
    queryKey: ["admin-pasien", q],
    queryFn: () => apiFetch<{ data: PasienRow[] }>(`/admin/pasien?q=${encodeURIComponent(q)}`),
  });

  const tambah = useMutation({
    mutationFn: (v: Omit<Pasien, "id" | "user_id">) =>
      apiFetch<{ id: number }>("/admin/pasien", { method: "POST", body: JSON.stringify(v) }),
    onSuccess: () => {
      toast.success("Pasien ditambahkan");
      queryClient.invalidateQueries({ queryKey: ["admin-pasien"] });
    },
    onError: (e: Error) => toast.error(e.message),
  });

  return (
    <>
      <h1 className="text-2xl font-bold">Pasien & Riwayat</h1>

      <Kartu className="mt-4 p-5">
        <h2 className="text-lg font-semibold">Tambah pasien walk-in</h2>
        <form
          className="mt-3 grid gap-3 sm:grid-cols-4"
          onSubmit={(e) => {
            e.preventDefault();
            const f = e.currentTarget.elements as typeof e.currentTarget.elements & {
              nama_lengkap: HTMLInputElement;
              nik: HTMLInputElement;
              jenis_kelamin: HTMLSelectElement;
              no_hp: HTMLInputElement;
            };
            if (f.nik.value.length !== 16) {
              toast.error("NIK harus 16 angka");
              return;
            }
            tambah.mutate({
              nama_lengkap: f.nama_lengkap.value,
              nik: f.nik.value,
              jenis_kelamin: f.jenis_kelamin.value,
              no_hp: f.no_hp.value,
              tempat_lahir: "",
              tanggal_lahir: "",
              alamat: "",
            });
            e.currentTarget.reset();
          }}
        >
          <Input label="Nama lengkap" name="nama_lengkap" required />
          <Input label="NIK (16 angka)" name="nik" inputMode="numeric" pattern="\d{16}" required />
          <Select label="Jenis kelamin" name="jenis_kelamin">
            <option value="L">Laki-laki</option>
            <option value="P">Perempuan</option>
          </Select>
          <Input label="No. HP" name="no_hp" required />
          <div className="sm:col-span-4">
            <Tombol type="submit" disabled={tambah.isPending}>
              {tambah.isPending ? "Menyimpan..." : "Simpan Pasien"}
            </Tombol>
          </div>
        </form>
      </Kartu>

      <div className="mt-4 flex items-center gap-2 rounded-[var(--radius-k)] border border-[var(--color-garis)] bg-white px-3">
        <input
          type="search"
          aria-label="Cari pasien"
          placeholder="Cari nama atau NIK"
          value={q}
          onChange={(e) => setQ(e.target.value)}
          className="min-h-[44px] w-full bg-transparent text-sm outline-none"
        />
      </div>

      <div className="mt-4 space-y-3">
        {isLoading && <Memuat proses="Memuat pasien..." />}
        {isError && <ErrorState pesan={(error as Error).message} onCobaLagi={() => refetch()} />}
        {!isLoading && !isError && (data?.data.length ?? 0) === 0 && (
          <EmptyState
            judul="Belum ada pasien yang cocok"
            pesan="Pasien online mendaftar sendiri lewat halaman daftar; pasien walk-in diinput di atas."
          />
        )}
        {data?.data.map((p) => (
          <BarisPasien key={p.id} pasien={p} />
        ))}
      </div>
    </>
  );
}

function BarisPasien({ pasien }: { pasien: PasienRow }) {
  const [bukaRiwayat, setBukaRiwayat] = useState(false);
  const { data, isLoading, isError, error, refetch } = useQuery({
    queryKey: ["riwayat", pasien.id],
    queryFn: () => apiFetch<{ data: Reservasi[] }>(`/admin/pasien/${pasien.id}/riwayat`),
    enabled: bukaRiwayat,
  });

  return (
    <Kartu className="p-4">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <p className="text-sm font-semibold">{pasien.nama_lengkap}</p>
          <p className="text-xs tabular-nums text-[var(--color-teks-sec)]">
            NIK {pasien.nik} · {pasien.jenis_kelamin === "L" ? "Laki-laki" : "Perempuan"} · {pasien.no_hp}
            {pasien.user_id ? " · akun online" : " · walk-in"}
          </p>
        </div>
        <Tombol varian="outline" onClick={() => setBukaRiwayat((v) => !v)}>
          {bukaRiwayat ? "Tutup riwayat" : "Lihat riwayat kunjungan"}
        </Tombol>
      </div>

      {bukaRiwayat && (
        <div className="mt-4 border-t border-[var(--color-garis)] pt-4">
          {isLoading && <p className="text-sm text-[var(--color-teks-sec)]">Memuat riwayat...</p>}
          {isError && <ErrorState pesan={(error as Error).message} onCobaLagi={() => refetch()} />}
          {!isLoading && !isError && (data?.data.length ?? 0) === 0 && (
            <p className="text-sm text-[var(--color-teks-sec)]">
              Pasien ini belum pernah berkunjung.
            </p>
          )}
          <ul className="space-y-2">
            {data?.data.map((r) => (
              <li key={r.id} className="rounded-[var(--radius-k)] bg-stone-50 p-3 text-sm">
                <div className="flex items-center justify-between gap-2">
                  <span className="font-semibold tabular-nums">
                    {r.kode} · {r.tanggal_kunjungan}
                  </span>
                  <BadgeStatus status={r.status} />
                </div>
                <p className="mt-1 text-[var(--color-teks-sec)]">Keluhan: {r.keluhan}</p>
                {r.catatan_admin && <p>Catatan klinik: {r.catatan_admin}</p>}
              </li>
            ))}
          </ul>
        </div>
      )}
    </Kartu>
  );
}
