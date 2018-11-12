-- +migrate Up
DROP TABLE IF EXISTS daily_page_stats_old;
ALTER TABLE daily_page_stats RENAME TO daily_page_stats_old;
CREATE TABLE daily_page_stats(
   site_id INTEGER NOT NULL DEFAULT 1,
   hostname_id INTEGER NOT NULL,
   pathname_id INTEGER NOT NULL,
   pageviews INTEGER NOT NULL,
   visitors INTEGER NOT NULL,
   entries INTEGER NOT NULL,
   bounce_rate FLOAT NOT NULL,
   known_durations INTEGER NOT NULL DEFAULT 0,
   avg_duration FLOAT NOT NULL,
   date DATE NOT NULL
);
INSERT INTO daily_page_stats 
    SELECT site_id, h.id, p.id, pageviews, visitors, entries, bounce_rate, known_durations, avg_duration, date 
    FROM daily_page_stats_old s 
    LEFT JOIN hostnames h ON h.name = s.hostname 
    LEFT JOIN pathnames p ON p.name = s.pathname;
DROP TABLE daily_page_stats_old;

-- +migrate Down

