"use client";

import { useEffect, useState } from "react";
import { wsURL } from "@/lib/api";
import type { SnapshotAntrean } from "@/lib/tipe";

// Langganan WS antrean untuk satu dokter+tanggal: snapshot awal lalu update live.
export function useAntrean(dokterID: number, tanggal: string, aktif: boolean) {
  const [snap, setSnap] = useState<SnapshotAntrean | null>(null);
  const [terputus, setTerputus] = useState(false);

  useEffect(() => {
    if (!aktif || !dokterID || !tanggal) return;
    const ws = new WebSocket(`${wsURL(`/ws/antrean?dokter_id=${dokterID}&tanggal=${tanggal}`)}`);
    ws.onopen = () => setTerputus(false);
    ws.onmessage = (ev) => setSnap(JSON.parse(ev.data));
    ws.onclose = () => setTerputus(true);
    return () => ws.close();
  }, [dokterID, tanggal, aktif]);

  return { snap, terputus };
}
