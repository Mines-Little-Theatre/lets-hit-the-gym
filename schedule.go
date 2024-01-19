package main

import (
	"bytes"
	"fmt"
	"log"
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

	log.Println("opens at", day.OpenHour, "and closes at", day.CloseHour)

	channelID, err := store.GetChannelID()
	if err != nil {
		return fmt.Errorf("get channel ID: %w", err)
	}

	var buf bytes.Buffer
	err = templates.ExecuteTemplate(&buf, "schedule-message.txt", nil)
	if err != nil {
		return fmt.Errorf("execute template: %w", err)
	}
	content := strings.TrimSpace(buf.String())

	var embeds []*discordgo.MessageEmbed
	if day.Workout != nil {
		fields := make([]*discordgo.MessageEmbedField, 0, len(day.Workout.Routines))
		for _, routine := range day.Workout.Routines {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:  routine.Title,
				Value: routine.Description,
			})
		}
		embeds = []*discordgo.MessageEmbed{{
			Title:       day.Workout.Title,
			Description: day.Workout.Description,
			Color:       day.Workout.Color,
			Fields:      fields,
		}}
	}

	message, err := bot.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Content: content,
		Embeds:  embeds,
	})
	if err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	err = store.UpdateLastScheduleMessageID(message.ID)
	if err != nil {
		return fmt.Errorf("update last schedule message ID: %w", err)
	}

	return nil
}
