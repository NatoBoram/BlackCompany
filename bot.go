package main

import (
	"log"

	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/api"
)

// Bot holds the state of the bot.
type Bot struct {
	*scl.Bot

	observation *api.ResponseObservation
}

// Step is called at every step of the game. This is the main loop of the bot.
func (b *Bot) Step() {
	b.BuildWorker()
}

// Expand expands the bot's base whenever enough resources are available.
func (b *Bot) Expand() {
}

// Observe fetches the current observation from the game.
func (b *Bot) Observe() {
	obs, err := b.Client.Observation(api.RequestObservation{})
	if err != nil {
		log.Printf("Failed to observe: %v", err)
	}
	b.observation = obs
}

func OnUnitCreated(unit *scl.Unit) {
}
