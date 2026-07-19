"""
Vercel Python entrypoint. Vercel's Python runtime detects the `app` ASGI
callable in this file and serves it directly — no Mangum/adapter needed.
`vercel.json` rewrites every request here so FastAPI's own router (defined in
app/main.py) handles /api/newsletter, /api/download, /api/analytics, /api/health.
"""

import sys
from pathlib import Path

# Make the project root (the parent of this api/ folder) importable regardless
# of Vercel's working directory, so `from app.main import app` always resolves.
sys.path.append(str(Path(__file__).resolve().parent.parent))

from app.main import app  # noqa: E402
