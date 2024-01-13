# Let's hit the gym!

A Discord bot that helps coordinate the Mines Lifting Thespians.

## Installation and Usage

Set the `GYM_BOT_DB` environment variable to a path where the application can create an SQLite database.

Configure the application:

```sh
go run . config --token 'Bot ...' --channel-id ...
```

Check for any issues:

```sh
go run . doctor
```

Then set it up to run on a schedule using cron!

### Example crontab

```
0 7 * * * <executable> schedule

50 13 * * * <executable> remind 2️⃣
50 14 * * * <executable> remind 3️⃣
50 15 * * * <executable> remind 4️⃣
50 16 * * * <executable> remind 5️⃣
```
