-- +migrate Up
CREATE TABLE site_stats(
   site_id INTEGER NOT NULL DEFAULT 1,
   pageviews INTEGER NOT NULL,
   visitors INTEGER NOT NULL,
   sessions INTEGER NOT NULL,
   bounce_rate FLOAT NOT NULL,
   known_durations INTEGER NOT NULL DEFAULT 0,
   avg_duration FLOAT NOT NULL,
   ts DATETIME NOT NULL
) CHARACTER SET=utf8;
INSERT INTO site_stats 
    SELECT site_id, pageviews, visitors, sessions, bounce_rate, known_durations, avg_duration,  CONCAT(date, ' 00:00:00')
    FROM daily_site_stats s ;
DROP TABLE daily_site_stats;

-- +migrate Down