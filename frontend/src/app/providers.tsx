"use client";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { Toaster } from "react-hot-toast";
import { useState } from "react";

export function Providers({ children }: { children: React.ReactNode }) {
  const [client] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: { retry: 1, staleTime: 15_000, refetchOnWindowFocus: false },
        },
      }),
  );
  return (
    <QueryClientProvider client={client}>
      {children}
      {/* posisi bawah tengah, kontras tinggi, tanpa animasi berlebih (MOTION 1) */}
      <Toaster
        position="bottom-center"
        toastOptions={{
          style: { borderRadius: "var(--radius-m)", background: "#1c1917", color: "#fafaf9" },
        }}
      />
    </QueryClientProvider>
  );
}
