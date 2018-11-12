-- +migrate Up
INSERT OR IGNORE INTO hostnames(name) SELECT DISTINCT(hostname) FROM daily_page_stats;
INSERT OR IGNORE INTO hostnames(name) SELECT DISTINCT(hostname) FROM daily_referrer_stats; 

-- +migrate Down

