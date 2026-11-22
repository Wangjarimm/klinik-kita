// Kartu: permukaan putih, garis tipis, radius 12. Tanpa bayangan (bayangan hanya untuk elemen mengambang).
export function Kartu({
  className = "",
  ...props
}: React.HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      className={`rounded-[var(--radius-m)] border border-[var(--color-garis)] bg-[var(--color-surface)] ${className}`}
      {...props}
    />
  );
}
