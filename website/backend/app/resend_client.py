"""
Thin wrapper around the Resend SDK so the rest of the app doesn't import/
configure it directly.
"""

import resend

from app.config import settings


def send_email(*, to: str, subject: str, html: str) -> None:
    if not settings.resend_api_key:
        raise RuntimeError("RESEND_API_KEY must be set in the environment")

    resend.api_key = settings.resend_api_key
    resend.Emails.send(
        {
            "from": settings.from_email,
            "to": [to],
            "subject": subject,
            "html": html,
        }
    )
