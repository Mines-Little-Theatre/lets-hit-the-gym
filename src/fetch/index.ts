import { error, json, Router } from "itty-router";
import { Env } from "../env.js";
import { interactions } from "./interactions.js";

const router = Router<Request, [Env, ExecutionContext]>();

router.post("/interactions", interactions);

export async function fetch(request: Request, env: Env, ctx: ExecutionContext) {
  return router
    .handle(request, env, ctx)
    .then(json)
    .catch((e) => {
      console.error(e);
      return error(
        500,
        typeof e === "object" && e !== null ? (e as object) : String(e),
      );
    });
}
