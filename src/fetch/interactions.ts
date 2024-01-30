import {
  APIEmbed,
  APIInteraction,
  APIInteractionResponse,
  APIMessageComponentInteraction,
  ComponentType,
  InteractionResponseType,
  InteractionType,
  MessageFlags,
} from "discord-api-types/v10";
import { verifyKey } from "discord-interactions";
import { error } from "itty-router";
import { hourNames } from "../constants.js";
import { Env } from "../env.js";
import {
  getAllArrivals,
  getScheduleMessageID,
  setUserArrivals,
} from "../queries.js";

export async function interactions(
  request: Request,
  env: Env,
): Promise<Response | APIInteractionResponse> {
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

  const interaction: APIInteraction = JSON.parse(
    new TextDecoder("utf-8").decode(body),
  ) as APIInteraction;
  if (interaction.type === InteractionType.Ping) {
    return {
      type: InteractionResponseType.Pong,
    };
  } else if (interaction.type === InteractionType.MessageComponent) {
    const scheduleMessageID = await getScheduleMessageID(env.DB);
    if (interaction.message.id !== scheduleMessageID) {
      return {
        type: InteractionResponseType.ChannelMessageWithSource,
        data: {
          content: `You can’t change your schedule in the past! Try using today’s signup: https://discord.com/channels/${interaction.guild_id}/${env.CHANNEL_ID}/${scheduleMessageID}`,
          flags: MessageFlags.Ephemeral,
        },
      };
    } else {
      switch (interaction.data.custom_id) {
        case "signup_selection":
          return signupSelection(env, interaction);
        case "remove_signup":
          return removeSignup(env, interaction);
        default:
          return {
            type: InteractionResponseType.ChannelMessageWithSource,
            data: {
              content: `Looks like you somehow interacted with an invalid component \`${interaction.data.custom_id}\`.`,
              flags: MessageFlags.Ephemeral,
            },
          };
      }
    }
  } else {
    return error(400, "unexpected interaction type");
  }
}

async function signupSelection(
  env: Env,
  interaction: APIMessageComponentInteraction,
): Promise<Response | APIInteractionResponse> {
  if (interaction.data.component_type === ComponentType.StringSelect) {
    await setUserArrivals(
      env.DB,
      interaction.member!.user.id,
      interaction.data.values.map((v) => Number.parseInt(v)),
    );
    return modifySignupEmbed(env, interaction);
  } else {
    return error(400, "unexpected component type");
  }
}

async function removeSignup(
  env: Env,
  interaction: APIMessageComponentInteraction,
): Promise<Response | APIInteractionResponse> {
  await setUserArrivals(env.DB, interaction.member!.user.id, null);
  return modifySignupEmbed(env, interaction);
}

async function modifySignupEmbed(
  env: Env,
  interaction: APIMessageComponentInteraction,
): Promise<Response | APIInteractionResponse> {
  const arrivals = await getAllArrivals(env.DB);
  const signupEmbed: APIEmbed = {
    title: "Signups",
    color: 0x5865f2,
  };
  if (arrivals.length === 0) {
    signupEmbed.description = "No one has signed up yet!";
  } else {
    signupEmbed.fields = arrivals.map((hour) => ({
      name: hourNames[hour.hour] ?? "undefined",
      value: "<@" + hour.users.join(">\n<@") + ">",
      inline: true,
    }));
  }
  const embeds = interaction.message.embeds;
  const signupEmbedIndex = embeds.findIndex((e) => e.title === "Signups");
  if (signupEmbedIndex === -1) {
    embeds.push(signupEmbed);
  } else {
    embeds[signupEmbedIndex] = signupEmbed;
  }
  return {
    type: InteractionResponseType.UpdateMessage,
    data: {
      content: interaction.message.content,
      embeds,
      components: interaction.message.components,
    },
  };
}
