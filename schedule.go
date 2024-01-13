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
	Workout string `arg:"" optional:"" name:"workout" help:"Name of the workout card to display"`
}

func (c *ScheduleCmd) Run(store *store.Store) error {
	bot, err := connectBot(store)
	if err != nil {
		return err
	}

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
	if c.Workout != "" {
		workout, err := store.GetWorkout(c.Workout)
		if err != nil {
			log.Println("get workout:", err)
		} else {
			fields := make([]*discordgo.MessageEmbedField, 0, len(workout.Routines))
			for _, routine := range workout.Routines {
				fields = append(fields, &discordgo.MessageEmbedField{
					Name:  routine.Title,
					Value: routine.Description,
				})
			}
			embeds = []*discordgo.MessageEmbed{{
				Title:       workout.Title,
				Description: workout.Description,
				Color:       workout.Color,
				Fields:      fields,
			}}
		}
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
