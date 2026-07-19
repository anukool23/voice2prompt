export default function Logo({ gradId = "bg" }) {
  return (
    <svg viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
      <defs>
        <linearGradient id={gradId} x1="0" y1="0" x2="1024" y2="1024" gradientUnits="userSpaceOnUse">
          <stop offset="0" stopColor="#8B5CF6" />
          <stop offset="1" stopColor="#2563EB" />
        </linearGradient>
      </defs>
      <rect x="0" y="0" width="1024" height="1024" rx="228" fill={`url(#${gradId})`} />
      <g fill="none" stroke="#ffffff" strokeWidth="36" strokeLinecap="round" strokeLinejoin="round">
        <rect x="440" y="196" width="144" height="252" rx="72" fill="#ffffff" stroke="none" />
        <path d="M 372 420 A 140 140 0 0 0 652 420" />
        <line x1="512" y1="560" x2="512" y2="612" />
      </g>
      <g fill="none" stroke="#ffffff" strokeWidth="36" strokeLinecap="round" strokeLinejoin="round">
        <polyline points="452,648 512,700 452,752" />
        <line x1="556" y1="752" x2="656" y2="752" />
      </g>
    </svg>
  );
}
