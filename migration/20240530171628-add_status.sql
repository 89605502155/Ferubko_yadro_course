
-- +migrate Up
ALTER TABLE users ADD COLUMN status TEXT CHECK( status IN ('admin','user','content manager')) 
    NOT NULL DEFAULT 'user';

-- +migrate Down
ALTER TABLE users DROP COLUMN status;