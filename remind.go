package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type RemindCmd struct {
	Emoji string `arg:"" name:"emoji" help:"Reaction to check for (unicode emoji or name:id)"`
}

func (c *RemindCmd) Run(conn *Connections) error {
	channelID, err := conn.Store.GetChannelID()
	if err != nil {
		return fmt.Errorf("get channel ID: %w", err)
	}

	lastScheduleMessageID, err := conn.Store.GetLastScheduleMessageID()
	if err != nil {
		return fmt.Errorf("get last schedule message ID: %w", err)
	}

	userMentions := make([]string, 0)
	afterID := ""
	for {
		users, err := conn.Bot.MessageReactions(channelID, lastScheduleMessageID, c.Emoji, 100, "", afterID)
		if err != nil {
			return fmt.Errorf("get reactions: %w", err)
		}
		for _, u := range users {
			userMentions = append(userMentions, u.Mention())
		}
		if len(users) < 100 {
			break
		} else if len(users) > 0 {
			afterID = users[len(users)-1].ID
		}
	}
	if len(userMentions) > 0 {
		var buf bytes.Buffer
		err := templates.ExecuteTemplate(&buf, "remind-message.txt", userMentions)
		if err != nil {
			return fmt.Errorf("execute template: %w", err)
		}
		_, err = conn.Bot.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
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
