"use client";

// Input dengan label terlihat dan pesan error eksplisit (R-39, R-27).
export function Input({
  label,
  error,
  ...props
}: React.InputHTMLAttributes<HTMLInputElement> & { label: string; error?: string }) {
  const id = props.id ?? props.name;
  return (
    <div className="w-full">
      <label htmlFor={id} className="mb-1 block text-sm font-medium text-[var(--color-teks)]">
        {label}
      </label>
      <input
        id={id}
        className={`min-h-[44px] w-full rounded-[var(--radius-k)] border bg-white px-3 text-sm outline-none transition-colors placeholder:text-stone-400 ${
          error ? "border-[var(--color-bahaya)]" : "border-[var(--color-garis)] focus:border-[var(--color-aksi)]"
        }`}
        {...props}
      />
      {error && <p className="mt-1 text-xs text-[var(--color-bahaya)]">{error}</p>}
    </div>
  );
}

export function Textarea({
  label,
  error,
  ...props
}: React.TextareaHTMLAttributes<HTMLTextAreaElement> & { label: string; error?: string }) {
  const id = props.id ?? props.name;
  return (
    <div className="w-full">
      <label htmlFor={id} className="mb-1 block text-sm font-medium text-[var(--color-teks)]">
        {label}
      </label>
      <textarea
        id={id}
        className={`w-full rounded-[var(--radius-k)] border bg-white px-3 py-2 text-sm outline-none transition-colors placeholder:text-stone-400 ${
          error ? "border-[var(--color-bahaya)]" : "border-[var(--color-garis)] focus:border-[var(--color-aksi)]"
        }`}
        {...props}
      />
      {error && <p className="mt-1 text-xs text-[var(--color-bahaya)]">{error}</p>}
    </div>
  );
}

export function Select({
  label,
  error,
  children,
  ...props
}: React.SelectHTMLAttributes<HTMLSelectElement> & { label: string; error?: string }) {
  const id = props.id ?? props.name;
  return (
    <div className="w-full">
      <label htmlFor={id} className="mb-1 block text-sm font-medium text-[var(--color-teks)]">
        {label}
      </label>
      <select
        id={id}
        className={`min-h-[44px] w-full rounded-[var(--radius-k)] border bg-white px-3 text-sm outline-none ${
          error ? "border-[var(--color-bahaya)]" : "border-[var(--color-garis)]"
        }`}
        {...props}
      >
        {children}
      </select>
      {error && <p className="mt-1 text-xs text-[var(--color-bahaya)]">{error}</p>}
    </div>
  );
}
