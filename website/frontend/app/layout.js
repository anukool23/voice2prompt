import "./globals.css";
import Analytics from "@/components/Analytics";

// TODO: swap this for your real production domain once you deploy (Vercel gives you
// a *.vercel.app URL by default) — it drives canonical links, OG/Twitter tags and the sitemap.
const SITE_URL = "https://voice2prompt.app";

export const metadata = {
  metadataBase: new URL(SITE_URL),
  title: {
    default: "Voice2Prompt — Open-Source, On-Device Voice Dictation",
    template: "%s · Voice2Prompt",
  },
  description:
    "Voice2Prompt is a free, open-source, on-device voice dictation tool for macOS and Windows. Hold a key, speak, and your words are typed instantly — no cloud, no per-word cost, 100% private.",
  keywords: [
    "voice dictation",
    "speech to text",
    "open source dictation",
    "on-device speech recognition",
    "whisper.cpp",
    "ollama",
    "privacy-first dictation",
    "macOS dictation app",
    "Windows dictation app",
    "voice to text software",
  ],
  authors: [{ name: "Voice2Prompt" }],
  creator: "Voice2Prompt",
  applicationName: "Voice2Prompt",
  category: "technology",
  openGraph: {
    title: "Voice2Prompt — Open-Source, On-Device Voice Dictation",
    description:
      "Free, open-source, on-device voice dictation for macOS + Windows. No cloud calls, no per-word cost, 100% private.",
    url: SITE_URL,
    siteName: "Voice2Prompt",
    locale: "en_US",
    type: "website",
  },
  twitter: {
    card: "summary_large_image",
    title: "Voice2Prompt — Open-Source, On-Device Voice Dictation",
    description: "Free, open-source, on-device voice dictation for macOS + Windows.",
  },
  robots: {
    index: true,
    follow: true,
    googleBot: {
      index: true,
      follow: true,
      "max-image-preview": "large",
      "max-snippet": -1,
    },
  },
  alternates: {
    canonical: SITE_URL,
  },
};

export const viewport = {
  width: "device-width",
  initialScale: 1,
  maximumScale: 5,
  themeColor: "#07070c",
};

const jsonLd = {
  "@context": "https://schema.org",
  "@type": "SoftwareApplication",
  name: "Voice2Prompt",
  applicationCategory: "UtilitiesApplication",
  operatingSystem: "macOS, Windows",
  description:
    "A free, open-source, on-device voice dictation tool for macOS and Windows. No cloud calls, no per-word cost.",
  offers: {
    "@type": "Offer",
    price: "0",
    priceCurrency: "USD",
  },
  url: SITE_URL,
  license: "https://github.com/anukool23/voice2prompt/blob/main/LICENSE",
};

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <head>
        <link rel="preconnect" href="https://fonts.googleapis.com" />
        <link rel="preconnect" href="https://fonts.gstatic.com" crossOrigin="anonymous" />
        <link
          href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800&family=JetBrains+Mono:wght@400;500;600;700&display=swap"
          rel="stylesheet"
        />
        <script
          type="application/ld+json"
          // eslint-disable-next-line react/no-danger
          dangerouslySetInnerHTML={{ __html: JSON.stringify(jsonLd) }}
        />
      </head>
      <body>
        <Analytics />
        {children}
      </body>
    </html>
  );
}
