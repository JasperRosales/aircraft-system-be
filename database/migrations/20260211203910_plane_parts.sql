-- +goose Up
SELECT 'up SQL query';
CREATE TABLE plane_parts (
    id SERIAL PRIMARY KEY,
    plane_id INTEGER NOT NULL REFERENCES planes(id) ON DELETE CASCADE,

    part_name VARCHAR(255) NOT NULL,
    serial_number VARCHAR(100) UNIQUE NOT NULL,
    category VARCHAR(150) NOT NULL, 

    usage_hours NUMERIC(10,2) DEFAULT 0,
    usage_limit_hours NUMERIC(10,2) NOT NULL,

    installed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
SELECT 'down SQL query';
DROP TABLE IF EXISTS plane_parts;
