package main

import (
	"log"

	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/api"
)

// Bot holds the state of the bot.
type Bot struct {
	*scl.Bot
}

// Step is called at every step of the game. This is the main loop of the bot.
func (b *Bot) Step() {
	b.Cmds = &scl.CommandsStack{}
	b.Loop = int(b.Obs.GameLoop)

	// Skip repeated frames
	if b.Loop < b.LastLoop+b.FramesPerOrder {
		return
	} else {
		b.LastLoop = b.Loop
	}

	b.ParseData()

	b.BuildWorker()
	b.BuildSupplyDepot()

	b.Cmds.Process(&b.Actions)
	if len(b.Actions) > 0 {
		log.Printf("Loop %d, actions %d, %v", b.Loop, len(b.Actions), b.Actions)

		if _, err := b.Client.Action(api.RequestAction{Actions: b.Actions}); err != nil {
			log.Printf("Failed to send actions: %v", err)
		}

		b.Actions = nil
	}
}

// Expand expands the bot's base whenever enough resources are available.
func (b *Bot) Expand() {
}

// Observe fetches the current observation from the game.
func (b *Bot) Observe() {
	o, err := b.Client.Observation(api.RequestObservation{})
	if err != nil {
		log.Printf("Failed to observe: %v", err)
		return
	}

	b.Obs = o.Observation
	b.Chat = o.Chat
	b.Result = o.PlayerResult
	b.Errors = o.ActionErrors
}

func OnUnitCreated(unit *scl.Unit) {
}

func (b *Bot) ParseData() {
	if b.Info == nil {
		log.Printf("Info is nil")
		return
	}
	if b.Obs == nil {
		log.Printf("Observation is nil")
		return
	}
	if b.Obs.RawData == nil {
		log.Printf("RawData is nil")
		return
	}
	if b.Obs.RawData.MapState == nil {
		log.Printf("MapState is nil")
		return
	}

	b.ParseObservation()
	b.ParseUnits()
	b.ParseOrders()
	b.DetectEnemyRace()
	b.FindClusters()
}
