-- +migrate Up
CREATE UNIQUE INDEX unique_pathnames_name ON pathnames(name(100));

-- +migrate Down
DROP INDEX IF EXISTS unique_pathnames_name;
