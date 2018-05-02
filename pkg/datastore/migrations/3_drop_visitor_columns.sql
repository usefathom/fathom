-- +migrate Up
ALTER TABLE visitors DROP COLUMN ip_address;

-- +migrate Down
ALTER TABLE visitors ADD COLUMN ip_address VARCHAR(100) NULL;
