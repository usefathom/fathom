CREATE TABLE visits (
  id INTEGER AUTO_INCREMENT PRIMARY KEY,
  path TEXT NOT NULL,
  ip_address VARCHAR(100) NOT NULL,
  referrer_keyword VARCHAR(255) NULL,
  referrer_type VARCHAR(255) NULL,
  referrer_url TEXT NULL,
  device_brand VARCHAR(100) NULL,
  device_model VARCHAR(100) NULL,
  device_type VARCHAR(100) NULL,
  device_os VARCHAR(100) NULL,
  browser_name VARCHAR(31) NULL,
  browser_version VARCHAR(31) NULL,
  browser_language VARCHAR(31) NULL,
  screen_resolution VARCHAR(9) NULL,
  visitor_returning TINYINT(1) DEFAULT 0,
  country CHAR(3) NULL,
  timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)
