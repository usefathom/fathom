-- +migrate Up

ALTER TABLE pageviews ADD COLUMN is_finished BOOLEAN NOT NULL DEFAULT FALSE;

-- +migrate Down

ALTER TABLE pageviews DROP COLUMN is_finished;
