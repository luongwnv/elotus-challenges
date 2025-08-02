-- Create file_uploads table
CREATE TABLE IF NOT EXISTS "authentication-app"."file_uploads" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    filename VARCHAR(255) NOT NULL,
    original_name VARCHAR(255) NOT NULL,
    content_type VARCHAR(255) NOT NULL,
    size BIGINT NOT NULL CHECK (size > 0),
    file_path VARCHAR(500) NOT NULL UNIQUE,
    user_agent TEXT,
    ip_address VARCHAR(45),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NULL,
    FOREIGN KEY (user_id) REFERENCES "authentication-app"."users" (id) ON DELETE CASCADE
);

-- Create indexes for file_uploads
CREATE INDEX IF NOT EXISTS idx_file_uploads_user_id ON "authentication-app"."file_uploads" (user_id);
CREATE INDEX IF NOT EXISTS idx_file_uploads_created_at ON "authentication-app"."file_uploads" (created_at);
