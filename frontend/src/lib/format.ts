import type { Jadwal, Reservasi } from "./tipe";

export const NAMA_HARI = ["Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"];

export const LABEL_STATUS: Record<string, string> = {
  menunggu: "Menunggu",
  dikonfirmasi: "Dikonfirmasi",
  dipanggil: "Dipanggil",
  selesai: "Selesai",
  batal: "Batal",
};

export function tanggalIndo(iso: string): string {
  const [y, m, d] = iso.split("-").map(Number);
  const t = new Date(y, m - 1, d);
  return `${NAMA_HARI[t.getDay()]}, ${d} ${namaBulan(m - 1)} ${y}`;
}

export function namaBulan(i: number): string {
  return [
    "Januari", "Februari", "Maret", "April", "Mei", "Juni",
    "Juli", "Agustus", "September", "Oktober", "November", "Desember",
  ][i];
}

// Hari praktik ke teks, mis. "Senin, 08:00-12:00"
export function jadwalTeks(j: Pick<Jadwal, "hari" | "jam_mulai" | "jam_selesai">): string {
  return `${NAMA_HARI[j.hari]}, ${j.jam_mulai}-${j.jam_selesai}`;
}

export function bulanIniWIB(): string {
  const now = new Date();
  const wib = new Date(now.getTime() + 7 * 3600_000);
  return wib.toISOString().slice(0, 7);
}
