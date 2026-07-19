"""
Centralized environment/config access. Everything secret comes from the process
environment — locally via a `.env` file (see .env.example), and on Vercel via
Project Settings -> Environment Variables. Nothing here is ever hardcoded.
"""

import os
from functools import lru_cache

from dotenv import load_dotenv

# Loads variables from a local .env file into os.environ, if one exists.
# No-op on Vercel (no .env file there — env vars are injected directly by the
# platform), so this is safe to call unconditionally in both environments.
load_dotenv()


def _split_csv(value: str) -> list[str]:
    return [item.strip() for item in value.split(",") if item.strip()]


def _normalize_supabase_url(url: str) -> str:
    """supabase-py appends /rest/v1 itself, so SUPABASE_URL must be the bare
    project URL. Supabase's dashboard also shows a "REST API URL" that already
    has /rest/v1 baked in — if that got pasted in by mistake, strip it (and any
    trailing slash) so the client doesn't end up hitting .../rest/v1/rest/v1/...
    """
    return url.strip().rstrip("/").removesuffix("/rest/v1")


class Settings:
    # --- Supabase (server-side only; service role key must stay secret) ---
    supabase_url: str = _normalize_supabase_url(os.environ.get("SUPABASE_URL", ""))
    supabase_service_role_key: str = os.environ.get("SUPABASE_SERVICE_ROLE_KEY", "")

    # --- Resend ---
    resend_api_key: str = os.environ.get("RESEND_API_KEY", "")
    from_email: str = os.environ.get("FROM_EMAIL", "Voice2Prompt <onboarding@resend.dev>")

    # --- Download links, per platform ---
    download_url_mac: str = os.environ.get(
        "DOWNLOAD_URL_MAC", "https://github.com/anukool23/voice2prompt/releases"
    )
    download_url_windows: str = os.environ.get(
        "DOWNLOAD_URL_WINDOWS", "https://github.com/anukool23/voice2prompt/releases"
    )

    # --- CORS: comma-separated list of allowed origins, or "*" for all ---
    _cors_origins_raw: str = os.environ.get("CORS_ORIGINS", "*")

    @property
    def cors_origins(self) -> list[str]:
        if self._cors_origins_raw.strip() == "*":
            return ["*"]
        return _split_csv(self._cors_origins_raw)


@lru_cache
def get_settings() -> Settings:
    return Settings()


settings = get_settings()
