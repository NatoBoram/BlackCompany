package main

import (
	"log"

	"github.com/aiseeq/s2l/protocol/api"
	"github.com/aiseeq/s2l/protocol/client"
)

type Bot struct {
	c *client.Client

	observation *api.ResponseObservation
}

func (b *Bot) Step() {
}

func (b *Bot) Observe() {
	obs, err := b.c.Observation(api.RequestObservation{})
	if err != nil {
		log.Printf("Failed to observe: %v", err)
	}
	b.observation = obs
}
