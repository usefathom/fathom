-- +migrate Up
CREATE UNIQUE INDEX unique_hostnames_name ON hostnames(name);

-- +migrate Down
DROP INDEX IF EXISTS unique_hostnames_name;
