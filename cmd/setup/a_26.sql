ALTER TABLE IF EXISTS auth.user_sessions ADD COLUMN IF NOT EXISTS external_user_id TEXT NOT NULL DEFAULT '';