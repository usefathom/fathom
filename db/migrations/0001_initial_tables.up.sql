
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


CREATE TABLE pages(
  `id` INTEGER UNSIGNED AUTO_INCREMENT PRIMARY KEY NOT NULL,
  `hostname` VARCHAR(63) NOT NULL,
  `path` VARCHAR(255) NOT NULL,
  `title` VARCHAR(255) NULL
);

CREATE TABLE users (
  `id` INTEGER UNSIGNED AUTO_INCREMENT PRIMARY KEY NOT NULL,
  `email` VARCHAR(255) NOT NULL,
  `password` VARCHAR(255) NOT NULL
);

CREATE TABLE `archive` (
  `id` INTEGER UNSIGNED AUTO_INCREMENT PRIMARY KEY NOT NULL,
  `metric` VARCHAR(31) NOT NULL,
  `value` VARCHAR(31) NULL,
  `count` INTEGER UNSIGNED NOT NULL,
  `date` DATE NOT NULL
);

CREATE INDEX archive_metric ON archive(`metric`);
CREATE INDEX archive_metric_count ON archive(`metric`, `count`);
CREATE INDEX archive_metric_date ON archive(`metric`, `date`);
