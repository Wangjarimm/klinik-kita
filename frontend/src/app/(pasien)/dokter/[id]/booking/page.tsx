"use client";

import Link from "next/link";
import { Suspense } from "react";
import { useParams, useRouter, useSearchParams } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useMutation, useQuery } from "@tanstack/react-query";
import toast from "react-hot-toast";

import { apiFetch } from "@/lib/api";
import type { Dokter, Jadwal, Reservasi } from "@/lib/tipe";
import { Textarea } from "@/components/Input";
import { Tombol } from "@/components/Tombol";
import { Kartu } from "@/components/Kartu";
import { ErrorState, Memuat } from "@/components/EmptyState";
import { tanggalIndo } from "@/lib/format";

const schema = z.object({
  keluhan: z.string().min(3, "Ceritakan singkat keluhan Anda").max(500, "Maksimal 500 karakter"),
});
type FormNilai = z.infer<typeof schema>;

// Ringkasan pilihan + form keluhan. Sukses: kode, nomor antrian, invoice.
export default function Booking() {
  return (
    <Suspense fallback={<Memuat proses="Memuat formulir reservasi..." />}>
      <BookingKonten />
    </Suspense>
  );
}

function BookingKonten() {
  const { id } = useParams<{ id: string }>();
  const params = useSearchParams();
  const router = useRouter();
  const tanggal = params.get("tanggal") ?? "";

  const { data: detail, isLoading, isError, error, refetch } = useQuery({
    queryKey: ["dokter-detail", id],
    queryFn: () => apiFetch<{ data: { dokter: Dokter; jadwal: Jadwal[] } }>(`/dokter/${id}`),
    enabled: !!id && !!tanggal,
  });

  // jadwal praktik yang cocok dengan hari kunjungan
  const [y, m, d] = tanggal.split("-").map(Number);
  const weekday = tanggal ? new Date(y, m - 1, d).getDay() : -1;
  const jadwalCocok = detail?.data.jadwal.find((j) => j.hari === weekday);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<FormNilai>({ resolver: zodResolver(schema) });

  const mutasi = useMutation({
    mutationFn: (nilai: FormNilai) =>
      apiFetch<{ data: Reservasi }>("/reservasi", {
        method: "POST",
        body: JSON.stringify({ jadwal_id: jadwalCocok!.id, tanggal, keluhan: nilai.keluhan }),
      }),
    onSuccess: (hasil) => {
      toast.success("Reservasi berhasil dibuat");
      router.push(`/reservasi?sukses=${hasil.data.kode}`);
    },
    onError: (e: Error) => toast.error(e.message),
  });

  if (!tanggal) {
    return (
      <Kartu className="mx-auto max-w-md p-6 text-center">
        <h1 className="text-lg font-semibold">Tanggal belum dipilih</h1>
        <p className="mt-2 text-sm text-[var(--color-teks-sec)]">
          Pilih tanggal praktik pada kalender dokter terlebih dahulu.
        </p>
        <Link href={`/dokter/${id}`}>
          <Tombol className="mt-4">Pilih Tanggal</Tombol>
        </Link>
      </Kartu>
    );
  }

  if (isLoading) return <Memuat proses="Memuat jadwal dokter..." />;
  if (isError) return <ErrorState pesan={(error as Error).message} onCobaLagi={() => refetch()} />;

  if (!jadwalCocok) {
    return (
      <Kartu className="mx-auto max-w-md p-6 text-center">
        <h1 className="text-lg font-semibold">Dokter tidak praktik pada tanggal tersebut</h1>
        <p className="mt-2 text-sm text-[var(--color-teks-sec)]">
          Pilih ulang tanggal di kalender pada hari praktik {detail?.data.dokter.nama}.
        </p>
        <Link href={`/dokter/${id}`}>
          <Tombol className="mt-4">Pilih Tanggal</Tombol>
        </Link>
      </Kartu>
    );
  }

  return (
    <div className="mx-auto max-w-md">
      <h1 className="text-2xl font-bold">Konfirmasi Reservasi</h1>
      <Kartu className="mt-4 p-5">
        <dl className="space-y-2 text-sm">
          <div className="flex justify-between">
            <dt className="text-[var(--color-teks-sec)]">Tanggal kunjungan</dt>
            <dd className="font-semibold">{tanggalIndo(tanggal)}</dd>
          </div>
          <div className="flex justify-between">
            <dt className="text-[var(--color-teks-sec)]">Dokter</dt>
            <dd className="font-semibold">
              {detail?.data.dokter.nama} ({detail?.data.dokter.keahlian})
            </dd>
          </div>
          <div className="flex justify-between">
            <dt className="text-[var(--color-teks-sec)]">Jam praktik</dt>
            <dd className="font-semibold">
              {jadwalCocok.jam_mulai}-{jadwalCocok.jam_selesai}
            </dd>
          </div>
        </dl>
        <form className="mt-5 space-y-4" onSubmit={handleSubmit((v) => mutasi.mutate(v))}>
          <Textarea
            label="Keluhan"
            rows={4}
            placeholder="Contoh: demam 2 hari, batuk dan pilek"
            {...register("keluhan")}
            error={errors.keluhan?.message}
          />
          <Tombol type="submit" className="w-full" disabled={mutasi.isPending}>
            {mutasi.isPending ? "Memproses reservasi..." : "Buat Reservasi"}
          </Tombol>
        </form>
      </Kartu>
      <p className="mt-3 text-xs text-[var(--color-teks-sec)]">
        Nomor antrian Anda ditentukan saat reservasi dibuat dan tercetak di invoice.
      </p>
    </div>
  );
}
