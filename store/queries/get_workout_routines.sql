-- ?1: workout name
-- returns any number of rows (text, text): title, description
SELECT r.title, r.description
FROM workout_routines AS x, routines AS r
WHERE x.workout_name = ?1 AND x.routine_id = r.id
ORDER BY x.ordinal;
