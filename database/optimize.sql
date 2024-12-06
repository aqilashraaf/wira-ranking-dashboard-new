-- Create index for username search
CREATE INDEX IF NOT EXISTS idx_accounts_username ON accounts USING gin (username gin_trgm_ops);

-- Create extension if it doesn't exist
CREATE EXTENSION IF NOT EXISTS pg_trgm;
