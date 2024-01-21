package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Mines-Little-Theatre/lets-hit-the-gym/store"
	"github.com/bwmarrin/discordgo"
)

type ScheduleCmd struct {
	Day string `arg:"" name:"day" help:"The name of the day"`
}

func (c *ScheduleCmd) Run(store *store.Store) error {
	bot, err := connectBot(store)
	if err != nil {
		return err
	}

	day, err := store.GetDay(c.Day)
	if err != nil {
		return fmt.Errorf("get day: %w", err)
	} else if day == nil {
		return fmt.Errorf("invalid day: %s", c.Day)
	}

	if day.OpenHour >= len(hourNames) || day.CloseHour > len(hourNames) || day.CloseHour <= day.OpenHour {
		return fmt.Errorf("day %s has invalid hours %d-%d (max is %d-%d)", c.Day, day.OpenHour, day.CloseHour, len(hourNames)-1, len(hourNames))
	}

	channelID, err := store.GetChannelID()
	if err != nil {
		return fmt.Errorf("get channel ID: %w", err)
	}

	messageSend := new(discordgo.MessageSend)

	var buf bytes.Buffer
	err = templates.ExecuteTemplate(&buf, "schedule-message.txt", nil)
	if err != nil {
		return fmt.Errorf("execute template: %w", err)
	}
	messageSend.Content = strings.TrimSpace(buf.String())

	if day.Workout != nil {
		fields := make([]*discordgo.MessageEmbedField, 0, len(day.Workout.Routines))
		for _, routine := range day.Workout.Routines {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:  routine.Title,
				Value: routine.Description,
			})
		}
		messageSend.Embeds = append(messageSend.Embeds, &discordgo.MessageEmbed{
			Title:       day.Workout.Title,
			Description: day.Workout.Description,
			Color:       day.Workout.Color,
			Fields:      fields,
		})
	}

	signupEmbedIndex := len(messageSend.Embeds)
	messageSend.Embeds = append(messageSend.Embeds, &discordgo.MessageEmbed{
		Title:       "Signups",
		Description: "No one has signed up yet!",
		Color:       0x5865f2,
	})

	options := make([]discordgo.SelectMenuOption, 0, day.CloseHour-day.OpenHour+1)
	for i := day.OpenHour; i < day.CloseHour; i++ {
		options = append(options, discordgo.SelectMenuOption{
			Label: hourNames[i],
			Value: fmt.Sprint(i),
		})
	}
	messageSend.Components = []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					MenuType:    discordgo.StringSelectMenu,
					CustomID:    "signup_selection",
					Options:     options,
					Placeholder: "When are you working out today?",
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Remove my signup",
					Style:    discordgo.SecondaryButton,
					CustomID: "remove_signup",
				},
			},
		},
	}

	message, err := bot.ChannelMessageSendComplex(channelID, messageSend)
	if err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	err = store.UpdateLastScheduleMessage(message.ID, signupEmbedIndex)
	if err != nil {
		return fmt.Errorf("update last schedule message ID: %w", err)
	}

	return nil
}
