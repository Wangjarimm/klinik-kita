"use client";

import { useEffect, useState } from "react";
import { wsURL } from "@/lib/api";
import { useMe } from "@/lib/auth";
import toast from "react-hot-toast";

// Langganan WS notifikasi milik user: pesan masuk jadi toast dan tersimpan untuk halaman bell.
export function useNotifikasi(aktif: boolean) {
  const { data: me } = useMe();
  const [daftar, setDaftar] = useState<{ judul: string; pesan: string }[]>([]);

  useEffect(() => {
    if (!aktif || !me) return;
    const ws = new WebSocket(wsURL("/ws/notifikasi"));
    ws.onmessage = (ev) => {
      const msg = JSON.parse(ev.data);
      if (msg.tipe === "notifikasi") {
        setDaftar((d) => [{ judul: msg.judul, pesan: msg.pesan }, ...d]);
        toast(msg.pesan);
      }
    };
    return () => ws.close();
  }, [aktif, me]);

  return daftar;
}
