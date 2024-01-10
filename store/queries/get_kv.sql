-- ?1: key
-- returns zero or one row (any?): value

SELECT value FROM kv_store WHERE key = ?1;
