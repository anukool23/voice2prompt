# Voice2Prompt website — backend (FastAPI)

Three endpoints backing the frontend, on Supabase (storage) + Resend (email),
deployed to Vercel as a Python serverless function.

## Endpoints

- `POST /api/newsletter` — `{ "email": "..." }` → upserts into `subscribers`
  as `{ email, created_at }`. Subscribing twice with the same email is a no-op.
- `POST /api/download` — `{ "email": "...", "type": "mac" | "windows" }` →
  records the request in `download_requests`, then emails the current download
  link for that platform via Resend.
- `POST /api/analytics` — `{ "event": "pageview" | "click", "target"?, "location"?, "page"?, "referrer"? }`
  → inserts into `analytics_events`. Fire one `pageview` per site visit, and one
  `click` whenever someone clicks a Download / GitHub / Developer link anywhere
  on the site (`target` says which, `location` says where — header, footer, hero, etc.).
- `GET /api/health` — `{ "ok": true }`, for sanity checks.

## One-time setup

1. **Supabase** — create a project, then run `schema.sql` in the SQL editor
   (Project → SQL Editor → New query → paste → Run). Grab `Project URL` and the
   **service role** key from Project Settings → API.
2. **Resend** — create an API key, and verify a sending domain/address (or use
   their `onboarding@resend.dev` sandbox sender while testing).
3. Copy `.env.example` to `.env` and fill in the values above.

## Run locally

```sh
python -m venv .venv && source .venv/bin/activate   # or your preferred env tool
pip install -r requirements.txt
uvicorn app.main:app --reload --port 8000
```

Then point the frontend at it: in `../frontend/.env.local`, set
`NEXT_PUBLIC_API_URL=http://localhost:8000`.

## Deploy to Vercel

Deploy this `backend/` folder as its **own** Vercel project (separate from the
Next.js frontend project):

```sh
cd backend
vercel        # first deploy, follow the prompts
vercel --prod
```

Vercel auto-detects `api/index.py` as a Python serverless function; `vercel.json`
rewrites every request to it, and FastAPI's own router (in `app/main.py`) takes
it from there. Add every variable from `.env.example` under the new project's
Settings → Environment Variables (for Production **and** Preview).

Once deployed, set `NEXT_PUBLIC_API_URL` on the **frontend's** Vercel project to
this backend's URL (e.g. `https://voice2prompt-api.vercel.app`), and redeploy
the frontend.

## Notes

- CORS is controlled by `CORS_ORIGINS` — `*` works for testing, but once both
  projects are live, set it to the frontend's real domain(s) so only your site
  can call this API from a browser.
- `download_requests` and `analytics_events` are append-only logs — nothing here
  currently reads them back. Query them directly in the Supabase dashboard, or
  build a small internal view later if you want charts.
- Row Level Security is enabled on all three tables with no policies, so the
  public `anon` key can't touch them — only the service role key used here can.
