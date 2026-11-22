export const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

export class ApiError extends Error {
  status: number;
  constructor(status: number, message: string) {
    super(message);
    this.status = status;
  }
}

// Semua panggilan REST lewat sini: cookie ikut, error jadi ApiError.
export async function apiFetch<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${API_URL}/api/v1${path}`, {
    credentials: "include",
    headers: { "Content-Type": "application/json", ...(init?.headers ?? {}) },
    ...init,
  });
  if (!res.ok) {
    let pesan = `Terjadi kesalahan (${res.status})`;
    try {
      const body = await res.json();
      if (body.error) pesan = body.error;
    } catch {
      // body bukan JSON, pakai pesan default
    }
    throw new ApiError(res.status, pesan);
  }
  if (res.status === 204) return undefined as T;
  return res.json() as Promise<T>;
}

// URL WS untuk antrean/notifikasi.
export function wsURL(path: string): string {
  return API_URL.replace(/^http/, "ws") + path;
}
