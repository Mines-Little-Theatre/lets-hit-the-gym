package main

import (
	"bytes"
	"fmt"
	"strings"
)

type ScheduleCmd struct {
}

func (*ScheduleCmd) Run(conn *Connections) error {
	channelID, err := conn.Store.GetChannelID()
	if err != nil {
		return fmt.Errorf("get channel ID: %w", err)
	}

	var buf bytes.Buffer
	err = templates.ExecuteTemplate(&buf, "schedule-message.txt", nil)
	if err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	message, err := conn.Bot.ChannelMessageSend(channelID, strings.TrimSpace(buf.String()))
	if err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	err = conn.Store.UpdateLastScheduleMessageID(message.ID)
	if err != nil {
		return fmt.Errorf("update last schedule message ID: %w", err)
	}

	return nil
}
