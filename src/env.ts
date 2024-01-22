export interface Env {
  // vars
  TIME_ZONE: string;

  // secrets
  DISCORD_TOKEN: string;
  DISCORD_PUBLIC_KEY: string;
  DISCORD_APPLICATION_ID: string;
  CHANNEL_ID: string;

  // bindings
  DB: D1Database;
}
