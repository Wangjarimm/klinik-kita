"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import toast from "react-hot-toast";
import Swal from "sweetalert2";
import { useForm } from "react-hook-form";

import { apiFetch } from "@/lib/api";
import type { Dokter, Jadwal } from "@/lib/tipe";
import { Kartu } from "@/components/Kartu";
import { Tombol } from "@/components/Tombol";
import { Input, Select } from "@/components/Input";
import { EmptyState, ErrorState, Memuat } from "@/components/EmptyState";
import { NAMA_HARI } from "@/lib/format";

// CRUD dokter + jadwal/kuota per dokter, semuanya dengan validasi dan konfirmasi hapus.
export default function AdminDokter() {
  const queryClient = useQueryClient();
  const { data, isLoading, isError, error, refetch } = useQuery({
    queryKey: ["admin-dokter"],
    queryFn: () => apiFetch<{ data: Dokter[] }>("/admin/dokter"),
  });

  const buatDokter = useForm<{
    nid: string; nama: string; keahlian: string; no_hp: string;
  }>();

  const tambah = useMutation({
    mutationFn: (v: { nid: string; nama: string; keahlian: string; no_hp: string }) =>
      apiFetch<{ id: number }>("/admin/dokter", { method: "POST", body: JSON.stringify(v) }),
    onSuccess: () => {
      toast.success("Dokter ditambahkan");
      queryClient.invalidateQueries({ queryKey: ["admin-dokter"] });
      buatDokter.reset();
    },
    onError: (e: Error) => toast.error(e.message),
  });

  return (
    <>
      <h1 className="text-2xl font-bold">Dokter & Jadwal</h1>

      <Kartu className="mt-4 p-5">
        <h2 className="text-lg font-semibold">Tambah dokter</h2>
        <form
          className="mt-3 grid gap-3 sm:grid-cols-4"
          onSubmit={buatDokter.handleSubmit(
            (v) => tambah.mutate(v),
            () => toast.error("Lengkapi NID dan nama dokter"),
          )}
        >
          <Input label="NID" placeholder="D003" {...buatDokter.register("nid", { required: true })} />
          <Input label="Nama" placeholder="dr. Nama Lengkap" {...buatDokter.register("nama", { required: true })} />
          <Input label="Keahlian" placeholder="Anak" {...buatDokter.register("keahlian")} />
          <Input label="No. HP" placeholder="0812..." {...buatDokter.register("no_hp")} />
          <div className="sm:col-span-4">
            <Tombol type="submit" disabled={tambah.isPending}>
              {tambah.isPending ? "Menyimpan..." : "Simpan Dokter"}
            </Tombol>
          </div>
        </form>
      </Kartu>

      <div className="mt-6 space-y-4">
        {isLoading && <Memuat proses="Memuat dokter..." />}
        {isError && <ErrorState pesan={(error as Error).message} onCobaLagi={() => refetch()} />}
        {!isLoading && !isError && (data?.data.length ?? 0) === 0 && (
          <EmptyState
            judul="Belum ada dokter"
            pesan="Tambahkan dokter pertama melalui formulir di atas, lalu atur jadwal praktiknya."
          />
        )}
        {data?.data.map((d) => (
          <DokterBaris key={d.id} dokter={d} />
        ))}
      </div>
    </>
  );
}

