-- Add OAuth support to users table
ALTER TABLE users
    ALTER COLUMN password_hash DROP NOT NULL;

ALTER TABLE users
    ADD COLUMN oauth_provider VARCHAR(50),
    ADD COLUMN oauth_id VARCHAR(255);

-- Create unique index for OAuth provider and ID combination
CREATE UNIQUE INDEX idx_users_oauth_provider_id 
    ON users(oauth_provider, oauth_id) 
    WHERE oauth_provider IS NOT NULL AND oauth_id IS NOT NULL;