-- +goose Up
SELECT 'up SQL query';
CREATE INDEX idx_plane_parts_plane_id
ON plane_parts(plane_id);

CREATE INDEX idx_plane_parts_category
ON plane_parts(category);

CREATE INDEX idx_plane_parts_usage_ratio
ON plane_parts ((usage_hours / usage_limit_hours));

ALTER TABLE plane_parts
ADD COLUMN usage_percent NUMERIC GENERATED ALWAYS AS
((usage_hours / usage_limit_hours) * 100) STORED;

CREATE INDEX idx_plane_parts_usage_percent
ON plane_parts(usage_percent);


-- +goose Down
SELECT 'down SQL query';
DROP INDEX IF EXISTS idx_plane_parts_usage_percent;
ALTER TABLE plane_parts
DROP COLUMN IF EXISTS usage_percent;
