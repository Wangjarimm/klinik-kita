"use client";

import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useMutation } from "@tanstack/react-query";
import toast from "react-hot-toast";
import Link from "next/link";

import { apiFetch } from "@/lib/api";
import { Input } from "@/components/Input";
import { Tombol } from "@/components/Tombol";
import { Kartu } from "@/components/Kartu";

const schema = z.object({
  username: z.string().min(1, "Username wajib diisi"),
  password: z.string().min(1, "Password wajib diisi"),
});
type FormNilai = z.infer<typeof schema>;

export default function Login() {
  const router = useRouter();
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<FormNilai>({ resolver: zodResolver(schema) });

  const mutasi = useMutation({
    mutationFn: (nilai: FormNilai) =>
      apiFetch<{ role: string }>("/auth/login", { method: "POST", body: JSON.stringify(nilai) }),
    onSuccess: (hasil) => {
      toast.success("Login berhasil");
      router.push(hasil.role === "petugas" ? "/admin" : "/dokter");
    },
    onError: (e: Error) => toast.error(e.message),
  });

  return (
    <main className="flex min-h-screen items-center justify-center px-4">
      <Kartu className="w-full max-w-sm p-6">
        <h1 className="text-xl font-bold">Masuk ke Klinik Kita</h1>
        <p className="mt-1 text-sm text-[var(--color-teks-sec)]">
          Akun pasien maupun petugas klinik memakai halaman ini.
        </p>
        <form
          className="mt-5 space-y-4"
          onSubmit={handleSubmit((v) => mutasi.mutate(v))}
        >
          <Input label="Username" {...register("username")} error={errors.username?.message} />
          <Input
            label="Password"
            type="password"
            {...register("password")}
            error={errors.password?.message}
          />
          <Tombol type="submit" className="w-full" disabled={mutasi.isPending}>
            {mutasi.isPending ? "Memproses..." : "Masuk"}
          </Tombol>
        </form>
        <p className="mt-4 text-center text-sm text-[var(--color-teks-sec)]">
          Belum punya akun?{" "}
          <Link href="/register" className="font-semibold text-[var(--color-aksi)]">
            Daftar di sini
          </Link>
        </p>
      </Kartu>
    </main>
  );
}
