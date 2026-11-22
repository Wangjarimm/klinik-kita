// Badge status reservasi: warna menandai keadaan nyata (DESIGN.md), bukan dekorasi.
const WARNA: Record<string, string> = {
  menunggu: "bg-[var(--color-amber-soft)] text-[var(--color-amber)]",
  dikonfirmasi: "border border-[var(--color-aksi)] text-[var(--color-aksi)] bg-white",
  dipanggil: "bg-[var(--color-aksi)] text-white",
  selesai: "bg-stone-100 text-stone-500",
  batal: "bg-[var(--color-bahaya-soft)] text-[var(--color-bahaya)]",
};

const LABEL: Record<string, string> = {
  menunggu: "Menunggu",
  dikonfirmasi: "Dikonfirmasi",
  dipanggil: "Dipanggil",
  selesai: "Selesai",
  batal: "Batal",
};

export function BadgeStatus({ status }: { status: string }) {
  return (
    <span className={`inline-block rounded-[var(--radius-k)] px-2 py-0.5 text-xs font-semibold ${WARNA[status] ?? "bg-stone-100 text-stone-500"}`}>
      {LABEL[status] ?? status}
    </span>
  );
}
