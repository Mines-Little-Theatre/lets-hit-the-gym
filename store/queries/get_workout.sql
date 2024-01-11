-- ?1: workout name
-- returns zero or one row (text, text, integer): title, description, color
SELECT title, description, color FROM workouts WHERE name = ?1;
