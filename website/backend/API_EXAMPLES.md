# Voice2Prompt backend — cURL examples

Replace `$API` with your deployed URL (e.g. `https://voice2prompt-api.vercel.app`)
or `http://localhost:8000` when running locally via `uvicorn app.main:app --reload`.

```sh
export API="https://your-backend.vercel.app"
```

## Health check

```sh
curl "$API/api/health"
```
Response: `{"ok": true}`

## Newsletter — subscribe

```sh
curl -X POST "$API/api/newsletter" \
  -H "Content-Type: application/json" \
  -d '{"email": "you@example.com"}'
```
Response: `{"ok": true}`
Upserts into `subscribers`; subscribing the same email twice is a no-op, not an error.

## Download — request a link by email

```sh
# macOS
curl -X POST "$API/api/download" \
  -H "Content-Type: application/json" \
  -d '{"email": "you@example.com", "type": "mac"}'

# Windows
curl -X POST "$API/api/download" \
  -H "Content-Type: application/json" \
  -d '{"email": "you@example.com", "type": "windows"}'
```
Response: `{"ok": true, "downloadUrl": "https://..."}`
`type` must be exactly `"mac"` or `"windows"` (Pydantic will reject anything else
with a 422). Records the request in `download_requests` and emails the link via
Resend — for `mac`, the email also includes the Gatekeeper/permissions walkthrough.

## Analytics — pageview

```sh
curl -X POST "$API/api/analytics" \
  -H "Content-Type: application/json" \
  -d '{"event": "pageview", "page": "/", "referrer": "https://google.com"}'
```

## Analytics — click

```sh
curl -X POST "$API/api/analytics" \
  -H "Content-Type: application/json" \
  -d '{"event": "click", "target": "download", "location": "header"}'
```
`event` must be `"pageview"` or `"click"`. `target`/`location`/`page`/`referrer`
are all optional strings — `target` is typically `"download" | "github" | "developer"`,
`location` is typically `"header" | "footer" | "hero" | "download_section"`, etc.
(nothing is enforced server-side beyond that, so any string works).
Response for all three: `{"ok": true}`

## Error responses

- Invalid/missing email, bad `type`, or bad `event` → `422 Unprocessable Entity`
  with Pydantic's validation detail.
- Supabase write failure → `500 Internal Server Error`, `{"detail": "..."}`
- CORS: if calling from a browser and getting blocked, check `CORS_ORIGINS` on
  the backend's Vercel project matches your frontend's actual origin.
