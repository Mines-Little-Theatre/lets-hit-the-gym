-- ?1: key
-- ?2: value
-- returns no rows, upserts (key, value) into kv_store

INSERT INTO kv_store (key, value) VALUES (?1, ?2)
  ON CONFLICT (key) DO UPDATE SET value = excluded.value;
