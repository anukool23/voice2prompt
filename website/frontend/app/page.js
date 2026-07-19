import BgFx from "@/components/BgFx";
import Header from "@/components/Header";
import Hero from "@/components/Hero";
import TrustStrip from "@/components/TrustStrip";
import Features from "@/components/Features";
import Architecture from "@/components/Architecture";
import Roadmap from "@/components/Roadmap";
import DownloadSection from "@/components/DownloadSection";
import Newsletter from "@/components/Newsletter";
import Footer from "@/components/Footer";

export default function Home() {
  return (
    <>
      <BgFx />
      <Header />
      <main id="top">
        <Hero />
        <TrustStrip />
        <Features />
        <Architecture />
        <Roadmap />
        <DownloadSection />
        <Newsletter />
      </main>
      <Footer />
    </>
  );
}
