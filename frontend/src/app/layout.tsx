import type { Metadata } from "next";
import { Plus_Jakarta_Sans } from "next/font/google";
import "./globals.css";
import { Providers } from "./providers";

const jakarta = Plus_Jakarta_Sans({
  subsets: ["latin"],
  variable: "--font-jakarta",
  display: "swap",
});

export const metadata: Metadata = {
  title: "Klinik Kita",
  description: "Reservasi dokter online Klinik Kita",
};

export default function RootLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="id" className={jakarta.variable}>
      <body className="font-sans antialiased" style={{ fontFamily: "var(--font-jakarta), sans-serif" }}>
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}
