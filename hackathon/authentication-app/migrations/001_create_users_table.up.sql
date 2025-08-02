-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create schema if not exists
CREATE SCHEMA IF NOT EXISTS "authentication-app";

-- Create users table
CREATE TABLE IF NOT EXISTS "authentication-app"."users" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NULL
);

-- Create index on username
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON "authentication-app"."users" (username);
