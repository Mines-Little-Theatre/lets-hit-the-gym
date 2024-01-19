-- ?1: day name
-- returns any number of rows (
--   integer: open hour, integer: close hour,
--   text?: workout title, text?: workout description, integer?: workout color,
--   text?: routine title, text?: routine description
-- )
SELECT d.open_hour, d.close_hour,
       w.title, w.description, w.color,
       r.title, r.description
FROM days AS d
  LEFT JOIN workouts AS w ON d.workout_id = w.id
  LEFT JOIN workout_routines AS x ON w.id = x.workout_id
  LEFT JOIN routines AS r ON x.routine_id = r.id
WHERE d.name = ?1
ORDER BY x.ordinal;
