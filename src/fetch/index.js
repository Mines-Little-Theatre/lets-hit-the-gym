import { error, json, Router } from "itty-router";
import { interactions } from "./interactions.js";

const router = Router();

router.post("/interactions", interactions);

export async function fetch(...args) {
  return router
    .handle(...args)
    .then(json)
    .catch((e) => {
      console.error(e);
      return error(e);
    });
}
