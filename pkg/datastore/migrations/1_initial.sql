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
   pageviews INT NOT NULL,
   visitors INT NOT NULL,
   entries INT NOT NULL,
   bounce_rate FLOAT NOT NULL,
   avg_duration INT NOT NULL,
   date DATE NOT NULL
);

CREATE TABLE daily_site_stats(
   pageviews INT NOT NULL,
   visitors INT NOT NULL,
   sessions INT NOT NULL,
   bounce_rate FLOAT NOT NULL,
   avg_duration INT NOT NULL,
   date DATE NOT NULL
);

CREATE TABLE daily_referrer_stats(
   url VARCHAR(255) NOT NULL,
   pageviews INT NOT NULL,
   visitors INT NOT NULL,
   bounce_rate FLOAT NOT NULL,
   avg_duration INT NOT NULL,
   date DATE NOT NULL
);

CREATE UNIQUE INDEX unique_user_email ON users(email);
CREATE UNIQUE INDEX unique_daily_site_stats ON daily_site_stats(date);
CREATE UNIQUE INDEX unique_daily_page_stats ON daily_page_stats(pathname, date);
CREATE UNIQUE INDEX unique_daily_referrer_stats ON daily_referrer_stats(url, date);

-- +migrate Down

DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS pageviews;
DROP TABLE IF EXISTS daily_page_stats;
DROP TABLE IF EXISTS daily_site_stats;
DROP TABLE IF EXISTS daily_referrer_stats;

