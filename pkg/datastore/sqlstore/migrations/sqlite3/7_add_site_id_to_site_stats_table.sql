-- +migrate Up

ALTER TABLE daily_site_stats ADD COLUMN site_id INTEGER NOT NULL DEFAULT 1;

-- +migrate Down


