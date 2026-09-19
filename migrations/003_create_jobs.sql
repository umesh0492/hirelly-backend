-- Migration 003: Create Jobs Table
CREATE TABLE IF NOT EXISTS public.jobs (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    company TEXT NOT NULL,
    icon TEXT DEFAULT 'fa-briefcase',
    location TEXT NOT NULL,
    type TEXT NOT NULL DEFAULT 'Full-time',
    salary TEXT,
    salary_val INTEGER DEFAULT 0,
    tags TEXT[] DEFAULT '{}',
    posted TEXT DEFAULT 'Recently',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_jobs_created_at ON public.jobs(created_at DESC);
