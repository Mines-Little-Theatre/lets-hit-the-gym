BEGIN;

PRAGMA user_version = 2;

CREATE TABLE workouts (
  name TEXT PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT NOT NULL,
  color INTEGER NOT NULL
) WITHOUT ROWID;

CREATE TABLE routines (
  id INTEGER PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT NOT NULL
);

CREATE TABLE workout_routines (
  workout_name INTEGER NOT NULL REFERENCES workouts(name),
  routine_id INTEGER NOT NULL REFERENCES routines(id),
  ordinal INTEGER NOT NULL
);
CREATE INDEX workout_routine_idx ON workout_routines(workout_name, ordinal, routine_id);

COMMIT;
