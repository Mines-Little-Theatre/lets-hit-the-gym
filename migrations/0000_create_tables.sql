CREATE TABLE kv_store (
  key TEXT PRIMARY KEY,
  value ANY
) STRICT, WITHOUT ROWID;

CREATE TABLE workouts (
  id INTEGER PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT NOT NULL,
  color INTEGER NOT NULL
) STRICT;

CREATE TABLE routines (
  id INTEGER PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT NOT NULL
) STRICT;

CREATE TABLE workout_routines (
  workout_id INTEGER NOT NULL REFERENCES workouts(id),
  routine_id INTEGER NOT NULL REFERENCES routines(id),
  ordinal INTEGER NOT NULL
) STRICT;
CREATE INDEX workout_routine_idx ON workout_routines(workout_id, ordinal, routine_id);

CREATE TABLE weekdays (
  id INTEGER PRIMARY KEY,
  post_hour INTEGER NOT NULL,
  open_hour INTEGER NOT NULL,
  close_hour INTEGER NOT NULL,
  workout_id INTEGER REFERENCES workouts(id)
) STRICT;

CREATE TABLE arrivals (
  hour INTEGER NOT NULL,
  user_id TEXT NOT NULL
) STRICT;
CREATE UNIQUE INDEX arrivals_idx ON arrivals(hour, user_id);
CREATE INDEX arrivals_user_id_idx ON arrivals(user_id);
