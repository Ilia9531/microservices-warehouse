-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    login  VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    email         VARCHAR(255) UNIQUE NOT NULL,
    notification_methods  JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- Индекс для быстрого поиска по uuid
CREATE INDEX idx_users_uuid ON users(uuid);

-- Индекс для быстрого поиска по логину
CREATE INDEX idx_users_login ON users(login);

-- +goose Down
DROP INDEX IF EXISTS idx_users_uuid;
DROP INDEX IF EXISTS idx_users_login;
DROP TABLE IF EXISTS users;