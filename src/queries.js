const queries = {
  add_arrival: `
    INSERT INTO arrivals (hour, user_id) VALUES (?1, ?2);
  `,
  clear_arrivals: `
    DELETE FROM arrivals;
  `,
  clear_user_arrivals: `
    DELETE FROM arrivals WHERE user_id = ?1;
  `,
  get_arrivals: `
    SELECT hour, user_id FROM arrivals ORDER BY hour;
  `,
  get_kv: `
    SELECT value FROM kv_store WHERE key = ?1;
  `,
  get_weekday: `
    SELECT post_hour, open_hour, close_hour, workout_id
    FROM weekdays WHERE id = ?1;
  `,
  get_workout: `
    SELECT w.title AS workout_title,
      w.description AS workout_description,
      w.color AS workout_color,
      r.title AS routine_title,
      r.description AS routine_description
    FROM workouts AS w
      LEFT JOIN workout_routines AS x ON w.id = x.workout_id
      LEFT JOIN routines AS r ON x.routine_id = r.id
    WHERE w.id = ?1
    ORDER BY x.ordinal;
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

export async function updateScheduleMessageID(db, messageID) {
  await putKV(db, "schedule_message_id", messageID);
}

export async function getAllArrivals(db) {
  const result = [];
  let currentHour = null;
  const stmt = prepareStatement(db, "get_arrivals");
  for (const row of (await stmt.all()).results) {
    if (!currentHour || currentHour.hour !== row.hour) {
      result.push((currentHour = { hour: row.hour, users: [] }));
    }
    currentHour.users.push(row.user_id);
  }
  return result;
}

export async function setUserArrivals(db, userID, hours) {
  const clearStmt = prepareStatement(db, "clear_user_arrivals");
  await clearStmt.bind(userID).run();
  if (hours && hours.length > 0) {
    const addStmt = prepareStatement(db, "add_arrival");
    for (const hour of hours) {
      await addStmt.bind(hour, userID).run();
    }
  }
}

export async function clearArrivals(db) {
  const stmt = prepareStatement(db, "clear_arrivals");
  await stmt.run();
}

export async function getWeekday(db, id) {
  const stmt = prepareStatement(db, "get_weekday");
  return stmt.bind(id).first();
}

export async function getWorkout(db, id) {
  const stmt = prepareStatement(db, "get_workout");
  const { results } = await stmt.bind(id).all();
  if (results.length <= 0) {
    return null;
  }
  const workout = {
    title: results[0].workout_title,
    description: results[0].workout_description,
    color: results[0].workout_color,
    routines: [],
  };
  if (results[0].routine_title !== null) {
    for (const row of results) {
      workout.routines.push({
        title: row.routine_title,
        description: row.routine_description,
      });
    }
  }
  return workout;
}
