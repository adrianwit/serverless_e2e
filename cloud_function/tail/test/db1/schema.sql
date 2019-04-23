DROP TABLE IF EXISTS events;
CREATE TABLE events (
  id         INTEGER NOT NULL PRIMARY KEY,
  timestamp  DATETIME,
  event_type INTEGER,
  ua         VARCHAR(255) DEFAULT NULL,
  quantity   DECIMAL(7,2),
  modified   TIMESTAMP    DEFAULT CURRENT_TIMESTAMP
);