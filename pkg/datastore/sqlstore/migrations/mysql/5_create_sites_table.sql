-- +migrate Up
CREATE TABLE sites (
    id INTEGER PRIMARY KEY NOT NULL,
    tracking_id VARCHAR(8) UNIQUE,
    name VARCHAR(100) NOT NULL
) CHARACTER SET=utf8;

-- +migrate Down
DROP TABLE IF EXISTS sites;