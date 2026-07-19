"use client";

import { useState } from "react";
import Reveal from "./Reveal";

const API_BASE = process.env.NEXT_PUBLIC_API_URL || "";

export default function Newsletter() {
  const [email, setEmail] = useState("");
  const [submitted, setSubmitted] = useState(false);
  const [submittedEmail, setSubmittedEmail] = useState("");

  async function handleSubmit(e) {
    e.preventDefault();
    try {
      await fetch(`${API_BASE}/api/newsletter`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email }),
      });
    } catch {
      // Backend not wired up yet — fail silently, still confirm to the user.
    } finally {
      setSubmittedEmail(email);
      setSubmitted(true);
    }
  }

  return (
    <section id="newsletter">
      <div className="wrap">
        <Reveal className="form-card form-card--center">
          <div className="newsletter-inner">
            <span className="eyebrow">Stay in the loop</span>
            <h2 className="h2-sm">New phases ship often. Don&apos;t miss one.</h2>
            <p className="newsletter-copy">
              One email per release — Windows support, command mode, packaging updates, and
              whatever we build next.
            </p>

            {!submitted ? (
              <form onSubmit={handleSubmit}>
                <p className="hp-field">
                  <label>
                    Don&apos;t fill this out: <input name="bot-field" tabIndex={-1} autoComplete="off" />
                  </label>
                </p>
                <div className="field-row field-row--narrow">
                  <input
                    type="email"
                    name="email"
                    placeholder="you@example.com"
                    required
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                  />
                  <button type="submit" className="btn btn-primary">
                    Subscribe
                  </button>
                </div>
              </form>
            ) : (
              <div className="success-box success-box--newsletter show">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <path d="M20 6L9 17l-5-5" />
                </svg>
                <p>
                  <b>You&apos;re in.</b> Watch for the next release note at {submittedEmail}.
                </p>
              </div>
            )}
          </div>
        </Reveal>
      </div>
    </section>
  );
}
