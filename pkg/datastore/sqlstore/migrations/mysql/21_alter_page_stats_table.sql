-- +migrate Up
CREATE TABLE page_stats(
   site_id INTEGER NOT NULL DEFAULT 1,
   hostname_id INTEGER NOT NULL,
   pathname_id INTEGER NOT NULL,
   pageviews INTEGER NOT NULL,
   visitors INTEGER NOT NULL,
   entries INTEGER NOT NULL,
   bounce_rate FLOAT NOT NULL,
   known_durations INTEGER NOT NULL DEFAULT 0,
   avg_duration FLOAT NOT NULL,
   ts DATETIME NOT NULL
) CHARACTER SET=utf8;
INSERT INTO page_stats 
    SELECT site_id, hostname_id, pathname_id, pageviews, visitors, entries, bounce_rate, known_durations, avg_duration, CONCAT(date, ' 00:00:00')
    FROM daily_page_stats s ;
DROP TABLE daily_page_stats;

-- +migrate Down