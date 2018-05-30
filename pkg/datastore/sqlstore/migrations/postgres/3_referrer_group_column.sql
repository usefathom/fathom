-- +migrate Up

ALTER TABLE daily_referrer_stats ADD COLUMN groupname VARCHAR(255) NULL;
ALTER TABLE daily_referrer_stats ADD COLUMN hostname VARCHAR(255);
ALTER TABLE daily_referrer_stats ADD COLUMN pathname VARCHAR(255);

UPDATE daily_referrer_stats SET hostname = CONCAT( SPLIT_PART(url, '://', 1), '://', SPLIT_PART(SPLIT_PART(url, '://', 2), '/', 1) ) WHERE url != '' AND hostname = '';
UPDATE daily_referrer_stats SET pathname = SPLIT_PART( url, hostname, 2 ) WHERE url != '' AND pathname = '';

ALTER TABLE daily_referrer_stats DROP COLUMN url;

-- +migrate Down

ALTER TABLE daily_referrer_stats DROP COLUMN groupname;
ALTER TABLE daily_referrer_stats DROP COLUMN hostname;
ALTER TABLE daily_referrer_stats DROP COLUMN pathname;

ALTER TABLE daily_referrer_stats ADD COLUMN url VARCHAR(255) NOT NULL;
