"use client";

import { useEffect, useRef, useState } from "react";

/**
 * Wraps a section in the same scroll-reveal fade/slide-up behavior as the
 * original static site's IntersectionObserver script.
 */
export default function Reveal({ as: Tag = "div", className = "", style, children }) {
  const ref = useRef(null);
  const [inView, setInView] = useState(false);

  useEffect(() => {
    const el = ref.current;
    if (!el) return;

    // Safety net: if IntersectionObserver is unavailable, misbehaves, or simply
    // never fires (some embedded/preview contexts don't report intersections
    // reliably), force the section visible after a short delay instead of
    // leaving it stuck at opacity:0 forever.
    const fallback = setTimeout(() => setInView(true), 1200);

    if (typeof IntersectionObserver === "undefined") {
      return () => clearTimeout(fallback);
    }

    const io = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            setInView(true);
            clearTimeout(fallback);
            io.unobserve(entry.target);
          }
        });
      },
      { threshold: 0.12 }
    );
    io.observe(el);
    return () => {
      io.disconnect();
      clearTimeout(fallback);
    };
  }, []);

  return (
    <Tag ref={ref} className={`reveal${inView ? " in" : ""}${className ? ` ${className}` : ""}`} style={style}>
      {children}
    </Tag>
  );
}
