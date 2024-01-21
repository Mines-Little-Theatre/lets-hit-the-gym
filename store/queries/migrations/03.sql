PRAGMA user_version = 3;

CREATE TEMPORARY TABLE workout_routines_old AS SELECT * FROM workout_routines;
DROP TABLE workout_routines;
CREATE TEMPORARY TABLE workouts_old AS SELECT * FROM workouts;
DROP TABLE workouts;

CREATE TEMPORARY TABLE workout_id_name_map (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL
);

INSERT INTO workout_id_name_map (name)
  SELECT name FROM workouts_old;

CREATE TABLE workouts (
  id INTEGER PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT NOT NULL,
  color INTEGER NOT NULL
);

INSERT INTO workouts (id, title, description, color)
  SELECT m.id, w.title, w.description, w.color
  FROM workout_id_name_map AS m
  JOIN workouts_old AS w USING (name);

CREATE TABLE workout_routines (
  workout_id INTEGER NOT NULL REFERENCES workouts(id),
  routine_id INTEGER NOT NULL REFERENCES routines(id),
  ordinal INTEGER NOT NULL
);

INSERT INTO workout_routines (workout_id, routine_id, ordinal)
  SELECT m.id, x.routine_id, x.ordinal
  FROM workout_id_name_map AS m
  JOIN workout_routines_old AS x ON m.name = x.workout_name;

CREATE INDEX workout_routine_idx ON workout_routines(workout_id, ordinal, routine_id);

DROP TABLE workout_id_name_map;
DROP TABLE workout_routines_old;
DROP TABLE workouts_old;

CREATE TABLE days (
  name TEXT PRIMARY KEY,
  open_hour INTEGER NOT NULL,
  close_hour INTEGER NOT NULL,
  workout_id INTEGER REFERENCES workouts(id)
) WITHOUT ROWID;
