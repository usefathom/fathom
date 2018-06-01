-- +migrate Up

ALTER TABLE daily_referrer_stats ADD COLUMN groupname VARCHAR(255);
ALTER TABLE daily_referrer_stats ADD COLUMN hostname VARCHAR(255);
ALTER TABLE daily_referrer_stats ADD COLUMN pathname VARCHAR(255);

UPDATE daily_referrer_stats SET hostname = SUBSTR(url, 0, (INSTR(url, '://')+3+INSTR(SUBSTR(url, INSTR(url, '://')+3), '/')-1)) WHERE url != '' AND (hostname = '' OR hostname IS NULL);
UPDATE daily_referrer_stats SET pathname = SUBSTR(url, LENGTH(hostname)+1) WHERE url != '' AND (pathname = '' OR pathname IS NULL);

-- drop `url` column... oh sqlite
ALTER TABLE daily_referrer_stats RENAME TO daily_referrer_stats_old;
CREATE TABLE daily_referrer_stats(
   hostname VARCHAR(255) NOT NULL,
   pathname VARCHAR(255) NOT NULL,
   groupname VARCHAR(255) NULL,
   pageviews INTEGER NOT NULL,
   visitors INTEGER NOT NULL,
   bounce_rate FLOAT NOT NULL,
   avg_duration FLOAT NOT NULL,
   known_durations INTEGER NOT NULL DEFAULT 0,
   date DATE NOT NULL
);
INSERT INTO daily_referrer_stats SELECT hostname, pathname, groupname, pageviews, visitors, bounce_rate, avg_duration, known_durations, date FROM daily_referrer_stats_old;

-- +migrate Down

-- TODO....
