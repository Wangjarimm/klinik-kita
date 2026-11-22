"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import toast from "react-hot-toast";
import Swal from "sweetalert2";
import { useForm } from "react-hook-form";

import { apiFetch } from "@/lib/api";
import type { Dokter, Penutupan } from "@/lib/tipe";
import { Kartu } from "@/components/Kartu";
import { Tombol } from "@/components/Tombol";
import { Input, Select } from "@/components/Input";
import { EmptyState, ErrorState, Memuat } from "@/components/EmptyState";
import { tanggalIndo } from "@/lib/format";

// Tutup jadwal: per dokter atau seluruh klinik. Tanggal yang ditutup terkunci di kalender pasien.
export default function AdminPenutupan() {
  const queryClient = useQueryClient();
  const { data, isLoading, isError, error, refetch } = useQuery({
    queryKey: ["penutupan"],
    queryFn: () => apiFetch<{ data: Penutupan[] }>("/admin/penutupan"),
  });
  const { data: dokterData } = useQuery({
    queryKey: ["admin-dokter"],
    queryFn: () => apiFetch<{ data: Dokter[] }>("/admin/dokter"),
  });

  const form = useForm<{ dokter_id: string; tanggal: string; alasan: string }>({
    defaultValues: { dokter_id: "", tanggal: "", alasan: "" },
  });

  const tambah = useMutation({
    mutationFn: (v: { dokter_id: string; tanggal: string; alasan: string }) =>
      apiFetch<{ id: number }>("/admin/penutupan", {
        method: "POST",
        body: JSON.stringify({
          dokter_id: v.dokter_id ? Number(v.dokter_id) : null,
          tanggal: v.tanggal,
          alasan: v.alasan,
        }),
      }),
    onSuccess: () => {
      toast.success("Tanggal ditutup");
      queryClient.invalidateQueries({ queryKey: ["penutupan"] });
      form.reset({ dokter_id: "", tanggal: "", alasan: "" });
    },
    onError: (e: Error) => toast.error(e.message),
  });

  const hapus = useMutation({
    mutationFn: (id: number) => apiFetch(`/admin/penutupan/${id}`, { method: "DELETE" }),
    onSuccess: () => {
      toast.success("Penutupan dihapus");
      queryClient.invalidateQueries({ queryKey: ["penutupan"] });
    },
    onError: (e: Error) => toast.error(e.message),
  });

  return (
    <>
      <h1 className="text-2xl font-bold">Tutup Jadwal</h1>
      <p className="mt-1 text-sm text-[var(--color-teks-sec)]">
        Tanggal yang ditutup langsung terkunci di kalender pasien dengan pesan alasannya.
      </p>

      <Kartu className="mt-4 p-5">
        <h2 className="text-lg font-semibold">Tutup tanggal</h2>
        <form
          className="mt-3 grid gap-3 sm:grid-cols-4"
          onSubmit={form.handleSubmit(
            (v) => tambah.mutate(v),
            () => toast.error("Tanggal dan alasan wajib diisi"),
          )}
        >
          <Select label="Dokter" {...form.register("dokter_id")}>
            <option value="">Seluruh klinik</option>
            {dokterData?.data.map((d) => (
              <option key={d.id} value={d.id}>
                {d.nama}
              </option>
            ))}
          </Select>
          <Input label="Tanggal" type="date" {...form.register("tanggal", { required: true })} />
          <Input label="Alasan" placeholder="Dokter cuti" {...form.register("alasan", { required: true })} />
          <div className="flex items-end">
            <Tombol type="submit" className="w-full" disabled={tambah.isPending}>
              Tutup Tanggal
            </Tombol>
          </div>
        </form>
      </Kartu>

      <div className="mt-6 space-y-3">
        {isLoading && <Memuat proses="Memuat daftar penutupan..." />}
        {isError && <ErrorState pesan={(error as Error).message} onCobaLagi={() => refetch()} />}
        {!isLoading && !isError && (data?.data.length ?? 0) === 0 && (
          <EmptyState
            judul="Belum ada penutupan"
            pesan="Semua tanggal buka. Tutup tanggal ketika dokter cuti atau klinik libur."
          />
        )}
        {data?.data.map((p) => (
          <Kartu key={p.id} className="flex items-center justify-between p-4">
            <div>
              <p className="text-sm font-semibold">{tanggalIndo(p.tanggal)}</p>
              <p className="text-xs text-[var(--color-teks-sec)]">
                {p.dokter_id ? `Dokter ID ${p.dokter_id}` : "Seluruh klinik"} · {p.alasan}
              </p>
            </div>
            <Tombol
              varian="ghost"
              onClick={() => {
                Swal.fire({
                  title: "Buka kembali tanggal ini?",
                  icon: "question",
                  showCancelButton: true,
                  confirmButtonText: "Ya, buka",
                  cancelButtonText: "Batal",
                }).then((r) => {
                  if (r.isConfirmed) hapus.mutate(p.id);
                });
              }}
            >
              Buka kembali
            </Tombol>
          </Kartu>
        ))}
      </div>
    </>
  );
}
