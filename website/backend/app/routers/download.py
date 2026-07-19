from datetime import datetime, timezone

from fastapi import APIRouter, HTTPException

from app.config import settings
from app.email_templates import PLATFORM_LABELS, download_email_html
from app.resend_client import send_email
from app.schemas import DownloadIn
from app.supabase_client import get_supabase

router = APIRouter(tags=["download"])


@router.post("/download")
def request_download(payload: DownloadIn):
    """Record the request in Supabase and email the current download link via
    Resend. Email is mandatory (enforced by the frontend + the schema). For
    macOS, the email also includes the Gatekeeper/permissions walkthrough so
    the security prompt doesn't blindside anyone."""
    email = payload.email.lower()
    platform = payload.type
    download_url = (
        settings.download_url_mac if platform == "mac" else settings.download_url_windows
    )

    try:
        get_supabase().table("download_requests").insert(
            {
                "email": email,
                "platform": platform,
                "download_url": download_url,
                "created_at": datetime.now(timezone.utc).isoformat(),
            }
        ).execute()
    except Exception as exc:  # noqa: BLE001
        raise HTTPException(status_code=500, detail="Could not record download request") from exc

    label = PLATFORM_LABELS.get(platform, platform)
    try:
        send_email(
            to=email,
            subject=f"Your Voice2Prompt download link ({label})",
            html=download_email_html(platform=platform, download_url=download_url),
        )
    except Exception as exc:  # noqa: BLE001
        # The request is already saved above — don't fail the whole call just
        # because email delivery hiccuped. Surface it in logs for now.
        print(f"[download] email send failed for {email}: {exc}")

    return {"ok": True, "downloadUrl": download_url}
