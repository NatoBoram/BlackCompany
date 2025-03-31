package macro

import (
	"github.com/NatoBoram/BlackCompany/bot"
	"github.com/NatoBoram/BlackCompany/filter"
	"github.com/NatoBoram/BlackCompany/log"
	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/api"
)

func upgradeStep(name string, abilityId api.AbilityID, buildingId api.UnitTypeID) *bot.BuildStep {
	return &bot.BuildStep{
		Name: name,
		Predicate: func(b *bot.Bot) bool {
			if !b.CanBuy(abilityId) {
				return false
			}

			if b.Upgrades[abilityId] {
				return false
			}

			if b.Units.My.OfType(buildingId).Filter(filter.IsOrderedTo(abilityId)).Exists() {
				return false
			}

			return true
		},

		Execute: func(b *bot.Bot) {
			buildings := b.Units.My.OfType(buildingId).Filter(scl.Ready, scl.Idle)
			if buildings.Empty() {
				return
			}

			log.Info("Researching %s", name)
			buildings.First().Command(abilityId)
			b.DeductResources(abilityId)
		},

		Next: func(b *bot.Bot) bool {
			if b.Upgrades[abilityId] {
				return true
			}

			if b.Units.My.OfType(buildingId).Filter(filter.IsOrderedTo(abilityId)).Exists() {
				return true
			}

			return false
		},
	}
}
