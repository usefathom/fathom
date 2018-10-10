-- +migrate Up

DROP INDEX IF EXISTS unique_daily_site_stats;
DROP INDEX IF EXISTS unique_daily_page_stats;
DROP INDEX IF EXISTS unique_daily_referrer_stats;

CREATE UNIQUE INDEX unique_daily_site_stats ON daily_site_stats(site_id, date);
CREATE UNIQUE INDEX unique_daily_page_stats ON daily_page_stats(site_id, hostname, pathname, date);
CREATE UNIQUE INDEX unique_daily_referrer_stats ON daily_referrer_stats(site_id, hostname, pathname, date);

-- +migrate Down

DROP INDEX IF EXISTS unique_daily_site_stats;
DROP INDEX IF EXISTS unique_daily_page_stats;
DROP INDEX IF EXISTS unique_daily_referrer_stats;

CREATE UNIQUE INDEX unique_daily_site_stats ON daily_site_stats(date);
CREATE UNIQUE INDEX unique_daily_page_stats ON daily_page_stats(hostname, pathname, date);
CREATE UNIQUE INDEX unique_daily_referrer_stats ON daily_referrer_stats(hostname, pathname, date);