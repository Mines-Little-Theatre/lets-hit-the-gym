package main

import (
	"embed"
	"log"
	"text/template"

	"github.com/alecthomas/kong"
)

var (
	//go:embed templates/*.txt
	templateFS embed.FS
	templates  = template.Must(template.ParseFS(templateFS, "templates/*.txt"))
)

type CommandLine struct {
	Schedule ScheduleCmd `cmd:"" help:"Post a schedule message"`
	Remind   RemindCmd   `cmd:"" help:"Post a reminder message"`
}

func main() {
	k := kong.Parse(new(CommandLine))
	conn, err := Connect()
	if err != nil {
		log.Fatalln("connect:", err)
	}
	err = k.Run(conn)
	conn.Close()
	k.FatalIfErrorf(err)
}
