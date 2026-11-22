"use client";

import { useQuery } from "@tanstack/react-query";
import { Bar, BarChart, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";

import { apiFetch } from "@/lib/api";
import type { Statistik } from "@/lib/tipe";
import { Kartu } from "@/components/Kartu";
import { ErrorState, Memuat } from "@/components/EmptyState";
import { LABEL_STATUS } from "@/lib/format";

const STATUS_URUT = ["menunggu", "dikonfirmasi", "dipanggil", "selesai", "batal"];

// Dashboard admin: angka hari ini, tren 30 hari, dokter terpopuler. Semua dari API nyata.
export default function AdminStatistik() {
  const { data, isLoading, isError, error, refetch } = useQuery({
    queryKey: ["statistik"],
    queryFn: () => apiFetch<Statistik>("/admin/statistik"),
  });

  if (isLoading) return <Memuat proses="Menghitung statistik dari database..." />;
  if (isError) return <ErrorState pesan={(error as Error).message} onCobaLagi={() => refetch()} />;
  if (!data) return null;

  return (
    <>
      <h1 className="text-2xl font-bold">Statistik</h1>
      <p className="mt-1 text-sm text-[var(--color-teks-sec)]">
        Angka dihitung langsung dari data reservasi yang tersimpan.
      </p>

      <div className="mt-6 grid grid-cols-2 gap-4 sm:grid-cols-5">
        {STATUS_URUT.map((s) => (
          <Kartu key={s} className="p-4 text-center">
            <p className="text-xs text-[var(--color-teks-sec)]">{LABEL_STATUS[s]}</p>
            <p className="mt-1 text-3xl font-bold tabular-nums">{data.hari_ini[s] ?? 0}</p>
            <p className="text-xs text-[var(--color-teks-sec)]">hari ini</p>
          </Kartu>
        ))}
      </div>

      <div className="mt-6 grid gap-6 lg:grid-cols-3">
        <Kartu className="p-5 lg:col-span-2">
          <h2 className="text-lg font-semibold">Reservasi per hari, 30 hari terakhir</h2>
          <div className="mt-4 h-64">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={data.per_hari}>
                <XAxis
                  dataKey="tanggal"
                  tickFormatter={(v: string) => v.slice(8)}
                  fontSize={11}
                />
                <YAxis allowDecimals={false} fontSize={11} width={24} />
                <Tooltip
                  formatter={(v) => [`${v} reservasi`, "Total"]}
                  labelFormatter={(l) => `Tanggal ${l}`}
                />
                <Bar dataKey="total" fill="#0f766e" radius={[4, 4, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          </div>
        </Kartu>

        <Kartu className="p-5">
          <h2 className="text-lg font-semibold">Dokter terpopuler</h2>
          {data.dokter_terpopuler.length === 0 ? (
            <p className="mt-3 text-sm text-[var(--color-teks-sec)]">
              Belum ada reservasi untuk dihitung. Daftar muncul setelah ada pasien memesan.
            </p>
          ) : (
            <ol className="mt-3 space-y-2">
              {data.dokter_terpopuler.map((d, i) => (
                <li key={d.nama} className="flex items-center justify-between text-sm">
                  <span>
                    {i + 1}. {d.nama}
                  </span>
                  <span className="font-bold tabular-nums">{d.total}</span>
                </li>
              ))}
            </ol>
          )}
        </Kartu>
      </div>
    </>
  );
}
