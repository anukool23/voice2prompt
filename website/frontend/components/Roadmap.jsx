import Reveal from "./Reveal";

const PHASES = [
  { phase: "P0", title: "Feasibility spike", desc: "Hotkey → record → whisper → paste, latency measured end-to-end.", status: "done" },
  { phase: "P1", title: "macOS text injection", desc: "Accessibility-API caret insert with clipboard fallback.", status: "done" },
  { phase: "P2", title: "Local LLM cleanup", desc: "Punctuation, filler removal, context-aware tone via Ollama.", status: "done" },
  { phase: "P3", title: "Wails settings app", desc: "Onboarding, permissions, hotkey + dictionary editor.", status: "progress" },
  { phase: "P4", title: "Windows port", desc: "UI Automation-based injection for Windows.", status: "planned" },
  { phase: "P5", title: "Command mode", desc: "Editing keystrokes + LLM rewrite of your current selection.", status: "progress" },
  { phase: "P6", title: "Packaging", desc: "Self-contained .app, menu-bar tray, launch-at-login.", status: "progress" },
];

const STATUS_LABEL = {
  done: { cls: "status-done", label: "✓ Shipped" },
  progress: { cls: "status-progress", label: "◐ In progress" },
  planned: { cls: "status-planned", label: "○ Planned" },
};

export default function Roadmap() {
  return (
    <section id="roadmap">
      <div className="wrap">
        <Reveal className="section-head">
          <span className="eyebrow">Roadmap</span>
          <h2>Built in the open, shipping phase by phase</h2>
          <p>This is early — and that's the point. Star the repo or subscribe below to catch new releases as each phase lands.</p>
        </Reveal>

        <Reveal className="roadmap">
          {PHASES.map((p) => {
            const status = STATUS_LABEL[p.status];
            return (
              <div className="rm-row" key={p.phase}>
                <span className="rm-phase">{p.phase}</span>
                <div className="rm-body">
                  <h4>{p.title}</h4>
                  <p>{p.desc}</p>
                </div>
                <span className={`status ${status.cls}`}>{status.label}</span>
              </div>
            );
          })}
        </Reveal>

        <Reveal className="roadmap-cta">
          <p>New phases ship as releases — the fastest way to know is the newsletter below, or watch the repo on GitHub.</p>
          <a href="#newsletter" className="btn btn-ghost">
            Get notified about new releases
          </a>
        </Reveal>
      </div>
    </section>
  );
}
