"use client";

import { useEffect, useRef } from "react";
import { trackPageview } from "@/lib/analytics";

/**
 * Fires one "pageview" analytics event per site visit. Mounted once in the
 * root layout so it covers every route. Renders nothing.
 */
export default function Analytics() {
  const firedRef = useRef(false);

  useEffect(() => {
    if (firedRef.current) return;
    firedRef.current = true;
    trackPageview();
  }, []);

  return null;
}