// Baris dokter + daftar jadwalnya + form tambah jadwal (kuota per hari).
function DokterBaris({ dokter }: { dokter: Dokter }) {
  const queryClient = useQueryClient();
  const { data, isLoading, isError, error, refetch } = useQuery({
    queryKey: ["jadwal-dokter", dokter.id],
    queryFn: () => apiFetch<{ data: Jadwal[] }>(`/dokter/${dokter.id}/jadwal`),
  });

  const hapusDokter = useMutation({
    mutationFn: (id: number) => apiFetch(`/admin/dokter/${id}`, { method: "DELETE" }),
    onSuccess: () => {
      toast.success("Dokter dihapus");
      queryClient.invalidateQueries({ queryKey: ["admin-dokter"] });
    },
    onError: (e: Error) => toast.error(e.message),
  });

  const konfirmasiHapusDokter = () => {
    Swal.fire({
      title: `Hapus ${dokter.nama}?`,
      text: "Jadwal praktik dokter ini ikut terhapus.",
      icon: "warning",
      showCancelButton: true,
      confirmButtonText: "Hapus",
      cancelButtonText: "Batal",
      confirmButtonColor: "#b91c1c",
    }).then((r) => {
      if (r.isConfirmed) hapusDokter.mutate(dokter.id);
    });
  };

  const form = useForm<{ hari: string; jam_mulai: string; jam_selesai: string; kuota: string }>();

  const tambahJadwal = useMutation({
    mutationFn: (v: { hari: string; jam_mulai: string; jam_selesai: string; kuota: string }) =>
      apiFetch(`/admin/dokter/${dokter.id}/jadwal`, {
        method: "POST",
        body: JSON.stringify({
          hari: Number(v.hari),
          jam_mulai: v.jam_mulai,
          jam_selesai: v.jam_selesai,
          kuota: Number(v.kuota || 3),
        }),
      }),
    onSuccess: () => {
      toast.success("Jadwal ditambahkan");
      queryClient.invalidateQueries({ queryKey: ["jadwal-dokter", dokter.id] });
      form.reset({ hari: "1", jam_mulai: "08:00", jam_selesai: "12:00", kuota: "3" });
    },
    onError: (e: Error) => toast.error(e.message),
  });

  const hapusJadwal = useMutation({
    mutationFn: (id: number) => apiFetch(`/admin/jadwal/${id}`, { method: "DELETE" }),
    onSuccess: () => {
      toast.success("Jadwal dihapus");
      queryClient.invalidateQueries({ queryKey: ["jadwal-dokter", dokter.id] });
    },
    onError: (e: Error) => toast.error(e.message),
  });

  return (
    <Kartu className="p-5">
      <div className="flex items-start justify-between gap-3">
        <div>
          <h2 className="text-sm font-semibold">{dokter.nama}</h2>
          <p className="text-xs text-[var(--color-teks-sec)]">
            NID {dokter.nid} · {dokter.keahlian} · {dokter.no_hp}
          </p>
        </div>
        <Tombol varian="bahaya" onClick={konfirmasiHapusDokter}>
          Hapus Dokter
        </Tombol>
      </div>

      {isLoading && <p className="mt-3 text-sm text-[var(--color-teks-sec)]">Memuat jadwal...</p>}
      {isError && <ErrorState pesan={(error as Error).message} onCobaLagi={() => refetch()} />}
      {data?.data.length === 0 && (
        <p className="mt-3 text-sm text-[var(--color-teks-sec)]">
          Belum ada jadwal. Tambahkan hari praktik pertama di bawah.
        </p>
      )}
      <ul className="mt-3 space-y-2">
        {data?.data.map((j) => (
          <li key={j.id} className="flex items-center justify-between rounded-[var(--radius-k)] bg-stone-50 px-3 py-2 text-sm">
            <span>
              {NAMA_HARI[j.hari]}, {j.jam_mulai}-{j.jam_selesai} · kuota {j.kuota}
            </span>
            <Tombol
              varian="ghost"
              onClick={() => {
                Swal.fire({
                  title: "Hapus jadwal ini?",
                  icon: "warning",
                  showCancelButton: true,
                  confirmButtonText: "Hapus",
                  cancelButtonText: "Batal",
                  confirmButtonColor: "#b91c1c",
                }).then((r) => {
                  if (r.isConfirmed) hapusJadwal.mutate(j.id);
                });
              }}
            >
              Hapus
            </Tombol>
          </li>
        ))}
      </ul>

      <form
        className="mt-4 grid gap-3 sm:grid-cols-5"
        onSubmit={form.handleSubmit((v) => tambahJadwal.mutate(v))}
      >
        <Select label="Hari" {...form.register("hari")}>
          {NAMA_HARI.map((h, i) => (
            <option key={i} value={i}>
              {h}
            </option>
          ))}
        </Select>
        <Input label="Jam mulai" type="time" {...form.register("jam_mulai", { required: true })} />
        <Input label="Jam selesai" type="time" {...form.register("jam_selesai", { required: true })} />
        <Input label="Kuota/hari" type="number" min={1} {...form.register("kuota")} />
        <div className="flex items-end">
          <Tombol type="submit" className="w-full" disabled={tambahJadwal.isPending}>
            Tambah Jadwal
          </Tombol>
        </div>
      </form>
    </Kartu>
  );
}
