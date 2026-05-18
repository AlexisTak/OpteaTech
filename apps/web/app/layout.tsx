import type { Metadata } from "next";
import { Geist, Geist_Mono, Instrument_Serif } from "next/font/google";
import { Footer } from "@/components/shared/Footer";
import { Navbar } from "@/components/shared/Navbar";
import { PageCurtain } from "@/components/shared/PageCurtain";
import { AnimationProvider } from "@/providers/AnimationProvider";
import { LenisProvider } from "@/providers/LenisProvider";
import "./globals.css";

const instrumentSerif = Instrument_Serif({
  variable: "--font-display",
  subsets: ["latin"],
  weight: "400",
  display: "swap",
});

const geist = Geist({
  variable: "--font-body",
  subsets: ["latin"],
  display: "swap",
});

const geistMono = Geist_Mono({
  variable: "--font-mono",
  subsets: ["latin"],
  display: "swap",
});

export const metadata: Metadata = {
  metadataBase: new URL("https://opteaetch.fr"),
  title: {
    default: "optea.tech - Agence web et IA sur mesure",
    template: "%s | optea.tech",
  },
  description:
    "optea.tech conçoit des sites web, logiciels et experiences IA sur mesure. Expertise Next.js et Go. Devis gratuit sous 48h.",
  keywords: ["agence web", "developpement logiciel", "solutions IA", "Next.js", "sur mesure", "France"],
  openGraph: {
    type: "website",
    locale: "fr_FR",
    url: "https://opteaetch.fr",
    siteName: "optea.tech",
    images: [{ url: "/og-image.png", width: 1200, height: 630 }],
  },
  twitter: { card: "summary_large_image" },
  robots: { index: true, follow: true },
  alternates: { canonical: "https://opteaetch.fr" },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const orgJsonLd = {
    "@context": "https://schema.org",
    "@type": "Organization",
    name: "optea.tech",
    url: "https://opteaetch.fr",
    email: "hello@opteaetch.fr",
    sameAs: ["https://github.com", "https://www.linkedin.com"],
  };

  return (
    <html
      lang="fr"
      className={`${instrumentSerif.variable} ${geist.variable} ${geistMono.variable} h-full antialiased`}
    >
      <body suppressHydrationWarning className="min-h-full bg-[var(--off-white)] font-body text-[var(--ink)]">
        <LenisProvider>
          <AnimationProvider>
            <PageCurtain />
            <Navbar />
            <main>{children}</main>
            <Footer />
            <script
              type="application/ld+json"
              dangerouslySetInnerHTML={{ __html: JSON.stringify(orgJsonLd) }}
            />
          </AnimationProvider>
        </LenisProvider>
      </body>
    </html>
  );
}
