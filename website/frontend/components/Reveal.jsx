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
    const io = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            setInView(true);
            io.unobserve(entry.target);
          }
        });
      },
      { threshold: 0.12 }
    );
    io.observe(el);
    return () => io.disconnect();
  }, []);

  return (
    <Tag ref={ref} className={`reveal${inView ? " in" : ""}${className ? ` ${className}` : ""}`} style={style}>
      {children}
    </Tag>
  );
}
