-- +migrate Up

ALTER TABLE daily_site_stats ALTER COLUMN bounce_rate TYPE NUMERIC;
ALTER TABLE daily_page_stats ALTER COLUMN bounce_rate TYPE NUMERIC;
ALTER TABLE daily_referrer_stats ALTER COLUMN bounce_rate TYPE NUMERIC;
ALTER TABLE daily_site_stats ALTER COLUMN avg_duration TYPE NUMERIC;
ALTER TABLE daily_page_stats ALTER COLUMN avg_duration TYPE NUMERIC;
ALTER TABLE daily_referrer_stats ALTER COLUMN avg_duration TYPE NUMERIC;

-- +migrate Down


