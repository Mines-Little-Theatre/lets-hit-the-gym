package main

import (
	"errors"

	"github.com/Mines-Little-Theatre/lets-hit-the-gym/store"
	"github.com/Mines-Little-Theatre/lets-hit-the-gym/util"
	"github.com/bwmarrin/discordgo"
)

type Connections struct {
	Bot   *discordgo.Session
	Store *store.Store
}

func Connect() (*Connections, error) {
	botToken := util.ReadEnv("GYM_TOKEN")
	dbName := util.ReadEnv("GYM_DB")

	bot, err := discordgo.New(botToken)
	if err != nil {
		return nil, err
	}

	store, err := store.Open(dbName)
	if err != nil {
		bot.Close()
		return nil, err
	}

	return &Connections{Bot: bot, Store: store}, nil
}

func (c *Connections) Close() error {
	return errors.Join(
		c.Bot.Close(),
		c.Store.Close(),
	)
}
