-- +migrate Up

DROP INDEX unique_daily_referrer_stats ON daily_referrer_stats;

ALTER TABLE daily_referrer_stats ADD COLUMN groupname VARCHAR(255);
ALTER TABLE daily_referrer_stats ADD COLUMN hostname VARCHAR(255);
ALTER TABLE daily_referrer_stats ADD COLUMN pathname VARCHAR(255);


UPDATE daily_referrer_stats SET hostname = SUBSTRING_INDEX( url, "/", 3) WHERE url != "" AND ( hostname = "" OR hostname IS NULL);
UPDATE daily_referrer_stats SET pathname = REPLACE(url, hostname, "") WHERE url != "" AND (pathname = '' OR pathname IS NULL);

CREATE UNIQUE INDEX unique_daily_referrer_stats ON daily_referrer_stats(hostname(100), pathname(100), date);
ALTER TABLE daily_referrer_stats DROP COLUMN url;

-- +migrate Down

ALTER TABLE daily_referrer_stats DROP COLUMN groupname;
ALTER TABLE daily_referrer_stats DROP COLUMN hostname;
ALTER TABLE daily_referrer_stats DROP COLUMN pathname;

ALTER TABLE daily_referrer_stats ADD COLUMN url VARCHAR(255) NOT NULL;

