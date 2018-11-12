-- +migrate Up
CREATE UNIQUE INDEX unique_pathnames_name ON pathnames(name);

-- +migrate Down
DROP INDEX IF EXISTS unique_pathnames_name;
