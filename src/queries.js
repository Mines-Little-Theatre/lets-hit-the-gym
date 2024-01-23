const queries = {
  get_kv: `
    SELECT value FROM kv_store WHERE key = ?1
  `,
  put_kv: `
    INSERT INTO kv_store (key, value) VALUES (?1, ?2)
      ON CONFLICT (key) DO UPDATE SET value = excluded.value;
  `,
};

const preparedStatements = new WeakMap();

function prepareStatement(db, name) {
  let dbPrepared = preparedStatements.get(db);
  if (!dbPrepared) {
    preparedStatements.set(db, (dbPrepared = {}));
  }

  let stmt = dbPrepared[name];
  if (!stmt) {
    stmt = dbPrepared[name] = db.prepare(queries[name]);
  }
  return stmt;
}

async function getKV(db, key) {
  const stmt = prepareStatement(db, "get_kv");
  return stmt.bind(key).first("value");
}

async function putKV(db, key, value) {
  const stmt = prepareStatement(db, "put_kv");
  await stmt.bind(key, value).run();
}

export async function getScheduleMessageID(db) {
  return getKV(db, "schedule_message_id");
}

export async function updateScheduleMessageID(db, message_id) {
  await putKV(db, "schedule_message_id", message_id);
}
