-- +migrate Up
INSERT IGNORE INTO hostnames(name) SELECT DISTINCT(hostname) FROM daily_page_stats;
INSERT IGNORE INTO hostnames(name) SELECT DISTINCT(hostname) FROM daily_referrer_stats; 

-- +migrate Down

