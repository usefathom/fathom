-- +migrate Up

TRUNCATE pageviews; -- postgres will fail because of NULL values otherwise
ALTER TABLE pageviews DROP COLUMN session_id;
ALTER TABLE pageviews DROP COLUMN id;
ALTER TABLE pageviews ADD COLUMN id VARCHAR(31) NOT NULL;

-- +migrate Down

ALTER TABLE pageviews DROP COLUMN id;
ALTER TABLE pageviews ADD COLUMN id INT AUTO_INCREMENT PRIMARY KEY NOT NULL;
ALTER TABLE pageviews ADD COLUMN session_id VARCHAR(16) NOT NULL;



