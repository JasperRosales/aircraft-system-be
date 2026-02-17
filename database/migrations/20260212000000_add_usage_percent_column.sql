-- +goose Up
SELECT 'up SQL query';
ALTER TABLE plane_parts
ADD COLUMN IF NOT EXISTS usage_percent NUMERIC GENERATED ALWAYS AS
((usage_hours / NULLIF(usage_limit_hours, 0)) * 100) STORED;

CREATE INDEX IF NOT EXISTS idx_plane_parts_usage_percent
ON plane_parts(usage_percent);


-- +goose Down
SELECT 'down SQL query';
DROP INDEX IF EXISTS idx_plane_parts_usage_percent;
ALTER TABLE plane_parts
DROP COLUMN IF EXISTS usage_percent;

