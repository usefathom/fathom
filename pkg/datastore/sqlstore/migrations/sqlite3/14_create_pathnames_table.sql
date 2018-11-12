-- +migrate Up
CREATE TABLE pathnames(
   id INTEGER PRIMARY KEY,
   name VARCHAR(255) NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS pathnames;
