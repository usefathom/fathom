-- +migrate Up
DROP INDEX IF EXISTS unique_daily_page_stats;
DROP INDEX IF EXISTS unique_daily_referrer_stats;
CREATE UNIQUE INDEX unique_page_stats ON page_stats(site_id, hostname_id, pathname_id, ts);
CREATE UNIQUE INDEX unique_referrer_stats ON referrer_stats(site_id, hostname_id, pathname_id, ts);
CREATE UNIQUE INDEX unique_site_stats ON site_stats(site_id, ts);

-- +migrate Down
DROP INDEX IF EXISTS unique_page_stats;
DROP INDEX IF EXISTS unique_referrer_stats;
DROP INDEX IF EXISTS unique_site_stats;