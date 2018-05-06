-- +migrate Up

CREATE TABLE daily_page_stats(
   pathname VARCHAR(255) NOT NULL,
   views INT NOT NULL,
   unique_views INT NOT NULL,
   bounced INT NOT NULL,
   bounced_n INT NOT NULL,
   avg_duration INT NOT NULL,
   avg_duration_n INT NOT NULL,
   date DATE NOT NULL
);

-- +migrate Down

DROP TABLE IF EXISTS daily_page_stats;
