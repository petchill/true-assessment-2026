CREATE TABLE IF NOT EXISTS user_watch_history (
id BIGSERIAL PRIMARY KEY,
user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
content_id BIGINT NOT NULL REFERENCES content(id) ON DELETE CASCADE,
watched_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_watch_history_user ON user_watch_history(user_id);
CREATE INDEX IF NOT EXISTS idx_watch_history_content ON user_watch_history(content_id);
CREATE INDEX IF NOT EXISTS idx_watch_history_composite ON user_watch_history(user_id, watched_at DESC);