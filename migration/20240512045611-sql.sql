
-- +migrate Up
CREATE TABLE comics (
    id INTEGER PRIMARY KEY,
    comics_id VARCHAR(5) NOT NULL,
    url VARCHAR(2048) NOT NULL,
    keywords VARCHAR(64) NOT NULL
);
CREATE TABLE indexes (
    id INTEGER PRIMARY KEY,
    word VARCHAR(64) NOT NULL,
    comics_index INT NOT NULL,
    number_comics_of_index INT NOT NULL
);
-- +migrate Down
DROP TABLE indexes;
DROP TABLE comics;