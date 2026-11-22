// Bentuk data dari API (backend Go).

export type Dokter = {
  id: number;
  nid: string;
  nama: string;
  keahlian: string;
  no_hp: string;
};

export type Jadwal = {
  id: number;
  dokter_id: number;
  hari: number;
  jam_mulai: string;
  jam_selesai: string;
  kuota: number;
};

export type HariStatus = {
  tanggal: string;
  status: "tersedia" | "penuh" | "tutup" | "tidak_praktik";
  sisa: number;
  alasan: string;
};

export type Reservasi = {
  id: number;
  kode: string;
  tanggal_kunjungan: string;
  keluhan: string;
  status: "menunggu" | "dikonfirmasi" | "dipanggil" | "selesai" | "batal";
  nomor_antrian: number;
  catatan_admin: string;
  dokter?: { id: number; nama: string; keahlian: string };
  jadwal?: { hari: number; jam_mulai: string; jam_selesai: string };
};

export type Penutupan = {
  id: number;
  dokter_id: number | null;
  tanggal: string;
  alasan: string;
};

export type Pasien = {
  id: number;
  user_id: number | null;
  nama_lengkap: string;
  nik: string;
  jenis_kelamin: string;
  tempat_lahir: string;
  tanggal_lahir: string;
  alamat: string;
  no_hp: string;
};

export type SnapshotAntrean = {
  dokter_id: number;
  tanggal: string;
  sedang_melayani: number;
  sisa_menunggu: number;
};

export type Statistik = {
  hari_ini: Record<string, number>;
  per_hari: { tanggal: string; total: number }[];
  dokter_terpopuler: { nama: string; total: number }[];
};
