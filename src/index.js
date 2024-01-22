import {
  InteractionResponseFlags,
  InteractionResponseType,
  InteractionType,
  verifyKey,
} from "discord-interactions";

const exportedHandler = {
  async fetch(request, env) {
    if (request.method !== "POST") {
      return new Response(null, { status: 405, headers: { Allow: "POST" } });
    }

    const signature = request.headers.get("x-signature-ed25519");
    const timestamp = request.headers.get("x-signature-timestamp");
    const body = await request.arrayBuffer();
    if (
      signature === null ||
      timestamp === null ||
      !verifyKey(body, signature, timestamp, env.DISCORD_PUBLIC_KEY)
    ) {
      return new Response("Bad request signature.", { status: 401 });
    }

    const interaction = JSON.parse(new TextDecoder("utf-8").decode(body));
    if (interaction.type === InteractionType.PING) {
      return new Response(
        JSON.stringify(
          {
            type: InteractionResponseType.PONG,
          },
          { headers: { "Content-Type": "application/json" } },
        ),
      );
    } else if (interaction.type === InteractionType.MESSAGE_COMPONENT) {
      return new Response(
        JSON.stringify({
          type: InteractionResponseType.CHANNEL_MESSAGE_WITH_SOURCE,
          data: {
            content: interaction.data.custom_id,
            flags: InteractionResponseFlags.EPHEMERAL,
          },
        }),
        {
          headers: { "Content-Type": "application/json" },
        },
      );
    }
  },
  // async scheduled({ scheduledTime }, env) {},
};

export default exportedHandler;
