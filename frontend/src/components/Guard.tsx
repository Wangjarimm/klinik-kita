"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { useMe } from "@/lib/auth";
import { Memuat } from "./EmptyState";

// Guard sisi tampilan: pengalihan cepat ke login. Keamanan sebenarnya di API (role middleware).
export function Guard({ role, children }: { role: "pasien" | "petugas"; children: React.ReactNode }) {
  const { data: me, isError, isLoading } = useMe();
  const router = useRouter();

  useEffect(() => {
    if (isError) router.replace("/login");
  }, [isError, router]);

  if (isLoading) return <Memuat proses="Memeriksa sesi login..." />;
  if (isError) return null;
  if (!me) return null;
  if (me.role !== role) {
    return (
      <main className="flex min-h-screen items-center justify-center p-6">
        <div className="max-w-sm text-center">
          <h1 className="text-lg font-semibold">Akses tidak sesuai</h1>
          <p className="mt-2 text-sm text-[var(--color-teks-sec)]">
            Halaman ini untuk {role}. Anda login sebagai {me.role}.
          </p>
        </div>
      </main>
    );
  }
  return <>{children}</>;
}
