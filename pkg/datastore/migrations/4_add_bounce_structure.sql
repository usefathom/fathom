-- +migrate Up
ALTER TABLE pageviews ADD COLUMN bounced TINYINT(1) NULL;

CREATE TABLE total_bounced (
  id INT AUTO_INCREMENT PRIMARY KEY NOT NULL,
  page_id INTEGER NOT NULL,
  count INTEGER NOT NULL DEFAULT 0,
  count_unique INTEGER NOT NULL DEFAULT 0,
  date DATE NOT NULL
);

CREATE INDEX total_bounced_date ON total_bounced(date);

ALTER TABLE total_bounced ADD UNIQUE(page_id, date);


-- +migrate Down
DROP TABLE total_bounced;

ALTER TABLE pageviews DROP COLUMN bounced;
