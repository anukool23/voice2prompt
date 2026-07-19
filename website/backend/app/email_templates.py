"""
HTML email bodies. Kept out of the routers so the templates can be read/edited
on their own. Email clients don't reliably support modern CSS, so everything
here is inline styles + a table-based wrapper for max-width centering.
"""

_FONT = (
    "-apple-system, BlinkMacSystemFont, 'Segoe UI', Helvetica, Arial, sans-serif"
)
_MONO = "'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace"

_ACCENT = "#7C3AED"  # close to the site's --violet, safe as a flat email color
_TEXT = "#1c1c1e"
_DIM = "#6b6b76"
_BORDER = "#e5e5ea"
_CODE_BG = "#f4f3f1"

PLATFORM_LABELS = {"mac": "macOS", "windows": "Windows"}


def _wrapper(body_html: str) -> str:
    return f"""
<div style="background:#f5f5f7; padding:32px 16px; font-family:{_FONT};">
  <table role="presentation" width="100%" cellpadding="0" cellspacing="0"
         style="max-width:560px; margin:0 auto; background:#ffffff; border-radius:16px;
                overflow:hidden; border:1px solid {_BORDER};">
    <tr>
      <td style="padding:32px 32px 8px;">
        <div style="font-size:15px; font-weight:700; color:{_TEXT};">Voice2Prompt</div>
      </td>
    </tr>
    <tr>
      <td style="padding:8px 32px 32px; color:{_TEXT}; font-size:14.5px; line-height:1.6;">
        {body_html}
      </td>
    </tr>
  </table>
  <p style="max-width:560px; margin:16px auto 0; text-align:center; color:{_DIM}; font-size:12px;">
    You're getting this because you requested a Voice2Prompt download.
  </p>
</div>
""".strip()


def _button(url: str, label: str) -> str:
    return f"""
<a href="{url}" target="_blank" rel="noopener"
   style="display:inline-block; background:{_ACCENT}; color:#ffffff; text-decoration:none;
          font-weight:600; font-size:14.5px; padding:12px 22px; border-radius:999px; margin:4px 0 20px;">
  Download for {label}
</a>
""".strip()


def _step(number: int, title: str, body_html: str) -> str:
    return f"""
<tr>
  <td style="width:32px; vertical-align:top; padding:14px 12px 14px 0;">
    <div style="width:26px; height:26px; border-radius:50%; background:{_ACCENT}; color:#fff;
                font-size:12.5px; font-weight:700; text-align:center; line-height:26px;">{number}</div>
  </td>
  <td style="vertical-align:top; padding:14px 0; border-top:1px solid {_BORDER};">
    <div style="font-weight:600; font-size:14px; margin-bottom:4px;">{title}</div>
    <div style="color:{_DIM}; font-size:13px; line-height:1.55;">{body_html}</div>
  </td>
</tr>
""".strip()


def _code(text: str) -> str:
    return (
        f'<span style="font-family:{_MONO}; background:{_CODE_BG}; padding:2px 6px; '
        f'border-radius:5px; font-size:12.5px; color:{_TEXT};">{text}</span>'
    )


def _mac_install_steps_html() -> str:
    steps = [
        (
            "Open the DMG and install",
            f"Double-click {_code('Voice2Prompt.dmg')}, then drag {_code('Voice2Prompt.app')} "
            f"onto the {_code('Applications')} shortcut in the same window.",
        ),
        (
            "Bypass Gatekeeper (one time only)",
            "macOS will say it can't verify the developer — that's expected for an "
            "ad-hoc signed, open-source app. Right-click Voice2Prompt.app in "
            "Applications → <b>Open</b> → <b>Open</b> again in the dialog. Or: "
            "<b>System Settings → Privacy &amp; Security</b> → click "
            "<b>Open Anyway</b> next to the Voice2Prompt notice → confirm. Still "
            f"blocked? Run {_code('xattr -dr com.apple.quarantine /Applications/Voice2Prompt.app')} in Terminal.",
        ),
        (
            "Grant permissions",
            "<b>Microphone</b> — allow when prompted on first recording. "
            "<b>Accessibility</b> (required) — System Settings → Privacy &amp; "
            "Security → Accessibility → toggle Voice2Prompt on, then quit and "
            "reopen the app. <b>Input Monitoring</b> — only if you use the "
            "default Fn-key trigger (also set Keyboard → \"Press \U0001f310 key to\" "
            "→ Do Nothing). Skip it entirely if you pick a hotkey chord instead.<br><br>"
            "<b>Six triggers to choose from</b> in Settings → Trigger: "
            + _code("Fn / \U0001f310")
            + " (hold, double-tap to lock — needs Input Monitoring), or "
            + _code("Ctrl+Option+Space") + ", " + _code("Ctrl+Shift+Space") + ", "
            + _code("Cmd+Option+Space") + ", " + _code("F8") + ", " + _code("F9")
            + " (hold to talk, no extra permission).",
        ),
        (
            "Start dictating",
            "Hold your trigger, speak, release — the transcript types itself "
            "wherever your cursor is.",
        ),
    ]
    rows = "".join(_step(i, title, body) for i, (title, body) in enumerate(steps, start=1))
    return f"""
<div style="margin-top:8px; padding-top:20px; border-top:1px solid {_BORDER};">
  <div style="font-weight:700; font-size:14.5px; margin-bottom:4px;">First launch on macOS</div>
  <div style="color:{_DIM}; font-size:13px; margin-bottom:8px;">
    A quick heads-up so the security prompt doesn't catch you off guard:
  </div>
  <table role="presentation" width="100%" cellpadding="0" cellspacing="0">{rows}</table>
  <div style="color:{_DIM}; font-size:12px; margin-top:14px;">
    Every permission is scoped to this one app and can be revoked anytime from the
    same Privacy &amp; Security screens. Full write-up:
    <a href="https://github.com/anukool23/voice2prompt#installing-the-app-macos"
       style="color:{_ACCENT};">github.com/anukool23/voice2prompt</a>.
  </div>
</div>
""".strip()


def download_email_html(*, platform: str, download_url: str) -> str:
    label = PLATFORM_LABELS.get(platform, platform)
    body = f"""
<p style="margin:0 0 4px;">Hi,</p>
<p style="margin:0 0 4px;">Here's your {label} download link for Voice2Prompt:</p>
{_button(download_url, label)}
<p style="margin:0 0 4px; color:{_DIM}; font-size:13px;">
  This link always points to the latest release — feel free to re-download from it later.
</p>
{_mac_install_steps_html() if platform == "mac" else ""}
<p style="margin:20px 0 0;">— Voice2Prompt</p>
""".strip()
    return _wrapper(body)
