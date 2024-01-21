-- ?1: user id
-- returns no rows, removes all arrivals for user_id
DELETE FROM arrivals WHERE user_id = ?1;
