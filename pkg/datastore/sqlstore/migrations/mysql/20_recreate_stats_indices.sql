-- +migrate Up
CREATE UNIQUE INDEX unique_daily_page_stats ON daily_page_stats(site_id, hostname_id, pathname_id, date);
CREATE UNIQUE INDEX unique_daily_referrer_stats ON daily_referrer_stats(site_id, hostname_id, pathname_id, date);

-- +migrate Down
DROP INDEX unique_daily_page_stats ON daily_page_stats;
DROP INDEX unique_daily_referrer_stats ON daily_referrer_stats;