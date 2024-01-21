PRAGMA user_version = 4;

CREATE TABLE arrivals (
  user_id TEXT NOT NULL,
  hour INTEGER NOT NULL
);
CREATE UNIQUE INDEX arrivals_idx ON arrivals(hour, user_id);
CREATE INDEX arrivals_user_id_idx ON arrivals(user_id);
