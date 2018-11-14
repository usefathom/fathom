-- +migrate Up 
INSERT INTO hostnames(name) SELECT hostname FROM daily_page_stats UNION SELECT hostname FROM daily_referrer_stats;

-- +migrate Down

