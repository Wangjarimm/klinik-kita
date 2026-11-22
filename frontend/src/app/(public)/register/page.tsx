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
  username: z.string().min(3, "Username minimal 3 karakter"),
  email: z.string().email("Format email tidak valid"),
  password: z.string().min(8, "Password minimal 8 karakter"),
});
type FormNilai = z.infer<typeof schema>;

export default function Register() {
  const router = useRouter();
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<FormNilai>({ resolver: zodResolver(schema) });

  const mutasi = useMutation({
    mutationFn: (nilai: FormNilai) =>
      apiFetch<{ message: string }>("/auth/register", { method: "POST", body: JSON.stringify(nilai) }),
    onSuccess: () => {
      toast.success("Registrasi berhasil. Silakan masuk.");
      router.push("/login");
    },
    onError: (e: Error) => toast.error(e.message),
  });

  return (
    <main className="flex min-h-screen items-center justify-center px-4">
      <Kartu className="w-full max-w-sm p-6">
        <h1 className="text-xl font-bold">Daftar akun pasien</h1>
        <p className="mt-1 text-sm text-[var(--color-teks-sec)]">
          Setelah masuk, lengkapi profil dan buat reservasi pertama Anda.
        </p>
        <form className="mt-5 space-y-4" onSubmit={handleSubmit((v) => mutasi.mutate(v))}>
          <Input label="Username" {...register("username")} error={errors.username?.message} />
          <Input label="Email" type="email" {...register("email")} error={errors.email?.message} />
          <Input
            label="Password"
            type="password"
            {...register("password")}
            error={errors.password?.message}
          />
          <Tombol type="submit" className="w-full" disabled={mutasi.isPending}>
            {mutasi.isPending ? "Memproses..." : "Daftar"}
          </Tombol>
        </form>
        <p className="mt-4 text-center text-sm text-[var(--color-teks-sec)]">
          Sudah punya akun?{" "}
          <Link href="/login" className="font-semibold text-[var(--color-aksi)]">
            Masuk di sini
          </Link>
        </p>
      </Kartu>
    </main>
  );
}
