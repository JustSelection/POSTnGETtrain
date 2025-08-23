ALTER TABLE tasks ADD COLUMN user_id VARCHAR(50) REFERENCES users(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_tasks_user_id ON tasks (user_id) WHERE deleted_at IS NULL;