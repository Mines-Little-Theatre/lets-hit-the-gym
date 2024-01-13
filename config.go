package main

import (
	"errors"
	"fmt"

	"github.com/Mines-Little-Theatre/lets-hit-the-gym/store"
)

type ConfigCmd struct {
	Token     string
	ChannelID string
}

func (c *ConfigCmd) Run(store *store.Store) error {
	var errs []error

	if c.Token != "" {
		err := store.UpdateToken(c.Token)
		if err != nil {
			errs = append(errs, fmt.Errorf("update token: %w", err))
		}
	}

	if c.ChannelID != "" {
		err := store.UpdateChannelID(c.Token)
		if err != nil {
			errs = append(errs, fmt.Errorf("update channel id: %w", err))
		}
	}

	return errors.Join(errs...)
}
