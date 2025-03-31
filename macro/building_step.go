package macro

import (
	"github.com/NatoBoram/BlackCompany/bot"
	"github.com/NatoBoram/BlackCompany/filter"
	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/api"
)

func buildingStep(name string, buildingId api.UnitTypeID, abilityId api.AbilityID, quantity int, requirements ...api.UnitTypeID) *bot.BuildStep {
	return &bot.BuildStep{
		Name: name,
		Predicate: func(b *bot.Bot) bool {
			for _, requirement := range requirements {
				if b.Units.My.OfType(requirement).Empty() {
					return false
				}
			}

			if !b.CanBuy(abilityId) {
				return false
			}

			buildings := b.Units.My.OfType(buildingId)
			ordered := b.FindWorkers().Filter(filter.IsOrderedTo(abilityId))
			inProgress := buildings.Filter(filter.IsInProgress)
			notStarted := ordered.Len() - inProgress.Len()
			return buildings.Len()+notStarted < quantity
		},

		Execute: func(b *bot.Bot) {
			build(b, name, buildingId, abilityId, scl.S5x3)
		},

		Next: func(b *bot.Bot) bool {
			buildings := b.Units.My.OfType(buildingId)
			ordered := b.FindWorkers().Filter(filter.IsOrderedTo(abilityId))
			inProgress := buildings.Filter(filter.IsInProgress)
			notStarted := ordered.Len() - inProgress.Len()
			return buildings.Len()+notStarted >= quantity
		},
	}
}
