import Reveal from "./Reveal";

export default function InstallGuide() {
  return (
    <section id="install-steps">
      <div className="wrap">
        <Reveal className="section-head">
          <span className="eyebrow">First launch (macOS)</span>
          <h2>Getting past Gatekeeper, step by step</h2>
          <p>
            Voice2Prompt is ad-hoc signed, not notarized yet (that needs a paid Apple
            Developer account) — so macOS blocks it once and asks for a few
            permissions. None of this is unusual for an open-source app; here&apos;s
            exactly what to click.
          </p>
        </Reveal>

        <Reveal className="install-card">
          <div className="install-grid">

            <div className="install-step">
              <div className="install-step-num">1</div>
              <div>
                <h4>Open the DMG and install</h4>
                <p>
                  Double-click <code>Voice2Prompt.dmg</code>, then drag{" "}
                  <code>Voice2Prompt.app</code> onto the <code>Applications</code>{" "}
                  shortcut in the same window. Eject the DMG once it&apos;s copied.
                </p>
              </div>
            </div>

            <div className="install-step">
              <div className="install-step-num">2</div>
              <div>
                <h4>Bypass Gatekeeper (one time only)</h4>
                <p>Opening it normally shows this — it&apos;s expected, not a sign anything&apos;s wrong:</p>

                <div className="install-screens">
                  <figure>
                    <div className="mac-mock">
                      <div className="mac-mock-icon">🛡️</div>
                      <div className="mac-mock-title">&quot;Voice2Prompt&quot; Not Opened</div>
                      <div className="mac-mock-body">
                        Apple could not verify &quot;Voice2Prompt&quot; is free of malware
                        that may harm your Mac or compromise your privacy.
                      </div>
                      <div className="mac-mock-buttons">
                        <span className="mac-btn">Done</span>
                        <span className="mac-btn mac-btn-primary">Move to Bin</span>
                      </div>
                    </div>
                    <figcaption>What you&apos;ll see on a plain double-click</figcaption>
                  </figure>
                </div>

                <p style={{ marginTop: 14 }}>
                  Click <b className="install-highlight">Done</b> (not Move to Bin!) —
                  then bypass it properly, either way:
                </p>
                <div className="sub">
                  <p>
                    <b className="install-highlight">Right-click</b> (or Control-click){" "}
                    <code>Voice2Prompt.app</code> in Applications →{" "}
                    <b className="install-highlight">Open</b> → click{" "}
                    <b className="install-highlight">Open</b> again in the dialog that
                    appears.
                  </p>
                  <p>
                    <b className="install-highlight">— or —</b> go to{" "}
                    <b className="install-highlight">System Settings → Privacy &amp; Security</b>,
                    scroll to the security notice, and click{" "}
                    <b className="install-highlight">Open Anyway</b>:
                  </p>
                </div>

                <div className="install-screens">
                  <figure>
                    <div className="mac-mock mac-mock-wide">
                      <div className="mac-mock-row">
                        <span>Allow applications from</span>
                        <span className="mac-mock-pill">App Store &amp; Known Developers</span>
                      </div>
                      <div className="mac-mock-row">
                        <span>&quot;Voice2Prompt&quot; was blocked to protect your Mac.</span>
                        <span className="mac-btn">Open Anyway</span>
                      </div>
                    </div>
                    <figcaption>System Settings → Privacy &amp; Security</figcaption>
                  </figure>
                  <figure>
                    <div className="mac-mock">
                      <div className="mac-mock-icon">🛡️</div>
                      <div className="mac-mock-title">Open &quot;Voice2Prompt&quot;?</div>
                      <div className="mac-mock-body">
                        Apple is not able to verify that it is free from malware that
                        could harm your Mac or compromise your privacy.
                      </div>
                      <div className="mac-mock-buttons-stack">
                        <span className="mac-btn full">Move to Bin</span>
                        <span className="mac-btn mac-btn-primary full">Open Anyway</span>
                        <span className="mac-btn full">Done</span>
                      </div>
                    </div>
                    <figcaption>Confirm once more — click Open Anyway</figcaption>
                  </figure>
                </div>

                <p style={{ marginTop: 14 }}>
                  Still getting blocked (common right after a browser download)? Clear
                  the quarantine flag from Terminal:
                </p>
                <div className="sub">
                  <p>
                    <code>xattr -dr com.apple.quarantine /Applications/Voice2Prompt.app</code>
                  </p>
                </div>
              </div>
            </div>

            <div className="install-step">
              <div className="install-step-num">3</div>
              <div>
                <h4>Grant permissions</h4>
                <p>
                  <b className="install-highlight">Microphone</b> — macOS prompts
                  automatically the first time you record. Click Allow.
                </p>
                <p>
                  <b className="install-highlight">Accessibility</b> (required) —{" "}
                  <code>System Settings → Privacy &amp; Security → Accessibility</code>,
                  toggle Voice2Prompt on. Quit and reopen the app afterward — permission
                  changes only apply on relaunch.
                </p>
                <p>
                  <b className="install-highlight">Input Monitoring</b> (only if using
                  the default Fn-key trigger) —{" "}
                  <code>System Settings → Privacy &amp; Security → Input Monitoring</code>,
                  toggle Voice2Prompt on. Also set{" "}
                  <code>Keyboard → &quot;Press 🌐 key to&quot;</code> → Do Nothing. Picked
                  a hotkey chord instead? Skip this one entirely — none of the five need it.
                </p>
                <div className="sub">
                  <p>
                    <b className="install-highlight">Six triggers to choose from</b> in
                    Settings → Trigger:
                  </p>
                  <p>
                    <code>Fn / 🌐</code> — hold to talk, double-tap to lock into
                    hands-free (needs Input Monitoring)
                  </p>
                  <p>
                    <code>Ctrl+Option+Space</code> · <code>Ctrl+Shift+Space</code> ·{" "}
                    <code>Cmd+Option+Space</code> · <code>F8</code> · <code>F9</code> —
                    hold to talk, no extra permission
                  </p>
                </div>
              </div>
            </div>

            <div className="install-step">
              <div className="install-step-num">4</div>
              <div>
                <h4>Start dictating</h4>
                <p>
                  Hold your trigger, speak, release — the transcript types itself
                  wherever your cursor is.
                </p>
              </div>
            </div>

          </div>

          <div className="arch-note">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M12 2l8 4v6c0 5-3.5 8.5-8 10-4.5-1.5-8-5-8-10V6l8-4z" />
            </svg>
            <p>
              <b>Every permission is scoped to this one app</b> and can be revoked
              anytime from the same Privacy &amp; Security screens. Voice2Prompt is
              fully open source, so you can read exactly what each permission is used
              for in <code className="code-mono">internal/inject/</code>,{" "}
              <code className="code-mono">internal/audio/</code>, and{" "}
              <code className="code-mono">internal/hotkey/</code>.
            </p>
          </div>
        </Reveal>
      </div>
    </section>
  );
}
