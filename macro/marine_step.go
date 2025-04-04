package macro

import (
	"github.com/NatoBoram/BlackCompany/bot"
	"github.com/NatoBoram/BlackCompany/filter"
	"github.com/NatoBoram/BlackCompany/log"
	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/enums/ability"
	"github.com/aiseeq/s2l/protocol/enums/terran"
)

var marineStep = bot.BuildStep{
	Name: "Train Marine",
	Predicate: func(b *bot.Bot) bool {
		// Waiting for the minerals to catch up with the supply slows down the
		// army production to allow for the rest of the build order to execute.
		return b.CanBuy(ability.Train_Marine) && b.Minerals > b.FoodUsed
	},

	Execute: func(b *bot.Bot) {
		barracks := b.Units.My.OfType(terran.Barracks).Filter(scl.Ready, scl.Ground, scl.Idle, filter.IsNotTag(b.State.BuildingForAddOn))
		if barracks.Empty() {
			return
		}

		for _, barrack := range barracks {
			amount := deductMarines(b, barrack)
			if amount == 0 {
				break
			}

			if rally := rallyPoint(b); rally != nil {
				barrack.CommandPos(ability.Rally_Building, rally)
			}

			if amount == 1 {
				barrack.CommandQueue(ability.Train_Marine)
				log.Info("Training one marine at %v", barrack.Point())
			}

			if amount == 2 {
				barrack.CommandQueue(ability.Train_Marine)
				barrack.CommandQueue(ability.Train_Marine)
				log.Info("Training two marines at %v", barrack.Point())
			}
		}
	},

	Next: func(b *bot.Bot) bool {
		return true
	},
}
