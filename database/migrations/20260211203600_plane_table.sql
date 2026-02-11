-- +goose Up
SELECT 'up SQL query';
CREATE TABLE planes (
    id SERIAL PRIMARY KEY,
    tail_number VARCHAR(50) UNIQUE NOT NULL,
    model VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
SELECT 'down SQL query';
DROP TABLE IF EXISTS planes;
