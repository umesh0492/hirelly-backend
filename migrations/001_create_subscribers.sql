-- Migration 001: Create Subscribers Table
CREATE TABLE IF NOT EXISTS public.subscribers (
    id TEXT PRIMARY KEY,
    email TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'Executive',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    timestamp TEXT
);

CREATE INDEX IF NOT EXISTS idx_subscribers_email ON public.subscribers(email);
CREATE INDEX IF NOT EXISTS idx_subscribers_created_at ON public.subscribers(created_at DESC);
