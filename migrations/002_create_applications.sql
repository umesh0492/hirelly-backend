-- Migration 002: Create Applications / Executive Mandates Table
CREATE TABLE IF NOT EXISTS public.applications (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL DEFAULT 'Executive Contact',
    email TEXT NOT NULL,
    org TEXT NOT NULL DEFAULT 'Confidential',
    type TEXT NOT NULL DEFAULT 'Executive Search Request',
    message TEXT DEFAULT '',
    status TEXT NOT NULL DEFAULT 'NEW',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_applications_email ON public.applications(email);
CREATE INDEX IF NOT EXISTS idx_applications_created_at ON public.applications(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_applications_status ON public.applications(status);
