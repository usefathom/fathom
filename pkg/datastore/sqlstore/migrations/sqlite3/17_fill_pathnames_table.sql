-- +migrate Up
INSERT INTO pathnames(name) SELECT pathname FROM daily_page_stats UNION SELECT pathname FROM daily_referrer_stats;

-- +migrate Down

