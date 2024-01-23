import {
  InteractionResponseFlags,
  InteractionResponseType,
  InteractionType,
  verifyKey,
} from "discord-interactions";
import { error } from "itty-router";
import { getScheduleMessageID } from "../queries.js";

export async function interactions(request, env) {
  const signature = request.headers.get("x-signature-ed25519");
  const timestamp = request.headers.get("x-signature-timestamp");
  const body = await request.arrayBuffer();
  if (
    signature === null ||
    timestamp === null ||
    !verifyKey(body, signature, timestamp, env.DISCORD_PUBLIC_KEY)
  ) {
    return error(401);
  }

  const interaction = JSON.parse(new TextDecoder("utf-8").decode(body));
  if (interaction.type === InteractionType.PING) {
    return {
      type: InteractionResponseType.PONG,
    };
  } else if (interaction.type === InteractionType.MESSAGE_COMPONENT) {
    const scheduleMessageID = await getScheduleMessageID(env.DB);
    if (interaction.message.id !== scheduleMessageID) {
      return {
        type: InteractionResponseType.CHANNEL_MESSAGE_WITH_SOURCE,
        data: {
          content: `You can’t change your schedule in the past! Try using today’s signup: https://discord.com/channels/${interaction.guild_id}/${env.CHANNEL_ID}/${scheduleMessageID}`,
          flags: InteractionResponseFlags.EPHEMERAL,
        },
      };
    } else {
      return error(501, "not implemented");
    }
  } else {
    return error(400, "unexpected interaction type");
  }
}
