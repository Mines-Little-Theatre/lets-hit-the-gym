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
  get_hour_arrivals: `
    SELECT user_id FROM arrivals WHERE hour = ?1;
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
} as const;

type PreparedQueries = {
  -readonly [key in keyof typeof queries]?: D1PreparedStatement;
};

const preparedStatements = new WeakMap<D1Database, PreparedQueries>();

function prepareStatement(
  db: D1Database,
  name: keyof typeof queries,
): D1PreparedStatement {
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

async function getKV(db: D1Database, key: string): Promise<unknown> {
  const stmt = prepareStatement(db, "get_kv");
  return stmt.bind(key).first("value");
}

async function putKV(db: D1Database, key: string, value: unknown) {
  const stmt = prepareStatement(db, "put_kv");
  await stmt.bind(key, value).run();
}

export async function getScheduleMessageID(db: D1Database): Promise<string> {
  return getKV(db, "schedule_message_id") as Promise<string>;
}

export async function updateScheduleMessageID(
  db: D1Database,
  messageID: string,
) {
  await putKV(db, "schedule_message_id", messageID);
}

export interface HourArrivals {
  hour: number;
  users: string[];
}

export async function getAllArrivals(db: D1Database): Promise<HourArrivals[]> {
  const result = [];
  let currentHour: HourArrivals | undefined;
  const stmt = prepareStatement(db, "get_arrivals");
  for (const row of (await stmt.all()).results) {
    if (!currentHour || currentHour.hour !== row["hour"]) {
      result.push((currentHour = { hour: row["hour"] as number, users: [] }));
    }
    currentHour.users.push(row["user_id"] as string);
  }
  return result;
}

export async function getHourArrivals(
  db: D1Database,
  hour: number,
): Promise<string[]> {
  const stmt = prepareStatement(db, "get_hour_arrivals");
  return (await stmt.bind(hour).all()).results.map(
    (row) => row["user_id"] as string,
  );
}

export async function setUserArrivals(
  db: D1Database,
  userID: string,
  hours: readonly number[] | null,
) {
  const clearStmt = prepareStatement(db, "clear_user_arrivals");
  await clearStmt.bind(userID).run();
  if (hours && hours.length > 0) {
    const addStmt = prepareStatement(db, "add_arrival");
    for (const hour of hours) {
      await addStmt.bind(hour, userID).run();
    }
  }
}

export async function clearArrivals(db: D1Database) {
  const stmt = prepareStatement(db, "clear_arrivals");
  await stmt.run();
}

export interface Weekday {
  post_hour: number;
  open_hour: number;
  close_hour: number;
  workout_id: number | null;
}

export async function getWeekday(
  db: D1Database,
  id: number,
): Promise<Weekday | null> {
  const stmt = prepareStatement(db, "get_weekday");
  return stmt.bind(id).first();
}

export interface Workout {
  title: string;
  description: string;
  color: number;
  routines: {
    title: string;
    description: string;
  }[];
}

export async function getWorkout(
  db: D1Database,
  id: number,
): Promise<Workout | null> {
  const stmt = prepareStatement(db, "get_workout");
  const { results } = await stmt.bind(id).all();
  if (results.length <= 0) {
    return null;
  }
  const workout: Workout = {
    title: (results[0]?.["workout_title"] as string) ?? "",
    description: (results[0]?.["workout_description"] as string) ?? "",
    color: (results[0]?.["workout_color"] as number) ?? 0,
    routines: [],
  };
  if (!isNullOrUndefined(results[0]?.["routine_title"])) {
    for (const row of results) {
      workout.routines.push({
        title: row["routine_title"] as string,
        description: row["routine_description"] as string,
      });
    }
  }
  return workout;
}

function isNullOrUndefined(v: unknown): v is null | undefined {
  return v === null || v === undefined;
}
