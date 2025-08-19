CREATE TABLE IF NOT EXISTS users (
                                     id VARCHAR(50) PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE CHECK (email <> ''),
    password VARCHAR(255) NOT NULL CHECK (password <> ''),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP DEFAULT NULL,
    CONSTRAINT valid_dates CHECK (created_at <= updated_at)
    );

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email) WHERE deleted_at IS NULL;