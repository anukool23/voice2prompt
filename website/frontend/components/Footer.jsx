"use client";

import Logo from "./Logo";
import ExternalIcon from "./ExternalIcon";
import { trackClick } from "@/lib/analytics";

const GITHUB_URL = "https://github.com/anukool23/voice2prompt";
const PORTFOLIO_URL = "https://www.anukool.me";

export default function Footer() {
  const year = new Date().getFullYear();

  return (
    <footer>
      <div className="wrap">
        <div className="footer-grid">
          <div className="footer-brand">
            <a href="#top" className="brand">
              <Logo gradId="logo-footer" />
              Voice2Prompt
            </a>
            <p>
              An open-source, on-device voice-dictation tool for macOS + Windows. Go +
              whisper.cpp + Ollama + Wails/React.
            </p>
          </div>
          <div className="footer-col">
            <h5>Project</h5>
            <a
              href={GITHUB_URL}
              target="_blank"
              rel="noopener noreferrer"
              onClick={() => trackClick("github", "footer")}
            >
              GitHub<ExternalIcon />
            </a>
            <a href={`${GITHUB_URL}/issues`} target="_blank" rel="noopener noreferrer">Issues<ExternalIcon /></a>
            <a href={`${GITHUB_URL}/discussions`} target="_blank" rel="noopener noreferrer">Discussions<ExternalIcon /></a>
            <a href={`${GITHUB_URL}/blob/main/LICENSE`} target="_blank" rel="noopener noreferrer">MIT License<ExternalIcon /></a>
          </div>
          <div className="footer-col">
            <h5>Docs</h5>
            <a href="#how-it-works">Architecture</a>
            <a href="#roadmap">Roadmap</a>
            <a href={`${GITHUB_URL}#readme`} target="_blank" rel="noopener noreferrer">README<ExternalIcon /></a>
            <a href={`${GITHUB_URL}/blob/main/flow-clone-implementation-plan.md`} target="_blank" rel="noopener noreferrer">
              Implementation plan<ExternalIcon />
            </a>
          </div>
          <div className="footer-col">
            <h5>Get it</h5>
            <a href="#download" onClick={() => trackClick("download", "footer")}>Download</a>
            <a href="#newsletter">Newsletter</a>
            <a href={`${GITHUB_URL}/releases`} target="_blank" rel="noopener noreferrer">Releases<ExternalIcon /></a>
          </div>
        </div>
        <div className="footer-bottom">
          <p>
            © {year} Voice2Prompt. Built in the open, MIT licensed. · Built by{" "}
            <a
              href={PORTFOLIO_URL}
              target="_blank"
              rel="noopener noreferrer"
              onClick={() => trackClick("developer", "footer")}
            >
              Anukool<ExternalIcon />
            </a>
          </p>
          <div className="footer-social">
            <a
              href={GITHUB_URL}
              target="_blank"
              rel="noopener noreferrer"
              aria-label="GitHub"
              onClick={() => trackClick("github", "footer")}
            >
              <svg viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 .5C5.65.5.5 5.65.5 12c0 5.08 3.29 9.39 7.86 10.91.57.1.78-.25.78-.55 0-.27-.01-1.16-.02-2.11-3.2.7-3.87-1.36-3.87-1.36-.53-1.34-1.29-1.7-1.29-1.7-1.05-.72.08-.7.08-.7 1.16.08 1.77 1.19 1.77 1.19 1.03 1.77 2.71 1.26 3.37.96.1-.75.4-1.26.73-1.55-2.56-.29-5.26-1.28-5.26-5.7 0-1.26.45-2.29 1.19-3.1-.12-.29-.52-1.46.11-3.05 0 0 .97-.31 3.18 1.18a11 11 0 0 1 5.79 0c2.2-1.49 3.17-1.18 3.17-1.18.63 1.59.23 2.76.11 3.05.74.81 1.18 1.84 1.18 3.1 0 4.43-2.7 5.4-5.28 5.69.42.36.78 1.07.78 2.16 0 1.56-.01 2.82-.01 3.2 0 .31.2.66.79.55A10.52 10.52 0 0 0 23.5 12C23.5 5.65 18.35.5 12 .5z" />
              </svg>
            </a>
          </div>
        </div>
      </div>
    </footer>
  );
}
