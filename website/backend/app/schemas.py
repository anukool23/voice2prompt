from typing import Literal, Optional

from pydantic import BaseModel, EmailStr


class NewsletterIn(BaseModel):
    email: EmailStr


class DownloadIn(BaseModel):
    email: EmailStr
    type: Literal["mac", "windows"]


class AnalyticsIn(BaseModel):
    event: Literal["pageview", "click"]
    # For "click" events: what was clicked, e.g. "download", "github", "developer".
    target: Optional[str] = None
    # Where on the site it was clicked, e.g. "header", "footer", "hero", "download_section".
    location: Optional[str] = None
    page: Optional[str] = None
    referrer: Optional[str] = None
