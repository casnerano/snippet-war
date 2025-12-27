-- Create users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    telegram_user_id BIGINT NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

-- Add table comment
COMMENT ON TABLE users IS 'Stores Telegram users who use the Snippet War app';

-- Add column comments
COMMENT ON COLUMN users.id IS 'Unique identifier for the user';
COMMENT ON COLUMN users.telegram_user_id IS 'Telegram user ID (unique identifier from Telegram)';
COMMENT ON COLUMN users.created_at IS 'Timestamp when the user was first created';

-- Create unique index on telegram_user_id for fast lookups
CREATE UNIQUE INDEX idx_users_telegram_user_id ON users(telegram_user_id);
