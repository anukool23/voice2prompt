"""
Lazy, cached Supabase client. Uses the service role key so the backend can read
and write freely regardless of Row Level Security policies — this key must never
be exposed to the frontend/browser.
"""

from functools import lru_cache

from supabase import Client, create_client

from app.config import settings


@lru_cache
def get_supabase() -> Client:
    missing = [
        name
        for name, value in (
            ("SUPABASE_URL", settings.supabase_url),
            ("SUPABASE_SERVICE_ROLE_KEY", settings.supabase_service_role_key),
        )
        if not value
    ]
    if missing:
        raise RuntimeError(
            f"Missing required env var(s): {', '.join(missing)}. Check that "
            f"website/backend/.env exists and has real (non-placeholder) values — "
            f"see .env.example."
        )
    return create_client(settings.supabase_url, settings.supabase_service_role_key)
