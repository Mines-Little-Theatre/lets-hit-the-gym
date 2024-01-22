import {
  InteractionResponseFlags,
  InteractionResponseType,
  InteractionType,
  verifyKey,
} from "discord-interactions";
import { error } from "itty-router";

export async function interactions(request, env) {
  const signature = request.headers.get("x-signature-ed25519");
  const timestamp = request.headers.get("x-signature-timestamp");
  const body = await request.arrayBuffer();
  if (
    signature === null ||
    timestamp === null ||
    !verifyKey(body, signature, timestamp, env.DISCORD_PUBLIC_KEY)
  ) {
    return error(401, "bad request signature");
  }

  const interaction = JSON.parse(new TextDecoder("utf-8").decode(body));
  if (interaction.type === InteractionType.PING) {
    return {
      type: InteractionResponseType.PONG,
    };
  } else if (interaction.type === InteractionType.MESSAGE_COMPONENT) {
    return {
      type: InteractionResponseType.CHANNEL_MESSAGE_WITH_SOURCE,
      data: {
        content: interaction.data.custom_id,
        flags: InteractionResponseFlags.EPHEMERAL,
      },
    };
  } else {
    return error(400, "unexpected interaction type");
  }
}
