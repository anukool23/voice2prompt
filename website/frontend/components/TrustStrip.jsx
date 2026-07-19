export default function TrustStrip() {
  return (
    <section className="trust">
      <div className="wrap trust-grid">
        <div className="trust-item">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <path d="M12 2l8 4v6c0 5-3.5 8.5-8 10-4.5-1.5-8-5-8-10V6l8-4z" />
          </svg>
          <span>
            <b>100% on-device</b> — nothing leaves localhost
          </span>
        </div>
        <div className="trust-item">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <path d="M4 12l6 6L20 6" />
          </svg>
          <span>
            <b>whisper.cpp + Ollama</b> — no per-word API cost
          </span>
        </div>
        <div className="trust-item">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <rect x="3" y="4" width="18" height="16" rx="2" />
            <path d="M3 9h18M8 2v4M16 2v4" />
          </svg>
          <span>
            <b>Go engine</b> — macOS today, Windows in progress
          </span>
        </div>
        <div className="trust-item">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <circle cx="12" cy="12" r="9" />
            <path d="M12 7v5l3 3" />
          </svg>
          <span>
            <b>Open source</b> — MIT, built in public
          </span>
        </div>
      </div>
    </section>
  );
}
