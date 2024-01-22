import { Env } from "./env.js";

const exportedHandler: ExportedHandler<Env> = {
  async fetch(request, env): Promise<Response> {},
  async scheduled({ scheduledTime }, env) {},
};

export default exportedHandler;
