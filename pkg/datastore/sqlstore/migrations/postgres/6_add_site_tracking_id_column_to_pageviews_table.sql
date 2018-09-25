-- +migrate Up

TRUNCATE pageviews; -- postgres will fail because of NULL values otherwise
ALTER TABLE pageviews ADD COLUMN site_tracking_id VARCHAR(8) NOT NULL;

-- +migrate Down

ALTER TABLE pageviews DROP COLUMN site_tracking_id;

