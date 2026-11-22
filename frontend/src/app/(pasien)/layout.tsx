"use client";

import Link from "next/link";
import { Bell, LogOut } from "lucide-react";
import { Guard } from "@/components/Guard";
import { Tombol } from "@/components/Tombol";
import { useNotifikasi } from "@/lib/useNotifikasi";
import { logout } from "@/lib/auth";
import { useState } from "react";

// Shell halaman pasien: navigasi atas + bell notifikasi + keluar.
export default function LayoutPasien({ children }: { children: React.ReactNode }) {
  return (
    <Guard role="pasien">
      <ShellPasien>{children}</ShellPasien>
    </Guard>
  );
}

function ShellPasien({ children }: { children: React.ReactNode }) {
  const daftarNotif = useNotifikasi(true);
  const [bukaBell, setBukaBell] = useState(false);

  return (
    <>
      <header>
        <nav className="border-b border-[var(--color-garis)] bg-white">
          <div className="mx-auto flex max-w-5xl items-center justify-between px-4 py-3">
            <Link href="/dokter" className="text-lg font-bold">
              Klinik Kita
            </Link>
            <div className="flex items-center gap-2">
              <Link href="/reservasi">
                <Tombol varian="ghost">Reservasi Saya</Tombol>
              </Link>
              <div className="relative">
                <button
                  onClick={() => setBukaBell((v) => !v)}
                  aria-label={`Notifikasi, ${daftarNotif.length} pesan baru`}
                  className="relative flex min-h-[44px] w-11 items-center justify-center rounded-[var(--radius-k)] text-[var(--color-teks-sec)] hover:bg-stone-100"
                >
                  <Bell size={20} aria-hidden />
                  {daftarNotif.length > 0 && (
                    <span className="absolute right-1.5 top-1.5 h-2.5 w-2.5 rounded-full bg-[var(--color-amber)]" />
                  )}
                </button>
                {bukaBell && (
                  <div className="absolute right-0 top-12 z-10 w-80 rounded-[var(--radius-m)] border border-[var(--color-garis)] bg-white p-3 shadow-lg">
                    {daftarNotif.length === 0 ? (
                      <p className="py-2 text-sm text-[var(--color-teks-sec)]">
                        Belum ada notifikasi baru. Status reservasi Anda akan muncul di sini.
                      </p>
                    ) : (
                      <ul className="max-h-72 space-y-2 overflow-auto">
                        {daftarNotif.map((n, i) => (
                          <li key={i} className="rounded-[var(--radius-k)] bg-stone-50 p-2 text-sm">
                            <p className="font-semibold">{n.judul}</p>
                            <p className="text-[var(--color-teks-sec)]">{n.pesan}</p>
                          </li>
                        ))}
                      </ul>
                    )}
                  </div>
                )}
              </div>
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
          </div>
        </nav>
      </header>
      <main className="mx-auto max-w-5xl px-4 py-6">{children}</main>
    </>
  );
}
