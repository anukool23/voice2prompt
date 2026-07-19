"use client";

import { useEffect, useRef, useState } from "react";

const DEMO_TEXT = "hello world this is a test";
const TRIGGER_KEYS = ["Ctrl+Option+Space", "Ctrl+Shift+Space", "Cmd+Option+Space", "F8", "F9", "Fn"];

export default function TerminalDemo() {
  const [typed, setTyped] = useState("");
  const [showResult, setShowResult] = useState(false);
  const [bars, setBars] = useState([]);
  const [triggerIdx, setTriggerIdx] = useState(0);

  // Waveform bars — randomized once on mount, same as the original inline script.
  useEffect(() => {
    const generated = Array.from({ length: 28 }, () => ({
      height: 10 + Math.random() * 24,
      delay: (Math.random() * 1.1).toFixed(2),
      duration: (0.7 + Math.random() * 0.8).toFixed(2),
    }));
    setBars(generated);
  }, []);

  // Cycle through all six triggers, one per second.
  useEffect(() => {
    const iv = setInterval(() => {
      setTriggerIdx((i) => (i + 1) % TRIGGER_KEYS.length);
    }, 1000);
    return () => clearInterval(iv);
  }, []);

  // Typing loop
  useEffect(() => {
    let interval;
    let restartTimeout;
    let resultTimeout;
    const startTimeout = setTimeout(function run() {
      let i = 0;
      setTyped("");
      setShowResult(false);
      interval = setInterval(() => {
        i++;
        setTyped(DEMO_TEXT.slice(0, i));
        if (i >= DEMO_TEXT.length) {
          clearInterval(interval);
          resultTimeout = setTimeout(() => setShowResult(true), 300);
          restartTimeout = setTimeout(run, 4200);
        }
      }, 55);
    }, 700);

    return () => {
      clearTimeout(startTimeout);
      clearTimeout(restartTimeout);
      clearTimeout(resultTimeout);
      clearInterval(interval);
    };
  }, []);

  return (
    <div className="terminal">
      <div className="terminal-bar">
        <span></span>
        <span></span>
        <span></span>
        <span className="terminal-title">voice2prompt — interactive</span>
      </div>
      <div className="terminal-body">
        <div className="prompt-line">
          <span className="arrow">➜</span> ./bin/voice2prompt
        </div>
        <div className="prompt-line mt-tight">
          Hold <b className="hl-white">{TRIGGER_KEYS[triggerIdx]}</b>, speak, release…
        </div>
        <div className="waveform">
          {bars.map((bar, i) => (
            <i
              key={i}
              style={{
                height: `${bar.height}px`,
                animationDelay: `${bar.delay}s`,
                animationDuration: `${bar.duration}s`,
              }}
            />
          ))}
        </div>
        <div className="typed-line">
          📝 &quot;<span>{typed}</span>
          <span className="typed-cursor"></span>&quot;
        </div>
        <div className={`result-line${showResult ? " show" : ""}`}>
          audio 2.1s | infer 190ms | TOTAL <b>260ms</b> | <span className="via">via accessibility</span> ✅ under 800ms
        </div>
      </div>
    </div>
  );
}
