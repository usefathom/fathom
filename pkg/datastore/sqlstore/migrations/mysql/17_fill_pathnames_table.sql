-- +migrate Up
INSERT IGNORE INTO pathnames(name) SELECT DISTINCT(pathname) FROM daily_page_stats;
INSERT IGNORE INTO pathnames(name) SELECT DISTINCT(pathname) FROM daily_referrer_stats;

-- +migrate Down

