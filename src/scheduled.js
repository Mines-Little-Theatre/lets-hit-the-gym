import { ButtonStyleTypes, MessageComponentTypes } from "discord-interactions";
import { DateTime } from "luxon";
import { hourNames } from "./constants.js";
import {
  clearArrivals,
  getWeekday,
  getWorkout,
  updateScheduleMessageID,
} from "./queries.js";

const DISCORD_API = "https://discord.com/api/v10";

export async function scheduled(event, env) {
  const { weekday, hour } = DateTime.fromMillis(event.scheduledTime)
    .setZone(env.TIME_ZONE)
    .plus({ hours: 1 });
  const weekdayInfo = await getWeekday(env.DB, weekday);
  if (hour === weekdayInfo.post_hour) {
    await postSchedule(env, weekdayInfo);
  } else if (weekdayInfo.open_hour <= hour && hour < weekdayInfo.close_hour) {
    await postReminder(env);
  }
}

async function postSchedule(env, weekdayInfo) {
  const signupOptions = [];
  for (let i = weekdayInfo.open_hour; i < weekdayInfo.close_hour; i++) {
    signupOptions.push({
      label: hourNames[i],
      value: String(i),
    });
  }

  const messageSend = {
    content:
      "Ready to work out today? Let us know when you’re arriving so others can join you!",
    embeds: [],
    components: [
      {
        type: MessageComponentTypes.ACTION_ROW,
        components: [
          {
            type: MessageComponentTypes.STRING_SELECT,
            custom_id: "signup_selection",
            options: signupOptions,
            placeholder: "When are you working out today?",
          },
        ],
      },
      {
        type: MessageComponentTypes.ACTION_ROW,
        components: [
          {
            type: MessageComponentTypes.BUTTON,
            label: "Remove my signup",
            style: ButtonStyleTypes.SECONDARY,
            custom_id: "remove_signup",
          },
        ],
      },
    ],
  };

  const workout = await getWorkout(env.DB, weekdayInfo.workout_id);
  if (workout !== null) {
    messageSend.embeds.push({
      title: workout.title,
      description: workout.description,
      color: workout.color,
      fields: workout.routines.map((r) => ({
        name: r.title,
        value: r.description,
      })),
    });
  }

  messageSend.embeds.push({
    title: "Signups",
    description: "No one has signed up yet!",
    color: 0x5865f2,
  });

  const message = await (
    await fetch(`${DISCORD_API}/channels/${env.CHANNEL_ID}/messages`, {
      method: "POST",
      body: JSON.stringify(messageSend),
      headers: {
        Authorization: env.DISCORD_TOKEN,
        "Content-Type": "application/json",
      },
    })
  ).json();

  await clearArrivals(env.DB);
  await updateScheduleMessageID(env.DB, message.id);
}

async function postReminder() {}
