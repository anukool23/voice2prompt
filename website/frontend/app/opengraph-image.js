import { ImageResponse } from "next/og";

export const runtime = "edge";
export const alt = "Voice2Prompt — Open-source, on-device voice dictation";
export const size = { width: 1200, height: 630 };
export const contentType = "image/png";

export default async function OgImage() {
  return new ImageResponse(
    (
      <div
        style={{
          width: "100%",
          height: "100%",
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
          justifyContent: "center",
          background: "#07070c",
          backgroundImage:
            "radial-gradient(circle at 20% 20%, rgba(139,92,246,0.35), transparent 45%), radial-gradient(circle at 85% 15%, rgba(37,99,235,0.35), transparent 45%)",
          fontFamily: "sans-serif",
        }}
      >
        <div
          style={{
            display: "flex",
            alignItems: "center",
            gap: 24,
            marginBottom: 36,
          }}
        >
          <div
            style={{
              width: 96,
              height: 96,
              borderRadius: 24,
              background: "linear-gradient(120deg, #8B5CF6, #2563EB)",
              display: "flex",
            }}
          />
          <div style={{ fontSize: 56, fontWeight: 800, color: "#ececf3", display: "flex" }}>
            Voice2Prompt
          </div>
        </div>
        <div
          style={{
            fontSize: 30,
            color: "#9c9cb0",
            maxWidth: 820,
            textAlign: "center",
            display: "flex",
          }}
        >
          Open-source, on-device voice dictation — no cloud, no per-word cost
        </div>
      </div>
    ),
    { ...size }
  );
}
