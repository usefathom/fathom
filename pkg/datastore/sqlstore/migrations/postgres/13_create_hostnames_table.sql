-- +migrate Up
CREATE TABLE hostnames(
   id SERIAL PRIMARY KEY NOT NULL,
   name VARCHAR(255) NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS hostnames;
