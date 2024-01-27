# Let's hit the gym!

A Discord bot that helps coordinate the Mines Lifting Thespians. Runs on Cloudflare Workers.

## Deployment

```sh
npm i # install dependencies, including "wrangler" tool
npx wrangler deploy # deploy to Cloudflare Workers
npx wrangler d1 migrations apply # create tables in the database
npx wrangler d1 execute prod-gym-bot --file data/src_hours.sql # populate the database with weekday info

# fill in with your Discord application's info
npx wrangler secret put DISCORD_APPLICATION_ID
npx wrangler secret put DISCORD_PUBLIC_KEY
npx wrangler secret put DISCORD_TOKEN # should be of the form "Bot ..."
npx wrangler secret put CHANNEL_ID # the channel to post in
```

Make sure the bot user is added to the server and has the "View Channel," "Send Messages," "Embed Links," and "Read Message History" permissions in the channel!
