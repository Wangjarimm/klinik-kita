"use client";

import { useQuery } from "@tanstack/react-query";
import { apiFetch } from "./api";

export type Me = { user_id: number; username: string; role: "pasien" | "petugas" };

export function useMe() {
  return useQuery({
    queryKey: ["me"],
    queryFn: () => apiFetch<Me>("/auth/me"),
    retry: false,
  });
}

// Logout menghapus sesi lalu pemanggil mengarahkan ke landing.
export async function logout(): Promise<void> {
  await apiFetch("/auth/logout", { method: "POST" });
}
