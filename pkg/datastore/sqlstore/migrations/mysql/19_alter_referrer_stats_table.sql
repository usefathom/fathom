-- +migrate Up
DROP TABLE IF EXISTS daily_referrer_stats_old;
RENAME TABLE daily_referrer_stats TO daily_referrer_stats_old;
CREATE TABLE daily_referrer_stats(
    site_id INTEGER NOT NULL DEFAULT 1, 
    hostname_id INTEGER NOT NULL,
    pathname_id INTEGER NOT NULL,
    groupname VARCHAR(255) NULL, 
    pageviews INTEGER NOT NULL, 
    visitors INTEGER NOT NULL, 
    bounce_rate FLOAT NOT NULL, 
    known_durations INTEGER NOT NULL DEFAULT 0, 
    avg_duration FLOAT NOT NULL, 
    date DATE NOT NULL 
) CHARACTER SET=utf8;
INSERT INTO daily_referrer_stats 
    SELECT site_id, h.id, p.id, groupname, pageviews, visitors, bounce_rate, known_durations, avg_duration, date 
    FROM daily_referrer_stats_old s 
    LEFT JOIN hostnames h ON h.name = s.hostname 
    LEFT JOIN pathnames p ON p.name = s.pathname;
DROP TABLE daily_referrer_stats_old;

-- +migrate Down

