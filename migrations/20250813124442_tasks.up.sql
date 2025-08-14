CREATE TABLE IF NOT EXISTS tasks (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255) NOT NULL CHECK (name <> ''),
    is_done BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP DEFAULT NULL,
    CONSTRAINT valid_dates CHECK (created_at <= updated_at)
);

CREATE INDEX IF NOT EXISTS idx_tasks_is_done ON tasks(is_done) 
WHERE deleted_at IS NULL;