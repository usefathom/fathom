-- +migrate Up

TRUNCATE pageviews;
ALTER TABLE pageviews ADD COLUMN site_tracking_id VARCHAR(8) NOT NULL;

-- +migrate Down

ALTER TABLE pageviews DROP COLUMN site_tracking_id;

