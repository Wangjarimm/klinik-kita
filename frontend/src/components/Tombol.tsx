"use client";

// Tombol: satu komponen untuk semua aksi. Radius 8 (bukan pill), aksen teal hanya di primary.
const gaya: Record<string, string> = {
  primary: "bg-[var(--color-aksi)] text-white hover:bg-[var(--color-aksi-hover)]",
  outline:
    "border border-[var(--color-garis)] bg-white text-[var(--color-teks)] hover:border-[var(--color-aksi)] hover:text-[var(--color-aksi)]",
  bahaya: "border border-[var(--color-bahaya)] bg-white text-[var(--color-bahaya)] hover:bg-[var(--color-bahaya-soft)]",
  ghost: "text-[var(--color-teks-sec)] hover:text-[var(--color-teks)] hover:bg-stone-100",
};

export function Tombol({
  varian = "primary",
  className = "",
  ...props
}: React.ButtonHTMLAttributes<HTMLButtonElement> & { varian?: keyof typeof gaya }) {
  return (
    <button
      className={`inline-flex min-h-[44px] items-center justify-center gap-2 rounded-[var(--radius-k)] px-4 text-sm font-semibold transition-colors disabled:cursor-not-allowed disabled:opacity-50 ${gaya[varian]} ${className}`}
      {...props}
    />
  );
}
