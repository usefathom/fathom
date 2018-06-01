-- +migrate Up

ALTER TABLE daily_referrer_stats ADD COLUMN groupname VARCHAR(255);
ALTER TABLE daily_referrer_stats ADD COLUMN hostname VARCHAR(255);
ALTER TABLE daily_referrer_stats ADD COLUMN pathname VARCHAR(255);

UPDATE daily_referrer_stats SET hostname = CONCAT( SPLIT_PART(url, '://', 1), '://', SPLIT_PART(SPLIT_PART(url, '://', 2), '/', 1) ) WHERE url != '' AND ( hostname = "" OR hostname IS NULL);
UPDATE daily_referrer_stats SET pathname = SPLIT_PART( url, hostname, 2 ) WHERE url != '' AND (pathname = '' OR pathname IS NULL);

DROP INDEX IF EXISTS unique_daily_referrer_stats;
ALTER TABLE daily_referrer_stats DROP COLUMN url;
CREATE UNIQUE INDEX unique_daily_referrer_stats ON daily_referrer_stats(hostname, pathname, date);

-- +migrate Down

ALTER TABLE daily_referrer_stats DROP COLUMN groupname;
ALTER TABLE daily_referrer_stats DROP COLUMN hostname;
ALTER TABLE daily_referrer_stats DROP COLUMN pathname;

ALTER TABLE daily_referrer_stats ADD COLUMN url VARCHAR(255) NOT NULL;
