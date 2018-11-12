-- +migrate Up
CREATE TABLE hostnames(
   id INTEGER PRIMARY KEY,
   name VARCHAR(255) NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS hostnames;
