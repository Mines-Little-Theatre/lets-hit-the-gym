export interface Env {
  // vars
  readonly TIME_ZONE: string;

  // secrets
  readonly CHANNEL_ID: string;
  readonly DISCORD_APPLICATION_ID: string;
  readonly DISCORD_PUBLIC_KEY: string;
  readonly DISCORD_TOKEN: string;

  // bindings
  readonly DB: D1Database;
}
