# Let's hit the gym!

A Discord bot that helps coordinate the Mines Lifting Thespians.

## Installation and Usage

You will need to set two environment variables:

- `GYM_BOT_TOKEN`: your authorization token, starting with "Bot"
- `GYM_BOT_DB`: a path to an SQLite database

Run the `schedule` program to create the database:

```sh
go run ./bin/schedule
```

It will exit with an error because the newly-created database does not have a channel ID set. Use the `sqlite3` program to connect to the database and set the channel ID:

```sql
INSERT INTO kv_store (key, value) VALUES ('channel_id', '{{ your channel ID }}');
```

Then set it up to run on a schedule using cron!

### Suggested crontab

```
0 12 * * MON-FRI .../schedule
50 13 * * MON-FRI .../remind 2️⃣
50 14 * * MON-FRI .../remind 3️⃣
50 15 * * MON-FRI .../remind 4️⃣
50 16 * * MON-FRI .../remind 5️⃣
```
