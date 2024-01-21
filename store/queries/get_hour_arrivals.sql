-- ?1: hour
-- returns any number of rows (text): user IDs
SELECT user_id FROM arrivals WHERE hour = ?1;
