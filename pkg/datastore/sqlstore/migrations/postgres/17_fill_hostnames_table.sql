-- +migrate Up
INSERT INTO hostnames(name) SELECT DISTINCT(hostname) FROM daily_page_stats;
INSERT INTO hostnames(name) SELECT DISTINCT(hostname) FROM daily_referrer_stats ON CONFLICT(name) DO NOTHING; 

-- +migrate Down

