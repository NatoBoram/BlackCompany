package main

import (
	"log"

	"github.com/aiseeq/s2l/protocol/api"
	"github.com/aiseeq/s2l/protocol/client"
)

// Bot holds the state of the bot.
type Bot struct {
	c *client.Client

	observation *api.ResponseObservation
}

// Step is called at every step of the game. This is the main loop of the bot.
func (b *Bot) Step() {
}

// Observe fetches the current observation from the game.
func (b *Bot) Observe() {
	obs, err := b.c.Observation(api.RequestObservation{})
	if err != nil {
		log.Printf("Failed to observe: %v", err)
	}
	b.observation = obs
}
