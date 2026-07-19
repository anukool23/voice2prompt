export default function manifest() {
  return {
    name: "Voice2Prompt",
    short_name: "Voice2Prompt",
    description:
      "Free, open-source, on-device voice dictation tool for macOS and Windows.",
    start_url: "/",
    display: "standalone",
    background_color: "#07070c",
    theme_color: "#07070c",
    icons: [
      {
        src: "/icon.png",
        sizes: "1024x1024",
        type: "image/png",
      },
    ],
  };
}
