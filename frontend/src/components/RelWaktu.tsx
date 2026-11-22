// RelWaktu: motif identitas dari DESIGN.md, garis waktu praktik jam_mulai -> jam_selesai.
// Data nyata (jam praktik), digambar konsisten di kartu jadwal, detail reservasi, dan display.
export function RelWaktu({ jamMulai, jamSelesai }: { jamMulai: string; jamSelesai: string }) {
  const [h1, m1] = jamMulai.split(":").map(Number);
  const [h2, m2] = jamSelesai.split(":").map(Number);
  const menitTotal = 14 * 60; // rentang tampil 07:00 sampai 21:00
  const awal = Math.max(0, ((h1 - 7) * 60 + m1) / menitTotal);
  const akhir = Math.min(1, ((h2 - 7) * 60 + m2) / menitTotal);
  return (
    <div aria-label={`Jam praktik ${jamMulai} sampai ${jamSelesai}`} className="w-full">
      <div className="relative h-2 w-full rounded-full bg-stone-100">
        <div
          className="absolute top-0 h-2 rounded-full bg-[var(--color-aksi)]"
          style={{ left: `${awal * 100}%`, width: `${Math.max(2, (akhir - awal) * 100)}%` }}
        />
      </div>
      <div className="mt-1 flex justify-between text-xs text-[var(--color-teks-sec)]">
        <span>{jamMulai}</span>
        <span>{jamSelesai}</span>
      </div>
    </div>
  );
}
