package macro

import (
	"github.com/NatoBoram/BlackCompany/bot"
	"github.com/NatoBoram/BlackCompany/filter"
	"github.com/aiseeq/s2l/protocol/enums/ability"
	"github.com/aiseeq/s2l/protocol/enums/terran"
)

func expandStep(quantity int) *bot.BuildStep {
	return &bot.BuildStep{
		Name: "Expand",
		Predicate: func(b *bot.Bot) bool {
			if !b.CanBuy(ability.Build_CommandCenter) {
				return false
			}

			if quantity != 0 &&
				b.FindTownHalls().Filter(filter.IsCcAtExpansion(b.State.CcForExp)).Len() >= quantity {
				return false
			}

			return true
		},

		Execute: func(b *bot.Bot) {
		},

		Next: func(b *bot.Bot) bool {
			return b.Units.My.OfType(terran.CommandCenter).Len() >= quantity
		},
	}
}
