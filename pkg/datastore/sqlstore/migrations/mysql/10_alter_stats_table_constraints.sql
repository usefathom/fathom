-- +migrate Up

DROP INDEX unique_daily_site_stats ON daily_site_stats; 
DROP INDEX unique_daily_page_stats ON daily_page_stats;
DROP INDEX unique_daily_referrer_stats ON daily_referrer_stats;

CREATE UNIQUE INDEX unique_daily_site_stats ON daily_site_stats(site_id, date);
CREATE UNIQUE INDEX unique_daily_page_stats ON daily_page_stats(site_id, hostname(100), pathname(100), date);
CREATE UNIQUE INDEX unique_daily_referrer_stats ON daily_referrer_stats(site_id, hostname(100), pathname(100), date);

-- +migrate Down

DROP INDEX unique_daily_site_stats ON daily_site_stats; 
DROP INDEX unique_daily_page_stats ON daily_page_stats;
DROP INDEX unique_daily_referrer_stats ON daily_referrer_stats;

CREATE UNIQUE INDEX unique_daily_site_stats ON daily_site_stats(date);
CREATE UNIQUE INDEX unique_daily_page_stats ON daily_page_stats(hostname(100), pathname(100), date);
CREATE UNIQUE INDEX unique_daily_referrer_stats ON daily_referrer_stats(hostname(100), pathname(100), date);