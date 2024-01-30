import {
  APIMessage,
  ButtonStyle,
  ComponentType,
  RESTPostAPIChannelMessageJSONBody,
} from "discord-api-types/v10";
import { DateTime } from "luxon";
import { hourNames } from "./constants.js";
import { Env } from "./env.js";
import {
  Weekday,
  clearArrivals,
  getHourArrivals,
  getScheduleMessageID,
  getWeekday,
  getWorkout,
  updateScheduleMessageID,
} from "./queries.js";

const DISCORD_API = "https://discord.com/api/v10";

export async function scheduled(event: ScheduledController, env: Env) {
  const { weekday, hour } = DateTime.fromMillis(event.scheduledTime)
    .setZone(env.TIME_ZONE)
    .plus({ hours: 1 });
  const weekdayInfo = await getWeekday(env.DB, weekday);
  if (weekdayInfo) {
    if (hour === weekdayInfo.post_hour) {
      await postSchedule(env, weekdayInfo);
    } else if (weekdayInfo.open_hour <= hour && hour < weekdayInfo.close_hour) {
      await postReminder(env, hour);
    }
  }
}

async function postSchedule(env: Env, weekdayInfo: Weekday) {
  const signupOptions = [];
  for (let i = weekdayInfo.open_hour; i < weekdayInfo.close_hour; i++) {
    signupOptions.push({
      label: hourNames[i] ?? "undefined",
      value: String(i),
    });
  }

  const messageSend: RESTPostAPIChannelMessageJSONBody = {
    content:
      "Ready to work out today? Let us know when you’re arriving so others can join you!",
    embeds: [],
    components: [
      {
        type: ComponentType.ActionRow,
        components: [
          {
            type: ComponentType.StringSelect,
            custom_id: "signup_selection",
            options: signupOptions,
            placeholder: "When are you working out today?",
          },
        ],
      },
      {
        type: ComponentType.ActionRow,
        components: [
          {
            type: ComponentType.Button,
            label: "Remove my signup",
            style: ButtonStyle.Secondary,
            custom_id: "remove_signup",
          },
        ],
      },
    ],
  };

  if (weekdayInfo.workout_id !== null) {
    const workout = await getWorkout(env.DB, weekdayInfo.workout_id);
    if (workout !== null) {
      messageSend.embeds!.push({
        title: workout.title,
        description: workout.description,
        color: workout.color,
        fields: workout.routines.map((r) => ({
          name: r.title,
          value: r.description,
        })),
      });
    }
  }

  messageSend.embeds!.push({
    title: "Signups",
    description: "No one has signed up yet!",
    color: 0x5865f2,
  });

  const message = (await (
    await fetch(`${DISCORD_API}/channels/${env.CHANNEL_ID}/messages`, {
      method: "POST",
      body: JSON.stringify(messageSend),
      headers: {
        Authorization: env.DISCORD_TOKEN,
        "Content-Type": "application/json",
      },
    })
  ).json()) as APIMessage;

  await clearArrivals(env.DB);
  await updateScheduleMessageID(env.DB, message.id);
}

async function postReminder(env: Env, hour: number) {
  const arrivingUsers = await getHourArrivals(env.DB, hour);
  if (arrivingUsers.length > 0) {
    const scheduleMessageID = await getScheduleMessageID(env.DB);
    const messageSend: RESTPostAPIChannelMessageJSONBody = {
      content:
        "Looks like we’ve got some people headed for the gym!\n- <@" +
        arrivingUsers.join(">\n- <@") +
        ">",
      allowed_mentions: {
        replied_user: false,
      },
      message_reference: {
        message_id: scheduleMessageID,
      },
    };
    await fetch(`${DISCORD_API}/channels/${env.CHANNEL_ID}/messages`, {
      method: "POST",
      body: JSON.stringify(messageSend),
      headers: {
        Authorization: env.DISCORD_TOKEN,
        "Content-Type": "application/json",
      },
    });
  }
}
