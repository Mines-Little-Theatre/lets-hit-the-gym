-- ?1: user id
-- ?2: hour
-- returns no rows, inserts (user id, hour) into arrivals
INSERT INTO arrivals (user_id, hour) VALUES (?1, ?2)
  ON CONFLICT DO NOTHING;
