package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Mines-Little-Theatre/lets-hit-the-gym/store"
	"github.com/bwmarrin/discordgo"
)

type RemindCmd struct {
	Hour int `arg:"" name:"hour" help:"Hour number to remind about signups"`
}

func (c *RemindCmd) Run(store *store.Store) error {
	bot, err := connectBot(store)
	if err != nil {
		return err
	}

	channelID, err := store.GetChannelID()
	if err != nil {
		return fmt.Errorf("get channel ID: %w", err)
	}

	userIDs, err := store.GetArrivingUsers(c.Hour)
	if err != nil {
		return fmt.Errorf("get arriving users: %w", err)
	}

	if len(userIDs) > 0 {
		lastScheduleMessageID, err := store.GetLastScheduleMessageID()
		if err != nil {
			return fmt.Errorf("get last schedule message ID: %w", err)
		}

		var buf bytes.Buffer
		err = templates.ExecuteTemplate(&buf, "remind-message.txt", userIDs)
		if err != nil {
			return fmt.Errorf("execute template: %w", err)
		}
		_, err = bot.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
			Content: strings.TrimSpace(buf.String()),
			AllowedMentions: &discordgo.MessageAllowedMentions{
				Parse:       []discordgo.AllowedMentionType{discordgo.AllowedMentionTypeUsers},
				RepliedUser: false,
			},
			Reference: &discordgo.MessageReference{MessageID: lastScheduleMessageID},
		})
		if err != nil {
			return fmt.Errorf("send message: %w", err)
		}
	}

	return nil
}
