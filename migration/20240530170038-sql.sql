
-- +migrate Up
CREATE TABLE users(
    id INTEGER PRIMARY KEY,
    username VARCHAR(20) unique NOT NULL,
    password_hash VARCHAR(255) NOT NULL
);

-- +migrate Down
DROP TABLE users;