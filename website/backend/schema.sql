-- Voice2Prompt backend — Supabase schema.
-- Run this once in the Supabase SQL editor (Project -> SQL Editor -> New query).

create extension if not exists pgcrypto;

-- 1) Newsletter subscribers: { email, createdAt }
create table if not exists subscribers (
  id uuid primary key default gen_random_uuid(),
  email text not null unique,
  created_at timestamptz not null default now()
);

-- 2) Download requests: which platform, which email, which link was sent
create table if not exists download_requests (
  id uuid primary key default gen_random_uuid(),
  email text not null,
  platform text not null check (platform in ('mac', 'windows')),
  download_url text,
  created_at timestamptz not null default now()
);

-- 3) Analytics: pageviews + clicks on download/github/developer links anywhere
--    on the site (header, footer, hero, download section, ...)
create table if not exists analytics_events (
  id uuid primary key default gen_random_uuid(),
  event_type text not null check (event_type in ('pageview', 'click')),
  target text,        -- e.g. 'download', 'github', 'developer' (click events only)
  location text,       -- e.g. 'header', 'footer', 'hero', 'download_section'
  page text,           -- path the event happened on, e.g. '/'
  referrer text,
  user_agent text,
  created_at timestamptz not null default now()
);

create index if not exists idx_subscribers_email on subscribers (email);
create index if not exists idx_download_requests_email on download_requests (email);
create index if not exists idx_download_requests_created_at on download_requests (created_at);
create index if not exists idx_analytics_events_type on analytics_events (event_type);
create index if not exists idx_analytics_events_target on analytics_events (target);
create index if not exists idx_analytics_events_created_at on analytics_events (created_at);

-- Lock every table down from the anon/public API. The backend only ever talks
-- to Supabase using the service role key, which bypasses RLS entirely — so no
-- policies are defined (or needed) here. This just makes sure nothing is
-- readable/writable via the public anon key if it ever leaks.
alter table subscribers enable row level security;
alter table download_requests enable row level security;
alter table analytics_events enable row level security;
