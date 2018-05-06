-- +migrate Up

CREATE TABLE raw_pageviews(
   id INT AUTO_INCREMENT PRIMARY KEY NOT NULL,
   session_id VARCHAR(16) NOT NULL,
   pathname VARCHAR(255) NOT NULL,
   is_new_visitor TINYINT(1) NOT NULL,
   is_unique TINYINT(1) NOT NULL,
   is_bounce TINYINT(1) NULL,
   referrer VARCHAR(255) NULL,
   duration INT(4) NULL,
   timestamp DATETIME NOT NULL
);

-- +migrate Down

DROP TABLE IF EXISTS raw_pageviews;
