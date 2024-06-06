
-- +migrate Up
CREATE TABLE users(
    username VARCHAR(30) PRIMARY KEY,
    password_hash VARCHAR(255) NOT NULL
);

-- +migrate Down
DROP TABLE users;