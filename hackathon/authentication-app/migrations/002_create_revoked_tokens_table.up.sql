-- Create revoked_tokens table
CREATE TABLE IF NOT EXISTS "authentication-app"."revoked_tokens" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    token_id VARCHAR(255) NOT NULL UNIQUE,
    user_id UUID NULL,
    revoked_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES "authentication-app"."users" (id) ON DELETE CASCADE
);

-- Create indexes for revoked_tokens
CREATE UNIQUE INDEX IF NOT EXISTS idx_revoked_tokens_token_id ON "authentication-app"."revoked_tokens" (token_id);
CREATE INDEX IF NOT EXISTS idx_revoked_tokens_user_id ON "authentication-app"."revoked_tokens" (user_id);
