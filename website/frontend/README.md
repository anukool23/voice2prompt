# Voice2Prompt website — frontend

Next.js (App Router) port of the marketing site, pixel-matched to `../../index.html`.

## Structure

- `app/` — routes, layout, SEO metadata (`layout.js`), `sitemap.js`, `robots.js`,
  `manifest.js`, and a dynamic `opengraph-image.js`.
- `components/` — one component per section (Header, Hero, Features, Architecture,
  Roadmap, DownloadSection, Newsletter, Footer) plus shared bits (`Logo`, `Reveal`,
  `TerminalDemo`, `BgFx`, `ExternalIcon`, `Analytics`).
- `app/globals.css` — the same CSS as the static site, unchanged, plus a few small
  utility classes that replace one-off inline `style="..."` attributes from the HTML
  version (JSX doesn't take raw style strings).
- `lib/analytics.js` — `trackEvent` / `trackPageview` / `trackClick` helpers that
  fire-and-forget events to the backend's `/api/analytics` (via `navigator.sendBeacon`,
  falling back to `fetch(..., { keepalive: true })`).

## Run locally

```sh
npm install
npm run dev      # http://localhost:3000
npm run build && npm run start   # production build
```

## Backend

`../backend` is a FastAPI app (Supabase + Resend) with three endpoints this frontend
calls, all under `NEXT_PUBLIC_API_URL` (see `.env.local.example` — copy it to
`.env.local` and point it at the backend running locally or deployed):

- `POST /api/newsletter` — `{ email }`, from `components/Newsletter.jsx`.
- `POST /api/download` — `{ email, type: "mac" | "windows" }`, from
  `components/DownloadSection.jsx`. If the request fails, the UI still shows the
  success state so a flaky network doesn't strand the user mid-flow.
- `POST /api/analytics` — `{ event: "pageview" | "click", target?, location?, page?, referrer? }`,
  fired by `components/Analytics.jsx` (once per visit) and by `onClick` handlers on
  every Download / GitHub / Developer link in `Header`, `Hero`, `Footer`, and
  `DownloadSection` (see `lib/analytics.js`).

## Before you deploy

- `app/layout.js`, `app/sitemap.js`, `app/robots.js` hardcode a placeholder
  `https://voice2prompt.app` — swap in your real domain.
- Replace the placeholder GitHub org (`github.com/anukool23/voice2prompt`) throughout
  `components/` with your actual repo URL.
