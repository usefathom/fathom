-- +migrate Up

CREATE TABLE users (
  id INT AUTO_INCREMENT PRIMARY KEY NOT NULL,
  email VARCHAR(255) NOT NULL,
  password VARCHAR(255) NOT NULL
);

CREATE TABLE pageviews(
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

CREATE TABLE daily_page_stats(
   pathname VARCHAR(255) NOT NULL,
   views INT NOT NULL,
   unique_views INT NOT NULL,
   entries INT NOT NULL,
   bounces INT NOT NULL,
   avg_duration INT NOT NULL,
   date DATE NOT NULL
);

CREATE TABLE daily_site_stats(
   visitors INT NOT NULL,
   pageviews INT NOT NULL,
   sessions INT NOT NULL,
   bounces INT NOT NULL,
   avg_duration INT NOT NULL,
   date DATE NOT NULL
);

CREATE TABLE daily_referrer_stats(
   url VARCHAR(255) NOT NULL,
   visitors INT NOT NULL,
   pageviews INT NOT NULL,
   bounces INT NOT NULL,
   avg_duration INT NOT NULL,
   date DATE NOT NULL
);

-- +migrate Down

DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS pageviews;
DROP TABLE IF EXISTS daily_page_stats;
DROP TABLE IF EXISTS daily_site_stats;
DROP TABLE IF EXISTS daily_referrer_stats;

