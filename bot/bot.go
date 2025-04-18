package bot

import (
	"github.com/NatoBoram/BlackCompany/adapter"
	"github.com/NatoBoram/BlackCompany/log"
	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/api"
)

// Bot holds the state of the bot.
type Bot struct {
	*scl.Bot

	miningInitialized bool

	State BotState
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
}

// Observe fetches the current observation from the game.
func (b *Bot) Observe() {
	o, err := b.Client.Observation(api.RequestObservation{})
	if err != nil {
		log.Info("Failed to observe: %v", err)
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
		log.Info("Info is nil")
		return
	}
	if b.Obs == nil {
		log.Info("Observation is nil")
		return
	}
	if b.Obs.RawData == nil {
		log.Info("RawData is nil")
		return
	}
	if b.Obs.RawData.MapState == nil {
		log.Info("MapState is nil")
		return
	}

	b.ParseObservation()
	b.ParseUnits()
	b.ParseOrders()
	b.DetectEnemyRace()

	if !b.miningInitialized {
		townHalls := b.FindTownHalls()
		resources := b.FindResourcesNearTownHalls(townHalls)
		turrets := b.findTurretsNearResourcesNearTownHalls(resources)
		b.InitMining(adapter.ToPoints(turrets))
		b.miningInitialized = true
	} else {
		b.acknowledgeMiners()
	}

	b.FindClusters() // Not used yet

	b.detectEnemyAirArmy()
}
