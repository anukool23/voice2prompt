from datetime import datetime, timezone

from fastapi import APIRouter, HTTPException, Request

from app.schemas import AnalyticsIn
from app.supabase_client import get_supabase

router = APIRouter(tags=["analytics"])


@router.post("/analytics")
def track_event(payload: AnalyticsIn, request: Request):
    """Fire-and-forget event log. Two shapes in practice:
    - {"event": "pageview", "page": "/", "referrer": "..."} — once per site visit.
    - {"event": "click", "target": "download" | "github" | "developer", "location": "header" | "footer" | "hero" | ...}
    """
    record = {
        "event_type": payload.event,
        "target": payload.target,
        "location": payload.location,
        "page": payload.page,
        "referrer": payload.referrer,
        "user_agent": request.headers.get("user-agent"),
        "created_at": datetime.now(timezone.utc).isoformat(),
    }

    try:
        get_supabase().table("analytics_events").insert(record).execute()
    except Exception as exc:  # noqa: BLE001
        print(f"[analytics] insert failed: {exc}")
        raise HTTPException(status_code=500, detail=f"Could not record event: {exc}") from exc

    return {"ok": True}
