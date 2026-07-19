from datetime import datetime, timezone

from fastapi import APIRouter, HTTPException

from app.config import settings
from app.resend_client import send_email
from app.schemas import DownloadIn
from app.supabase_client import get_supabase

router = APIRouter(tags=["download"])

PLATFORM_LABELS = {"mac": "macOS", "windows": "Windows"}


@router.post("/download")
def request_download(payload: DownloadIn):
    """Record the request in Supabase and email the current download link via
    Resend. Email is mandatory (enforced by the frontend + the schema)."""
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

    label = PLATFORM_LABELS[platform]
    try:
        send_email(
            to=email,
            subject=f"Your Voice2Prompt download link ({label})",
            html=(
                "<p>Hi,</p>"
                f"<p>Here's your {label} download link for Voice2Prompt:</p>"
                f'<p><a href="{download_url}">{download_url}</a></p>'
                "<p>This link always points to the latest release.</p>"
                "<p>— Voice2Prompt</p>"
            ),
        )
    except Exception as exc:  # noqa: BLE001
        # The request is already saved above — don't fail the whole call just
        # because email delivery hiccuped. Surface it in logs for now.
        print(f"[download] email send failed for {email}: {exc}")

    return {"ok": True, "downloadUrl": download_url}
