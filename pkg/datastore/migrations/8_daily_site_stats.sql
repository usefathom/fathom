-- +migrate Up

CREATE TABLE daily_site_stats(
   visitors INT NOT NULL,
   pageviews INT NOT NULL,
   bounced INT NOT NULL,
   bounced_n INT NOT NULL,
   avg_duration INT NOT NULL,
   avg_duration_n INT NOT NULL,
   date DATE NOT NULL
);

-- +migrate Down

DROP TABLE IF EXISTS daily_site_stats;
