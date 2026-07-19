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
    if not settings.supabase_url or not settings.supabase_service_role_key:
        raise RuntimeError(
            "SUPABASE_URL and SUPABASE_SERVICE_ROLE_KEY must be set in the environment"
        )
    return create_client(settings.supabase_url, settings.supabase_service_role_key)
