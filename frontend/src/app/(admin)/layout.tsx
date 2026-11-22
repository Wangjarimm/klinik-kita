"use client";

import Link from "next/link";
import { LogOut } from "lucide-react";
import { Guard } from "@/components/Guard";
import { Tombol } from "@/components/Tombol";
import { logout } from "@/lib/auth";

const MENU = [
  { href: "/admin", label: "Statistik" },
  { href: "/admin/reservasi", label: "Reservasi & Antrean" },
  { href: "/admin/dokter", label: "Dokter & Jadwal" },
  { href: "/admin/penutupan", label: "Tutup Jadwal" },
  { href: "/admin/pasien", label: "Pasien & Riwayat" },
];

export default function LayoutAdmin({ children }: { children: React.ReactNode }) {
  return (
    <Guard role="petugas">
      <div className="min-h-screen">
        <header>
          <nav className="border-b border-[var(--color-garis)] bg-white">
            <div className="mx-auto max-w-6xl px-4 py-3">
              <div className="flex items-center justify-between">
                <Link href="/admin" className="text-lg font-bold">
                  Klinik Kita (Petugas)
                </Link>
                <Tombol
                  varian="outline"
                  onClick={async () => {
                    await logout();
                    window.location.href = "/";
                  }}
                >
                  <LogOut size={16} aria-hidden />
                  Keluar
                </Tombol>
              </div>
              <ul className="mt-2 flex flex-wrap gap-1">
                {MENU.map((m) => (
                  <li key={m.href}>
                    <Link
                      href={m.href}
                      className="inline-flex min-h-[44px] items-center rounded-[var(--radius-k)] px-3 text-sm font-semibold text-[var(--color-teks-sec)] hover:bg-stone-100 hover:text-[var(--color-teks)]"
                    >
                      {m.label}
                    </Link>
                  </li>
                ))}
              </ul>
            </div>
          </nav>
        </header>
        <main className="mx-auto max-w-6xl px-4 py-6">{children}</main>
      </div>
    </Guard>
  );
}
