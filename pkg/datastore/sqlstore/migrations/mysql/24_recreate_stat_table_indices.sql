-- +migrate Up
CREATE UNIQUE INDEX unique_page_stats ON page_stats(site_id, hostname_id, pathname_id, ts);
CREATE UNIQUE INDEX unique_referrer_stats ON referrer_stats(site_id, hostname_id, pathname_id, ts);
CREATE UNIQUE INDEX unique_site_stats ON site_stats(site_id, ts);

-- +migrate Down
DROP INDEX unique_page_stats ON page_stats;
DROP INDEX unique_referrer_stats ON referrer_stats;
DROP INDEX unique_site_stats ON site_stats;
