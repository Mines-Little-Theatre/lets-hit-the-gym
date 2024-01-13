package main

import (
	"database/sql"
	"errors"
	"log"

	"github.com/Mines-Little-Theatre/lets-hit-the-gym/store"
	"github.com/Mines-Little-Theatre/lets-hit-the-gym/util"
	"github.com/bwmarrin/discordgo"
)

type Connections struct {
	Bot   *discordgo.Session
	Store *store.Store
}

func Connect() (*Connections, error) {
	store, err := store.Open(util.ReadEnv("GYM_BOT_DB"))
	if err != nil {
		return nil, err
	}

	botToken, err := store.GetToken()
	if err != nil {
		store.Close()
		if err == sql.ErrNoRows {
			log.Println("token not found")
		}
		return nil, err
	}

	bot, err := discordgo.New(botToken)
	if err != nil {
		store.Close()
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
