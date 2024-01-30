import { Env } from "./env.js";
import { fetch } from "./fetch/index.js";
import { scheduled } from "./scheduled.js";

const handlers: ExportedHandler<Env> = {
  fetch,
  scheduled,
};

export default handlers;
