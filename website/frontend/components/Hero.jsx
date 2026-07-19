"use client";

import TerminalDemo from "./TerminalDemo";
import ExternalIcon from "./ExternalIcon";
import { trackClick } from "@/lib/analytics";

const GITHUB_URL = "https://github.com/anukool23/voice2prompt";

export default function Hero() {
  return (
    <section className="hero">
      <div className="wrap hero-grid">
        <div>
          <span className="pill">
            <span className="dot"></span> v0.6 · actively shipping new phases
          </span>
          <h1>
            Dictation that never <span className="grad">leaves your machine.</span>
          </h1>
          <p className="lead">
            Voice2Prompt is a free, open-source, on-device voice dictation tool. Hold a key,
            speak, and your words land wherever your cursor is — transcribed and cleaned up
            entirely offline, on macOS and (soon) Windows.
          </p>
          <div className="hero-actions">
            <a
              href="#download"
              className="btn btn-primary"
              onClick={() => trackClick("download", "hero")}
            >
              ⤓ Get the Download Link
            </a>
            <a
              href={GITHUB_URL}
              target="_blank"
              rel="noopener noreferrer"
              className="btn btn-ghost"
              onClick={() => trackClick("github", "hero")}
            >
              View Source on GitHub<ExternalIcon />
            </a>
          </div>
          <div className="hero-meta">
            <div className="item">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M12 2l8 4v6c0 5-3.5 8.5-8 10-4.5-1.5-8-5-8-10V6l8-4z" />
              </svg>
              No network calls
            </div>
            <div className="item">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <circle cx="12" cy="12" r="10" />
                <path d="M12 6v6l4 2" />
              </svg>
              &lt;300ms round trip
            </div>
            <div className="item">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M20 6L9 17l-5-5" />
              </svg>
              MIT licensed
            </div>
          </div>
        </div>

        <div>
          <TerminalDemo />
        </div>
      </div>
    </section>
  );
}
