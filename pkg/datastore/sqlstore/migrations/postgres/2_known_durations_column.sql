-- +migrate Up

ALTER TABLE daily_site_stats ADD COLUMN known_durations INTEGER NOT NULL DEFAULT 0;
ALTER TABLE daily_page_stats ADD COLUMN known_durations INTEGER NOT NULL DEFAULT 0;
ALTER TABLE daily_referrer_stats ADD COLUMN known_durations INTEGER NOT NULL DEFAULT 0;

-- +migrate Down

ALTER TABLE daily_site_stats DROP COLUMN known_durations;
ALTER TABLE daily_page_stats DROP COLUMN known_durations;
ALTER TABLE daily_referrer_stats DROP COLUMN known_durations;


