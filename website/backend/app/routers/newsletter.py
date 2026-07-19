from datetime import datetime, timezone

from fastapi import APIRouter, HTTPException

from app.schemas import NewsletterIn
from app.supabase_client import get_supabase

router = APIRouter(tags=["newsletter"])


@router.post("/newsletter")
def subscribe_newsletter(payload: NewsletterIn):
    """Save { email, createdAt } to Supabase. Idempotent — subscribing twice with
    the same email is a no-op rather than an error."""
    email = payload.email.lower()

    try:
        get_supabase().table("subscribers").upsert(
            {"email": email, "created_at": datetime.now(timezone.utc).isoformat()},
            on_conflict="email",
            ignore_duplicates=True,
        ).execute()
    except Exception as exc:  # noqa: BLE001
        raise HTTPException(status_code=500, detail="Could not save subscriber") from exc

    return {"ok": True}
