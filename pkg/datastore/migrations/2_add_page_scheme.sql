-- +migrate Up
ALTER TABLE pages ADD COLUMN scheme ENUM("http", "https") DEFAULT "http";
ALTER TABLE pages DROP COLUMN title;

-- +migrate Down
ALTER TABLE pages DROP COLUMN scheme;
ALTER TABLE pages ADD COLUMN title VARCHAR(255) NULL;
