package main

import (
	"bytes"
	"embed"
	"log"
	"os"
	"strings"
	"text/template"

	gym "github.com/Mines-Little-Theatre/lets-hit-the-gym"
)

var (
	//go:embed templates/*.txt
	templateFS embed.FS
	templates  = template.Must(template.ParseFS(templateFS, "templates/*.txt"))
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("specify reaction in command line args (unicode emoji or name:id)")
	}
	targetReaction := os.Args[1]

	conn, err := gym.Connect()
	if err != nil {
		log.Fatalln("connect:", err)
	}

	channelID, err := conn.Store.GetChannelID()
	if err != nil {
		conn.Close()
		log.Fatalln("get channel ID:", err)
	}

	lastScheduleMessageID, err := conn.Store.GetLastScheduleMessageID()
	if err != nil {
		conn.Close()
		log.Fatalln("get last schedule message ID:", err)
	}

	userMentions := make([]string, 0)
	afterID := ""
	for {
		users, err := conn.Bot.MessageReactions(channelID, lastScheduleMessageID, targetReaction, 100, "", afterID)
		if err != nil {
			conn.Close()
			log.Fatalln("get reactions:", err)
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
		templates.ExecuteTemplate(&buf, "message.txt", userMentions)
		_, err := conn.Bot.ChannelMessageSend(channelID, strings.TrimSpace(buf.String()))
		if err != nil {
			conn.Close()
			log.Fatalln("send message:", err)
		}
	}

	conn.Close()
}
