"use client";

import { useQuery } from "@tanstack/react-query";
import { apiFetch } from "@/lib/api";
import type { HariStatus } from "@/lib/tipe";
import { namaBulan } from "@/lib/format";

type Status = HariStatus["status"];

const GAYA: Record<Status, string> = {
  tersedia: "bg-[var(--color-aksi-soft)] text-[var(--color-aksi)] hover:bg-[var(--color-aksi)] hover:text-white",
  penuh: "bg-stone-100 text-stone-400",
  tutup: "bg-stone-100 text-stone-400",
  tidak_praktik: "bg-white text-stone-300",
};

// Kalender bulan berjalan: status per tanggal dari API ketersediaan.
// Tanggal tak tersedia dinonaktifkan; kliknya tetap memberi alasan (toast oleh pemanggil).
export function Kalender({
  dokterID,
  bulan,
  onBulan,
  pilih,
  terpilih,
}: {
  dokterID: number;
  bulan: string; // YYYY-MM
  onBulan: (b: string) => void;
  pilih: (tanggal: string, status: Status, alasan: string, sisa: number) => void;
  terpilih?: string;
}) {
  const { data, isLoading, isError, error } = useQuery({
    queryKey: ["ketersediaan", dokterID, bulan],
    queryFn: () => apiFetch<{ hari: HariStatus[] }>(`/dokter/${dokterID}/ketersediaan?bulan=${bulan}`),
  });

  const map = new Map<string, HariStatus>();
  data?.hari.forEach((h) => map.set(h.tanggal, h));

  const [y, m] = bulan.split("-").map(Number);
  const firstDow = new Date(y, m - 1, 1).getDay(); // 0=Minggu
  const offset = (firstDow + 6) % 7; // kolom mulai Senin
  const jmlHari = new Date(y, m, 0).getDate();
  const hariIniISO = new Date(Date.now() + 7 * 3600_000).toISOString().slice(0, 10);

  const geserBulan = (delta: number) => {
    const d = new Date(y, m - 1 + delta, 1);
    onBulan(`${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, "0")}`);
  };

  return (
    <div>
      <div className="mb-3 flex items-center justify-between">
        <button
          onClick={() => geserBulan(-1)}
          className="min-h-[44px] rounded-[var(--radius-k)] px-3 text-sm font-semibold text-[var(--color-teks-sec)] hover:bg-stone-100"
          aria-label="Bulan sebelumnya"
        >
          ←
        </button>
        <h3 className="text-sm font-semibold">
          {namaBulan(m - 1)} {y}
        </h3>
        <button
          onClick={() => geserBulan(1)}
          className="min-h-[44px] rounded-[var(--radius-k)] px-3 text-sm font-semibold text-[var(--color-teks-sec)] hover:bg-stone-100"
          aria-label="Bulan berikutnya"
        >
          →
        </button>
      </div>

      {isLoading && <p className="py-6 text-center text-sm text-[var(--color-teks-sec)]">Memeriksa ketersediaan jadwal...</p>}
      {isError && (
        <p className="py-6 text-center text-sm text-[var(--color-bahaya)]">
          Gagal memuat ketersediaan: {(error as Error).message}
        </p>
      )}

      {!isLoading && !isError && (
        <>
          <div className="grid grid-cols-7 gap-1 text-center text-xs font-medium text-[var(--color-teks-sec)]">
            {["Sen", "Sel", "Rab", "Kam", "Jum", "Sab", "Min"].map((h) => (
              <div key={h} className="py-1">
                {h}
              </div>
            ))}
          </div>
          <div className="grid grid-cols-7 gap-1">
            {Array.from({ length: offset }).map((_, i) => (
              <div key={`kosong-${i}`} />
            ))}
            {data?.hari
              .filter((h) => h.tanggal.slice(0, 7) === bulan)
              .map((h) => {
                const nomor = Number(h.tanggal.slice(8, 10));
                const bisa = h.status === "tersedia";
                return (
                  <button
                    key={h.tanggal}
                    aria-disabled={!bisa}
                    onClick={() => pilih(h.tanggal, h.status, h.alasan, h.sisa)}
                    title={h.alasan || (bisa ? `${h.sisa} slot tersisa` : undefined)}
                    className={`min-h-[44px] rounded-[var(--radius-k)] text-sm font-semibold transition-colors ${GAYA[h.status]} ${
                      terpilih === h.tanggal ? "ring-2 ring-[var(--color-aksi)]" : ""
                    } ${bisa ? "cursor-pointer" : "cursor-not-allowed"}`}
                  >
                    {nomor}
                    {bisa && <span className="block text-[10px] font-normal">sisa {h.sisa}</span>}
                    {h.tanggal === hariIniISO && (
                      <span className="block text-[10px] font-normal underline">hari ini</span>
                    )}
                  </button>
                );
              })}
          </div>
        </>
      )}
    </div>
  );
}
