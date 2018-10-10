-- +migrate Up

ALTER TABLE pageviews ADD COLUMN is_finished TINYINT(1) NOT NULL DEFAULT 0;

-- +migrate Down

ALTER TABLE pageviews DROP COLUMN is_finished;
