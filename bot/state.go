package bot

import (
	"github.com/NatoBoram/BlackCompany/log"
	"github.com/NatoBoram/BlackCompany/quote"
	"github.com/NatoBoram/BlackCompany/wheel"
	"github.com/aiseeq/s2l/lib/point"
	"github.com/aiseeq/s2l/protocol/api"
)

type BotState struct {
	// CcForExp marks command centers that are reserved for new expansions.
	CcForExp CcForExp

	// CcForOrbitalCommand marks command centers that are reserved for upgrading
	// to orbital commands.
	CcForOrbitalCommand api.UnitTag

	// CcForPlanetaryFortress marks command centers that are reserved for
	// upgrading to planetary fortresses.
	CcForPlanetaryFortress api.UnitTag

	// BuildingForAddOn marks a barracks as reserved for building a reactor or
	// tech lab.
	BuildingForAddOn api.UnitTag

	// AttackWaves holds the groups of units that are used for attacking.
	AttackWaves AttackWaves

	// DetectedEnemyAirArmy saves whether the bot has seen any air units.
	DetectedEnemyAirArmy bool
}

func (b *Bot) InitState() {
	b.initCcForExp()
}

func (b *Bot) initCcForExp() {
	if b.State.CcForExp == nil {
		b.State.CcForExp = make(map[api.UnitTag]point.Point)
	}

	townHalls := b.FindTownHalls()
	if townHalls.Empty() {
		log.Warn("Couldn't initialize CcForExp because there are no town halls.")
		return
	}

	expansions := append(b.Locs.MyExps, b.Locs.MyStart)
	for _, expansion := range expansions {
		townHall := townHalls.ClosestTo(expansion)
		if townHall.IsCloserThan(1, expansion) {
			b.State.CcForExp[townHall.Tag] = expansion
		}
	}
}

func (b *Bot) detectEnemyAirArmy() {
	if b.State.DetectedEnemyAirArmy {
		return
	}

	army := b.FindEnemyAirArmy()
	if army.Exists() {
		log.Info("Detected enemy air army unit.")
		b.State.DetectedEnemyAirArmy = true

		message := wheel.RandomIn(quote.DetectedEnemyAirArmyQuotes)
		b.Actions.ChatSend(message, api.ActionChat_Broadcast)
	}
}
