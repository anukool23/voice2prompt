import Reveal from "./Reveal";

const NODES = [
  {
    title: "Hotkey",
    sub: "Fn / Carbon chord",
    icon: (
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
        <rect x="4" y="8" width="16" height="12" rx="2" />
        <path d="M8 8V6a4 4 0 0 1 8 0v2" />
      </svg>
    ),
  },
  {
    title: "Mic capture",
    sub: "malgo · 16kHz mono",
    icon: (
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
        <path d="M12 1a4 4 0 0 0-4 4v6a4 4 0 0 0 8 0V5a4 4 0 0 0-4-4z" />
        <path d="M19 10v1a7 7 0 0 1-14 0v-1" />
      </svg>
    ),
  },
  {
    title: "whisper-server",
    sub: "POST /inference",
    icon: (
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
        <path d="M9 18V5l12-2v13" />
        <circle cx="6" cy="18" r="3" />
        <circle cx="18" cy="16" r="3" />
      </svg>
    ),
  },
  {
    title: "Ollama",
    sub: "POST /api/chat",
    icon: (
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
        <path d="M12 2l2.4 7.2L22 12l-7.6 2.8L12 22l-2.4-7.2L2 12l7.6-2.8z" />
      </svg>
    ),
  },
  {
    title: "Inject",
    sub: "AX insert / ⌘V",
    icon: (
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
        <path d="M4 17l6-6-6-6M12 19h8" />
      </svg>
    ),
  },
];

export default function Architecture() {
  return (
    <section id="how-it-works">
      <div className="wrap">
        <Reveal className="section-head">
          <span className="eyebrow">Architecture</span>
          <h2>One pipeline, all of it on localhost</h2>
          <p>Every stage — capture, transcription, cleanup, injection — talks over 127.0.0.1. Nothing you say is sent anywhere.</p>
        </Reveal>

        <Reveal className="arch">
          <div className="arch-flow">
            {NODES.map((node, i) => (
              <div key={node.title} style={{ display: "contents" }}>
                <div className="arch-node">
                  <div className="n-icon">{node.icon}</div>
                  <h4>{node.title}</h4>
                  <span>{node.sub}</span>
                </div>
                {i < NODES.length - 1 && <div className="arch-arrow">→</div>}
              </div>
            ))}
          </div>
          <div className="arch-note">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M12 2l8 4v6c0 5-3.5 8.5-8 10-4.5-1.5-8-5-8-10V6l8-4z" />
            </svg>
            <p>
              <b>Both the STT and LLM engines are local HTTP sidecars on 127.0.0.1.</b> The Go
              core (<code className="code-mono">cmd/voice2prompt</code>) shares config and engine
              code with the Wails desktop app, so the CLI and the menu-bar app behave identically.
            </p>
          </div>
        </Reveal>
      </div>
    </section>
  );
}
