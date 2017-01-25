-- +migrate Up

CREATE TABLE visitors(
  `id` INTEGER UNSIGNED AUTO_INCREMENT PRIMARY KEY NOT NULL,
  `visitor_key` VARCHAR(255) NOT NULL,
  `ip_address` VARCHAR(100) NOT NULL,
  `device_os` VARCHAR(31) NULL,
  `browser_name` VARCHAR(31) NULL,
  `browser_version` VARCHAR(31) NULL,
  `browser_language` VARCHAR(31) NULL,
  `screen_resolution` VARCHAR(9) NULL,
  `country` CHAR(3) NULL
);
ALTER TABLE visitors ADD UNIQUE(`visitor_key`);

CREATE TABLE pages(
  `id` INTEGER UNSIGNED AUTO_INCREMENT PRIMARY KEY NOT NULL,
  `hostname` VARCHAR(63) NOT NULL,
  `path` VARCHAR(255) NOT NULL,
  `title` VARCHAR(255) NULL
);

CREATE TABLE pageviews(
  `id` INTEGER UNSIGNED AUTO_INCREMENT PRIMARY KEY NOT NULL,
  `page_id` INTEGER UNSIGNED NOT NULL,
  `visitor_id` INTEGER UNSIGNED NOT NULL,
  `referrer_keyword` TEXT NULL,
  `referrer_url` TEXT NULL,
  `timestamp` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE pageviews ADD FOREIGN KEY(`visitor_id`) REFERENCES visitors(`id`);
ALTER TABLE pageviews ADD FOREIGN KEY(`page_id`) REFERENCES pages(`id`);

CREATE TABLE users (
  `id` INTEGER UNSIGNED AUTO_INCREMENT PRIMARY KEY NOT NULL,
  `email` VARCHAR(255) NOT NULL,
  `password` VARCHAR(255) NOT NULL
);

CREATE TABLE options (
  `id` INTEGER UNSIGNED AUTO_INCREMENT PRIMARY KEY NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `value` VARCHAR(255) DEFAULT ''
);
ALTER TABLE options ADD UNIQUE(`name`);

CREATE TABLE `total_pageviews` (
  `id` INTEGER UNSIGNED AUTO_INCREMENT PRIMARY KEY NOT NULL,
  `page_id` INTEGER UNSIGNED NOT NULL,
  `count` INTEGER UNSIGNED NOT NULL DEFAULT 0,
  `count_unique` INTEGER UNSIGNED NOT NULL DEFAULT 0,
  `date` DATE NOT NULL
);
CREATE INDEX total_pageviews_date ON total_pageviews(`date`);
ALTER TABLE total_pageviews ADD UNIQUE(`page_id`, `date`);

CREATE TABLE `total_visitors` (
  `id` INTEGER UNSIGNED AUTO_INCREMENT PRIMARY KEY NOT NULL,
  `count` INTEGER UNSIGNED NOT NULL DEFAULT 0,
  `date` DATE NOT NULL
);
CREATE INDEX total_visitors_date ON total_visitors(`date`);
ALTER TABLE total_visitors ADD UNIQUE(`date`);

CREATE TABLE `total_screens` (
  `id` INTEGER UNSIGNED AUTO_INCREMENT PRIMARY KEY NOT NULL,
  `value` VARCHAR(12) NOT NULL,
  `count` INTEGER UNSIGNED NOT NULL DEFAULT 0,
  `count_unique` INTEGER UNSIGNED NOT NULL DEFAULT 0,
  `date` DATE NOT NULL
);
CREATE INDEX total_screens_date ON total_screens(`date`);
ALTER TABLE total_screens ADD UNIQUE(`value`, `date`);

CREATE TABLE `total_browser_names` (
  `id` INTEGER UNSIGNED AUTO_INCREMENT PRIMARY KEY NOT NULL,
  `value` VARCHAR(50) NOT NULL,
  `count` INTEGER UNSIGNED NOT NULL DEFAULT 0,
  `count_unique` INTEGER UNSIGNED NOT NULL DEFAULT 0,
  `date` DATE NOT NULL
);
CREATE INDEX total_browser_names_date ON total_browser_names(`date`);
ALTER TABLE total_browser_names ADD UNIQUE(`value`, `date`);

CREATE TABLE `total_browser_languages` (
  `id` INTEGER UNSIGNED AUTO_INCREMENT PRIMARY KEY NOT NULL,
  `value` VARCHAR(12) NOT NULL,
  `count` INTEGER UNSIGNED NOT NULL DEFAULT 0,
  `count_unique` INTEGER UNSIGNED NOT NULL DEFAULT 0,
  `date` DATE NOT NULL
);
CREATE INDEX total_browser_languages_date ON total_browser_languages(`date`);
ALTER TABLE total_browser_languages ADD UNIQUE(`value`, `date`);

CREATE TABLE `total_referrers` (
  `id` INTEGER UNSIGNED AUTO_INCREMENT PRIMARY KEY NOT NULL,
  `value` VARCHAR(510) NOT NULL,
  `count` INTEGER UNSIGNED NOT NULL DEFAULT 0,
  `count_unique` INTEGER UNSIGNED NOT NULL DEFAULT 0,
  `date` DATE NOT NULL
);
CREATE INDEX total_referrers_date ON total_referrers(`date`);
ALTER TABLE total_referrers ADD UNIQUE(`value`, `date`);

-- +migrate Down
DROP TABLE IF EXISTS pageviews;
DROP TABLE if exists visitors;
DROP TABLE IF EXISTS pages;
DROP TABLE IF EXISTS sites;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS total_pageviews;
DROP TABLE IF EXISTS total_visitors;
DROP TABLE IF EXISTS total_browser_languages;
DROP TABLE IF EXISTS total_screens;
DROP TABLE IF EXISTS total_browser_names;
DROP TABLE IF EXISTS total_referrers;
DROP TABLE IF EXISTS options;
