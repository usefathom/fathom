-- +migrate Up
CREATE TABLE sites (
    id INTEGER PRIMARY KEY,
    tracking_id VARCHAR(8) UNIQUE,
    name VARCHAR(100) NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS sites;