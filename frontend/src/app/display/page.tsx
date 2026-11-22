"use client";

import { Suspense } from "react";
import { useSearchParams } from "next/navigation";
import { useQuery } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { Dokter } from "@/lib/tipe";
import { useAntrean } from "@/lib/useAntrean";

// Layar antrian fullscreen untuk TV klinik: nomor sedang melayani, sisa menunggu.
// Akses via link tanpa login (nomor antrean bukan data sensitif).
export default function Display() {
  return (
    <Suspense fallback={<main className="min-h-screen bg-[var(--color-teks)]" />}>
      <DisplayKonten />
    </Suspense>
  );
}

function DisplayKonten() {
  const params = useSearchParams();
  const dokterID = Number(params.get("dokter") ?? 0);
  const tanggal = params.get("tanggal") ?? "";
  const { snap, terputus } = useAntrean(dokterID, tanggal, dokterID > 0 && tanggal !== "");
  const { data: detail } = useQuery({
    queryKey: ["dokter-detail", dokterID],
    queryFn: () => apiFetch<{ data: { dokter: Dokter } }>(`/dokter/${dokterID}`),
    enabled: dokterID > 0,
  });

  if (!dokterID || !tanggal) {
    return (
      <main className="flex min-h-screen items-center justify-center bg-[var(--color-teks)] p-6">
        <p className="max-w-md text-center text-lg text-stone-200">
          Halaman display memerlukan parameter dokter dan tanggal.
          <br />
          Format: <span className="text-white">/display?dokter=ID&amp;tanggal=YYYY-MM-DD</span>
        </p>
      </main>
    );
  }

  return (
    <main className="flex min-h-screen flex-col items-center justify-center gap-10 bg-[var(--color-teks)] p-8 text-center">
      <div>
        <p className="text-lg font-semibold text-stone-300">Klinik Kita</p>
        <p className="text-sm text-stone-400">
          {detail?.data.dokter.nama ?? "Dokter"} · {tanggal}
        </p>
      </div>
      <div>
        <p className="text-xl text-stone-300">Sedang melayani</p>
        {snap && snap.sedang_melayani > 0 ? (
          <p className="text-[10rem] font-bold leading-none tabular-nums text-white">
            {String(snap.sedang_melayani).padStart(2, "0")}
          </p>
        ) : (
          <p className="text-5xl font-semibold leading-none text-stone-400">Belum mulai</p>
        )}
      </div>
      <p className="text-2xl text-stone-300">
        Sisa menunggu: <span className="font-bold tabular-nums">{snap?.sisa_menunggu ?? "-"}</span>
      </p>
      {terputus && (
        <p className="text-sm text-[var(--color-amber-soft)]">
          Koneksi terputus. Muat ulang halaman untuk menyambung kembali.
        </p>
      )}
    </main>
  );
}
