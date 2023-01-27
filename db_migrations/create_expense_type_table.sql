CREATE TABLE IF NOT EXISTS expense_type
(
    id   uuid PRIMARY KEY,
    name VARCHAR(16),
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP
);