-- +migrate Up

CREATE TABLE daily_referrer_stats(
   url VARCHAR(255) NOT NULL,
   visitors INT NOT NULL,
   pageviews INT NOT NULL,
   bounced INT NOT NULL,
   bounced_n INT NOT NULL,
   avg_duration INT NOT NULL,
   avg_duration_n INT NOT NULL,
   date DATE NOT NULL
);

-- +migrate Down

DROP TABLE IF EXISTS daily_referrer_stats;
