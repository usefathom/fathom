-- +migrate Up
CREATE TABLE sites (
    id SERIAL PRIMARY KEY NOT NULL,
    tracking_id VARCHAR(8) UNIQUE,
    name VARCHAR(100) NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS sites;