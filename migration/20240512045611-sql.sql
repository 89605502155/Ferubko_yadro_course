
-- +migrate Up
CREATE TABLE comics (
    comics_id VARCHAR(5) PRIMARY KEY,
    url VARCHAR(2048) NOT NULL,
    keywords VARCHAR(64) NOT NULL
);
CREATE TABLE indexes (
    word VARCHAR(64) PRIMARY KEY,
    comics_index INT NOT NULL,
    number_comics_of_index INT NOT NULL
);
-- +migrate Down
DROP TABLE indexes;
DROP TABLE comics;