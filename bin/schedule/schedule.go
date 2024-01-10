package main

import (
	_ "embed"
	"log"
	"strings"

	gym "github.com/Mines-Little-Theatre/lets-hit-the-gym"
)

//go:embed message.txt
var messageContent string

func main() {
	conn, err := gym.Connect()
	if err != nil {
		log.Fatalln("connect:", err)
	}

	channelID, err := conn.Store.GetChannelID()
	if err != nil {
		conn.Close()
		log.Fatalln("get channel ID:", err)
	}

	message, err := conn.Bot.ChannelMessageSend(channelID, strings.TrimSpace(messageContent))
	if err != nil {
		conn.Close()
		log.Fatalln("send message:", err)
	}

	err = conn.Store.UpdateLastScheduleMessageID(message.ID)
	if err != nil {
		conn.Close()
		log.Fatalln("update last schedule message ID:", err)
	}

	conn.Close()
}
