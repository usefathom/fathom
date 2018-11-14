-- +migrate Up
CREATE TABLE referrer_stats(
    site_id INTEGER NOT NULL DEFAULT 1, 
    hostname_id INTEGER NOT NULL,
    pathname_id INTEGER NOT NULL,
    groupname VARCHAR(255) NULL, 
    pageviews INTEGER NOT NULL, 
    visitors INTEGER NOT NULL, 
    bounce_rate FLOAT NOT NULL, 
    known_durations INTEGER NOT NULL DEFAULT 0, 
    avg_duration FLOAT NOT NULL, 
    ts TIMESTAMP WITHOUT TIME ZONE NOT NULL 
);
INSERT INTO referrer_stats 
    SELECT site_id, hostname_id, pathname_id, groupname, pageviews, visitors, bounce_rate, known_durations, avg_duration, (date || ' 00:00:00')::timestamp
    FROM daily_referrer_stats s;
DROP TABLE daily_referrer_stats;

-- +migrate Down
