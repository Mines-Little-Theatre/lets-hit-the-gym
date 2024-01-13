package main

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/Mines-Little-Theatre/lets-hit-the-gym/store"
	"github.com/alecthomas/kong"
	"github.com/bwmarrin/discordgo"
)

var (
	//go:embed templates/*.txt
	templateFS embed.FS
	templates  = template.Must(template.ParseFS(templateFS, "templates/*.txt"))
)

type CommandLine struct {
	Config   ConfigCmd   `cmd:"" help:"Edit configuration"`
	Doctor   DoctorCmd   `cmd:"" help:"Check for issues with configuration, permissions, etc."`
	Remind   RemindCmd   `cmd:"" help:"Post a reminder message"`
	Schedule ScheduleCmd `cmd:"" help:"Post a schedule message"`
}

func main() {
	dataSourceName, ok := os.LookupEnv("GYM_BOT_DB")
	if !ok {
		log.Fatalln("GYM_BOT_DB environment variable must be set")
	}
	store, err := store.Open(dataSourceName)
	if err != nil {
		log.Fatalln("connect:", err)
	}
	k := kong.Parse(new(CommandLine))
	err = k.Run(store)
	store.Close()
	k.FatalIfErrorf(err)
}

func connectBot(store *store.Store) (*discordgo.Session, error) {
	token, err := store.GetToken()
	if err == sql.ErrNoRows {
		return nil, errors.New("token not found")
	} else if err != nil {
		return nil, fmt.Errorf("get token: %w", err)
	}

	bot, err := discordgo.New(token)
	if err != nil {
		return nil, fmt.Errorf("create bot session: %w", err)
	}

	return bot, nil
}
