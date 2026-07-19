"use client";

import { useState } from "react";
import Reveal from "./Reveal";
import ExternalIcon from "./ExternalIcon";
import { trackClick } from "@/lib/analytics";

const GITHUB_RELEASES_URL = "https://github.com/anukool23/voice2prompt/releases";
const API_BASE = process.env.NEXT_PUBLIC_API_URL || "";

export default function DownloadSection() {
  const [platform, setPlatform] = useState("mac");
  const [email, setEmail] = useState("");
  const [submitted, setSubmitted] = useState(false);
  const [submittedEmail, setSubmittedEmail] = useState("");
  const [result, setResult] = useState(null); // { downloadUrl, releaseUrl, emailSent }
  const [showModal, setShowModal] = useState(false);

  async function handleSubmit(e) {
    e.preventDefault();
    let data = null;
    try {
      const res = await fetch(`${API_BASE}/api/download`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, type: platform }),
      });
      data = await res.json().catch(() => null);
    } catch {
      // Backend unreachable — fail silently, still confirm to the user below.
    } finally {
      setSubmittedEmail(email);
      setSubmitted(true);
      setResult(data);
      setShowModal(true);
    }
  }

  function closeModal() {
    setShowModal(false);
  }

  return (
    <section id="download">
      <div className="wrap">
        <Reveal className="form-card">
          <div className="form-grid">
            <div className="form-copy">
              <span className="eyebrow">Download</span>
              <h2>Get Voice2Prompt</h2>
              <p>
                Builds are hand-packaged right now, so we email the current download link
                instead of hosting a public button — that way you always get a working link
                plus a heads-up when a new build replaces it. Grab the source anytime from{" "}
                <a
                  href={GITHUB_RELEASES_URL}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="link-accent"
                  onClick={() => trackClick("github", "download_section")}
                >
                  GitHub Releases<ExternalIcon />
                </a>
                .
              </p>
            </div>

            <div>
              <div className="platform-toggle">
                <button
                  type="button"
                  className={`platform-btn${platform === "mac" ? " active" : ""}`}
                  onClick={() => setPlatform("mac")}
                >
                  <svg viewBox="0 0 24 24" fill="currentColor">
                    <path d="M16.365 1.43c0 1.14-.493 2.27-1.177 3.08-.744.9-1.99 1.57-2.987 1.57-.12 0-.23-.02-.3-.03-.01-.06-.04-.22-.04-.39 0-1.15.572-2.27 1.206-2.98.804-.94 2.142-1.64 3.248-1.68.03.13.05.28.05.43zm4.565 15.71c-.03.07-.463 1.58-1.518 3.12-.945 1.34-1.94 2.71-3.43 2.71-1.517 0-1.9-.88-3.63-.88-1.698 0-2.302.91-3.67.91-1.377 0-2.332-1.26-3.4-2.8-1.233-1.8-2.235-4.6-2.235-7.26 0-4.26 2.77-6.52 5.5-6.52 1.4 0 2.57.92 3.45.92.84 0 2.16-.98 3.76-.98.61 0 2.81.06 4.26 2.13-.11.07-2.54 1.48-2.51 4.42.03 3.51 3.08 4.68 3.11 4.7-.02.07-.49 1.66-1.57 3.28z" />
                  </svg>
                  macOS
                </button>
                <button
                  type="button"
                  className="platform-btn"
                  disabled
                  onClick={() => setPlatform("windows")}
                >
                  <svg viewBox="0 0 24 24" fill="currentColor">
                    <path d="M3 5.5L10 4.5V11.5H3V5.5M11 4.35L21 3V11.5H11V4.35M3 12.5H10V19.5L3 18.5V12.5M11 12.5H21V21L11 19.65V12.5Z" />
                  </svg>
                  Windows <span className="soon">Phase 4</span>
                </button>
              </div>

              {!submitted ? (
                <form onSubmit={handleSubmit}>
                  <p className="hp-field">
                    <label>
                      Don&apos;t fill this out: <input name="bot-field" tabIndex={-1} autoComplete="off" />
                    </label>
                  </p>
                  <div className="field-row">
                    <input
                      type="email"
                      name="email"
                      placeholder="you@example.com"
                      required
                      value={email}
                      onChange={(e) => setEmail(e.target.value)}
                    />
                    <button type="submit" className="btn btn-primary">
                      Email me the link
                    </button>
                  </div>
                  <p className="fine-print">
                    Required — we send the current build link straight to your inbox. We&apos;ll
                    also let you know about major new releases. No spam, unsubscribe anytime.
                  </p>
                  <p className="fine-print">
                    macOS will ask for a couple of permissions on first launch (and block the
                    app once, since it isn&apos;t notarized yet) — see the{" "}
                    <a href="#install-steps">install &amp; permissions steps</a> below
                    before you open it.
                  </p>
                </form>
              ) : (
                <div className="success-box show">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <path d="M20 6L9 17l-5-5" />
                  </svg>
                  <p>
                    <b>Check your inbox.</b> We&apos;ve sent the {platform === "mac" ? "macOS" : "Windows"}{" "}
                    download link to {submittedEmail}.
                  </p>
                </div>
              )}
            </div>
          </div>
        </Reveal>
      </div>

      {showModal && (
        <div className="dl-modal-overlay" onClick={closeModal}>
          <div className="dl-modal" onClick={(e) => e.stopPropagation()}>
            <button type="button" className="dl-modal-close" onClick={closeModal} aria-label="Close">
              &times;
            </button>
            <div className="dl-modal-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M20 6L9 17l-5-5" />
              </svg>
            </div>
            <h3>{result?.emailSent === false ? "Request received" : "Email sent successfully"}</h3>
            <p>
              {result?.emailSent === false
                ? `We saved your request, but the confirmation email to ${submittedEmail} may be delayed. You can grab the build directly below in the meantime.`
                : `We've sent the ${platform === "mac" ? "macOS" : "Windows"} download link to ${submittedEmail}. Didn't get it? Check spam, or use the links below.`}
            </p>
            <div className="dl-modal-actions">
              <a
                href={result?.releaseUrl || GITHUB_RELEASES_URL}
                target="_blank"
                rel="noopener noreferrer"
                className="btn btn-primary btn-block"
                onClick={() => trackClick("github", "download_popup")}
              >
                Go to Release<ExternalIcon />
              </a>
              <a
                href={result?.downloadUrl || GITHUB_RELEASES_URL}
                target="_blank"
                rel="noopener noreferrer"
                className="btn btn-ghost btn-block"
                onClick={() => trackClick("download", "download_popup")}
              >
                Download Anyway<ExternalIcon />
              </a>
            </div>
          </div>
        </div>
      )}
    </section>
  );
}
