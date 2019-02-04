-- +migrate Up
DROP TABLE IF EXISTS sites_old;
ALTER TABLE sites RENAME TO sites_old;
CREATE TABLE sites (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `tracking_id` VARCHAR(8) UNIQUE,
    `name` VARCHAR(100) NOT NULL
);
INSERT INTO sites SELECT `id`, `tracking_id`, `name` FROM sites_old;

-- +migrate Down
DROP TABLE IF EXISTS sites_old;
ALTER TABLE sites RENAME TO sites_old;
CREATE TABLE sites (
    `id` INTEGER PRIMARY KEY,
    `tracking_id` VARCHAR(8) UNIQUE,
    `name` VARCHAR(100) NOT NULL
);
INSERT INTO sites SELECT `id`, `tracking_id`, `name` FROM sites_old;
