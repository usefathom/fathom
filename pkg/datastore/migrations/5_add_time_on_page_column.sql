-- +migrate Up
ALTER TABLE pageviews DROP COLUMN time_on_page;
ALTER TABLE pageviews ADD COLUMN time_on_page INT(4) DEFAULT 0;


-- +migrate Down
ALTER TABLE pageviews DROP COLUMN time_on_page;
