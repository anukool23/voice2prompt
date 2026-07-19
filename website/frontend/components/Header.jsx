"use client";

import { useState } from "react";
import Logo from "./Logo";
import ExternalIcon from "./ExternalIcon";
import { trackClick } from "@/lib/analytics";

const GITHUB_URL = "https://github.com/anukool23/voice2prompt";
const PORTFOLIO_URL = "https://www.anukool.me";

export default function Header() {
  const [open, setOpen] = useState(false);

  return (
    <header>
      <div className="wrap">
        <nav>
          <a href="#top" className="brand">
            <Logo gradId="logo-nav" />
            Voice2Prompt
          </a>
          <div className="nav-links">
            <a href="#features">Features</a>
            <a href="#how-it-works">How it Works</a>
            <a href="#roadmap">Roadmap</a>
            <a href="#newsletter">Newsletter</a>
            <a
              href={GITHUB_URL}
              target="_blank"
              rel="noopener noreferrer"
              onClick={() => trackClick("github", "header")}
            >
              GitHub<ExternalIcon />
            </a>
            <a
              href={PORTFOLIO_URL}
              target="_blank"
              rel="noopener noreferrer"
              onClick={() => trackClick("developer", "header")}
            >
              Developer<ExternalIcon />
            </a>
          </div>
          <div className="nav-cta">
            <a
              className="btn btn-ghost btn-sm"
              href={GITHUB_URL}
              target="_blank"
              rel="noopener noreferrer"
              onClick={() => trackClick("github", "header")}
            >
              ★ Star<ExternalIcon />
            </a>
            <a
              className="btn btn-primary btn-sm"
              href="#download"
              onClick={() => trackClick("download", "header")}
            >
              Download
            </a>
          </div>
          <button
            className="nav-mobile-toggle"
            aria-label="Toggle menu"
            aria-expanded={open}
            onClick={() => setOpen((v) => !v)}
          >
            {open ? "✕" : "☰"}
          </button>
        </nav>
        <div className={`mobile-menu${open ? " open" : ""}`}>
          <a href="#features" onClick={() => setOpen(false)}>Features</a>
          <a href="#how-it-works" onClick={() => setOpen(false)}>How it Works</a>
          <a href="#roadmap" onClick={() => setOpen(false)}>Roadmap</a>
          <a href="#newsletter" onClick={() => setOpen(false)}>Newsletter</a>
          <a
            href={GITHUB_URL}
            target="_blank"
            rel="noopener noreferrer"
            onClick={() => trackClick("github", "header")}
          >
            GitHub<ExternalIcon />
          </a>
          <a
            href={PORTFOLIO_URL}
            target="_blank"
            rel="noopener noreferrer"
            onClick={() => trackClick("developer", "header")}
          >
            Developer<ExternalIcon />
          </a>
          <a
            href="#download"
            onClick={() => {
              trackClick("download", "header");
              setOpen(false);
            }}
          >
            Download
          </a>
        </div>
      </div>
    </header>
  );
}
