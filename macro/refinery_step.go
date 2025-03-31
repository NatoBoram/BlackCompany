package macro

import (
	"math/rand"

	"github.com/NatoBoram/BlackCompany/bot"
	"github.com/NatoBoram/BlackCompany/filter"
	"github.com/NatoBoram/BlackCompany/log"
	"github.com/aiseeq/s2l/protocol/enums/ability"
	"github.com/aiseeq/s2l/protocol/enums/terran"
)

func refineryStep(quantity int) *bot.BuildStep {
	return &bot.BuildStep{
		Name: stepName("Refinery", quantity),
		Predicate: func(b *bot.Bot) bool {
			if !b.CanBuy(ability.Build_Refinery) {
				return false
			}

			if quantity == 0 {
				return true
			}

			refineries := b.Units.My.OfType(terran.Refinery, terran.RefineryRich)
			ordered := b.FindWorkers().Filter(filter.IsOrderedTo(ability.Build_Refinery))
			inProgress := refineries.Filter(filter.IsInProgress)
			notStarted := ordered.Len() - inProgress.Len()

			return refineries.Len()+notStarted < quantity
		},

		Execute: func(b *bot.Bot) {
			townHalls := b.FindTownHalls()
			if townHalls.Empty() {
				return
			}

			vespeneGeysers := b.FindVespeneGeysersNearTownHalls(townHalls)
			claimed := b.FindClaimedVespeneGeysersNearTownHalls(townHalls)
			buildable := vespeneGeysers.Filter(filter.NotCloserThan(1, claimed))
			if buildable.Empty() {
				return
			}

			ordered := b.FindWorkers().Filter(filter.IsOrderedTo(ability.Build_Refinery))
			unplanned := buildable.Filter(filter.IsNotOrderedOnTag(ability.Build_Refinery, ordered))
			if unplanned.Empty() {
				return
			}

			random := unplanned[rand.Intn(len(unplanned))]
			worker := b.FindIdleOrGatheringWorkers().ClosestTo(random)
			if worker == nil {
				return
			}

			log.Info("Building refinery at %v", random.Point())
			worker.CommandTag(ability.Build_Refinery, random.Tag)
			worker.CommandTagQueue(ability.Smart, random.Tag)
			b.DeductResources(ability.Build_Refinery)
			b.Miners.GasForMiner[worker.Tag] = random.Tag
		},

		Next: func(b *bot.Bot) bool {
			if quantity == 0 {
				return true
			}

			refineries := b.Units.My.OfType(terran.Refinery, terran.RefineryRich)
			ordered := b.FindWorkers().Filter(filter.IsOrderedTo(ability.Build_Refinery))
			inProgress := refineries.Filter(filter.IsInProgress)
			notStarted := ordered.Len() - inProgress.Len()
			return refineries.Len()+notStarted >= quantity
		},
	}
}
