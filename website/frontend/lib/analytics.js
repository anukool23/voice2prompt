const API_BASE = process.env.NEXT_PUBLIC_API_URL || "";

/**
 * Fire-and-forget event log. Never awaited by callers, never blocks navigation.
 * Uses navigator.sendBeacon when available (survives the tab navigating away,
 * e.g. clicking an external link), falling back to a keepalive fetch.
 *
 *   trackEvent({ event: "pageview" })
 *   trackEvent({ event: "click", target: "github", location: "header" })
 */
export function trackEvent({ event, target, location } = {}) {
  if (typeof window === "undefined") return;

  const payload = JSON.stringify({
    event,
    target,
    location,
    page: window.location.pathname,
    referrer: document.referrer || undefined,
  });

  const url = `${API_BASE}/api/analytics`;

  try {
    if (navigator.sendBeacon) {
      const blob = new Blob([payload], { type: "application/json" });
      const sent = navigator.sendBeacon(url, blob);
      if (sent) return;
    }
  } catch {
    // fall through to fetch
  }

  fetch(url, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: payload,
    keepalive: true,
  }).catch(() => {
    // Analytics failing should never affect the user's experience.
  });
}

export function trackPageview() {
  trackEvent({ event: "pageview" });
}

export function trackClick(target, location) {
  trackEvent({ event: "click", target, location });
}
