-- +migrate Up
INSERT INTO pathnames(name) SELECT DISTINCT(pathname) FROM daily_page_stats ON CONFLICT(name) DO NOTHING; ;
INSERT INTO pathnames(name) SELECT DISTINCT(pathname) FROM daily_referrer_stats ON CONFLICT(name) DO NOTHING; ;

-- +migrate Down

