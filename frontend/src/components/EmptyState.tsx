// EmptyState: menyebut penyebab dan satu aksi berikutnya (R-27), bukan sekadar "tidak ada data".
export function EmptyState({
  judul,
  pesan,
  aksi,
}: {
  judul: string;
  pesan: string;
  aksi?: React.ReactNode;
}) {
  return (
    <div className="flex flex-col items-center gap-2 rounded-[var(--radius-m)] border border-dashed border-[var(--color-garis)] bg-white px-6 py-10 text-center">
      <h3 className="text-sm font-semibold text-[var(--color-teks)]">{judul}</h3>
      <p className="max-w-sm text-sm text-[var(--color-teks-sec)]">{pesan}</p>
      {aksi && <div className="mt-2">{aksi}</div>}
    </div>
  );
}

export function ErrorState({ pesan, onCobaLagi }: { pesan: string; onCobaLagi?: () => void }) {
  return (
    <div className="flex flex-col items-center gap-2 rounded-[var(--radius-m)] border border-[var(--color-bahaya)] bg-[var(--color-bahaya-soft)] px-6 py-8 text-center">
      <h3 className="text-sm font-semibold text-[var(--color-bahaya)]">Gagal memuat data</h3>
      <p className="max-w-sm text-sm text-[var(--color-teks-sec)]">{pesan}</p>
      {onCobaLagi && (
        <button
          onClick={onCobaLagi}
          className="mt-1 min-h-[44px] rounded-[var(--radius-k)] border border-[var(--color-bahaya)] px-4 text-sm font-semibold text-[var(--color-bahaya)]"
        >
          Coba lagi
        </button>
      )}
    </div>
  );
}

// Loading state yang menyebut proses, bukan spinner kosong.
export function Memuat({ proses }: { proses: string }) {
  return (
    <div className="rounded-[var(--radius-m)] border border-dashed border-[var(--color-garis)] bg-white px-6 py-8 text-center text-sm text-[var(--color-teks-sec)]">
      {proses}
    </div>
  );
}
