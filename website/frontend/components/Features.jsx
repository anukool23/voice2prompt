import Reveal from "./Reveal";

const FEATURES = [
  {
    title: "Local speech-to-text",
    desc: "A whisper.cpp sidecar runs over localhost HTTP, keeping the model resident so every utterance skips reload time. Sub-second latency, zero cloud round-trip.",
    icon: (
      <svg viewBox="0 0 24 24" fill="none" stroke="#fff" strokeWidth="2">
        <path d="M12 1a4 4 0 0 0-4 4v6a4 4 0 0 0 8 0V5a4 4 0 0 0-4-4z" />
        <path d="M19 10v1a7 7 0 0 1-14 0v-1M12 18v4" />
      </svg>
    ),
  },
  {
    title: "LLM cleanup",
    desc: "A local Ollama model fixes punctuation, drops filler words, and adapts tone per app — casual for Slack, professional for Mail and Docs. Falls back to raw transcription if Ollama's off.",
    icon: (
      <svg viewBox="0 0 24 24" fill="none" stroke="#fff" strokeWidth="2">
        <path d="M12 2l2.4 7.2L22 12l-7.6 2.8L12 22l-2.4-7.2L2 12l7.6-2.8z" />
      </svg>
    ),
  },
  {
    title: "Smart text injection",
    desc: "Inserts at the caret via macOS's Accessibility API, and automatically falls back to clipboard + ⌘V for apps — like Slack, Notion, VS Code — that reject AX insertion.",
    icon: (
      <svg viewBox="0 0 24 24" fill="none" stroke="#fff" strokeWidth="2">
        <path d="M4 17l6-6-6-6M12 19h8" />
      </svg>
    ),
  },
  {
    title: "Voice commands",
    desc: '"Select all," "make this formal," "shorten this" — rule-based editing runs instantly, and rewrite commands hand your current selection to the LLM and replace it in place.',
    icon: (
      <svg viewBox="0 0 24 24" fill="none" stroke="#fff" strokeWidth="2">
        <path d="M9 18V5l12-2v13M9 9l12-2M6 21a3 3 0 1 0 0-6 3 3 0 0 0 0 6zM21 19a3 3 0 1 0 0-6 3 3 0 0 0 0 6z" />
      </svg>
    ),
  },
  {
    title: "Six ways to trigger it",
    desc: "Hold the 🌐 Fn key for push-to-talk, double-tap to lock into hands-free — or bind one of five hotkey chords instead, no extra permission required. Pick whichever fits your workflow in Settings → Trigger.",
    keys: [
      { label: "Fn / 🌐 — hold, double-tap to lock", isFn: true },
      { label: "Ctrl+Option+Space" },
      { label: "Ctrl+Shift+Space" },
      { label: "Cmd+Option+Space" },
      { label: "F8" },
      { label: "F9" },
    ],
    icon: (
      <svg viewBox="0 0 24 24" fill="none" stroke="#fff" strokeWidth="2">
        <rect x="4" y="8" width="16" height="12" rx="2" />
        <path d="M8 8V6a4 4 0 0 1 8 0v2" />
      </svg>
    ),
  },
  {
    title: "Native menu-bar app",
    desc: "Self-contained, ad-hoc signed .app with a menu-bar tray (Open / Start-Stop / Quit) and launch-at-login — one permission grant, no child-process mess.",
    icon: (
      <svg viewBox="0 0 24 24" fill="none" stroke="#fff" strokeWidth="2">
        <rect x="3" y="3" width="18" height="18" rx="3" />
        <circle cx="8" cy="8" r="1.2" fill="#fff" stroke="none" />
        <circle cx="8" cy="12" r="1.2" fill="#fff" stroke="none" />
        <circle cx="8" cy="16" r="1.2" fill="#fff" stroke="none" />
      </svg>
    ),
  },
];

export default function Features() {
  return (
    <section id="features">
      <div className="wrap">
        <Reveal className="section-head">
          <span className="eyebrow">Features</span>
          <h2>Everything a premium cloud dictation app does — running on your machine</h2>
          <p>Voice2Prompt ships fast because every layer is a swappable local sidecar. New capabilities land every release.</p>
        </Reveal>
        <div className="feature-grid">
          {FEATURES.map((f) => (
            <Reveal as="div" className="card" key={f.title}>
              <div className="card-icon">{f.icon}</div>
              <h3>{f.title}</h3>
              <p>{f.desc}</p>
              {f.keys && (
                <div className="trigger-keys">
                  {f.keys.map((k) => (
                    <span key={k.label} className={k.isFn ? "is-fn" : undefined}>
                      {k.label}
                    </span>
                  ))}
                </div>
              )}
            </Reveal>
          ))}
        </div>
      </div>
    </section>
  );
}
