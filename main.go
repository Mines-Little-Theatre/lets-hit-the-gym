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
	Doctor   DoctorCmd   `cmd:"" help:"Check the configuration for potential issues"`
	Remind   RemindCmd   `cmd:"" help:"Post a reminder message"`
	Schedule ScheduleCmd `cmd:"" help:"Post a schedule message"`
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
