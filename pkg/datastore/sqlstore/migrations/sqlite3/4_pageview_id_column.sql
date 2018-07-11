-- +migrate Up

DROP TABLE pageviews;
CREATE TABLE pageviews(
   id VARCHAR(31) NOT NULL,
   hostname VARCHAR(255) NOT NULL,
   pathname VARCHAR(255) NOT NULL,
   is_new_visitor TINYINT(1) NOT NULL,
   is_new_session TINYINT(1) NOT NULL,
   is_unique TINYINT(1) NOT NULL,
   is_bounce TINYINT(1) NULL,
   referrer VARCHAR(255) NULL,
   duration INTEGER(4) NULL,
   timestamp DATETIME NOT NULL
);

-- +migrate Down

DROP TABLE pageviews;
CREATE TABLE pageviews(
   id INTEGER PRIMARY KEY,
   hostname VARCHAR(255) NOT NULL,
   pathname VARCHAR(255) NOT NULL,
   session_id VARCHAR(16) NOT NULL,
   is_new_visitor TINYINT(1) NOT NULL,
   is_new_session TINYINT(1) NOT NULL,
   is_unique TINYINT(1) NOT NULL,
   is_bounce TINYINT(1) NULL,
   referrer VARCHAR(255) NULL,
   duration INTEGER(4) NULL,
   timestamp DATETIME NOT NULL
);
