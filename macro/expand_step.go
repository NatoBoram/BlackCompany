package macro

import (
	"github.com/NatoBoram/BlackCompany/bot"
	"github.com/NatoBoram/BlackCompany/log"
	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/enums/ability"
	"github.com/aiseeq/s2l/protocol/enums/terran"
)

func expandStep(quantity int) *bot.BuildStep {
	return &bot.BuildStep{
		Name: "Expand",
		Predicate: func(b *bot.Bot) bool {
			if !b.ShouldExpand() {
				return false
			}

			if quantity != 0 && b.FindTownHalls().Len() >= quantity {
				return false
			}

			return true
		},

		Execute: func(b *bot.Bot) {
			expansions := b.FindExpansionLocations()
			if expansions.Empty() {
				return
			}

			// Assign command centers to expansions
			available := b.FindAvailableCommandCenters()
			for i, cc := range available {
				if i >= len(expansions) {
					break
				}
				expansion := expansions[i]

				log.Debug("Assigning a town hall to expansion %s", expansion)
				b.State.CcForExp[cc.Tag] = expansion
			}
			if available.Exists() {
				return
			}

			// No available command centers, so let's build one
			expansion := expansions[0]

			worker := b.FindIdleOrGatheringWorkers().ClosestTo(expansion)
			if worker == nil {
				return
			}

			townHalls := b.FindTownHalls()

			// We're fucked, let's build one asap
			if townHalls.Empty() {
				log.Info("No town halls found, building Command Center at expansion %s", expansion)
				location := b.WhereToBuild(expansion, scl.S5x5, terran.CommandCenter, ability.Build_CommandCenter)
				worker.CommandPos(ability.Build_CommandCenter, location)
				b.DeductResources(ability.Build_CommandCenter)
				return
			}

			nearestTownHall := townHalls.ClosestTo(worker)
			towards := nearestTownHall.Point().Towards(expansion, nearestTownHall.SightRange())
			location := b.WhereToBuild(towards, scl.S5x5, terran.CommandCenter, ability.Build_CommandCenter)

			// So do I build it there or near worker then fly it over?
			if b.IsFlyingFaster(worker, location, expansion) {
				log.Info("Building Command Center at base %s to fly to expansion %s", location, expansion)

				worker.CommandPos(ability.Build_CommandCenter, location)
				b.DeductResources(ability.Build_CommandCenter)

				closestMineralField := b.Units.Minerals.All().ClosestTo(expansion)
				worker.CommandTagQueue(ability.Smart, closestMineralField.Tag)

				return
			}

			log.Info("Expanding to %s", expansion)

			worker.CommandPos(ability.Build_CommandCenter, expansion)
			b.DeductResources(ability.Build_CommandCenter)

			closestMineralField := b.Units.Minerals.All().ClosestTo(expansion)
			worker.CommandTagQueue(ability.Smart, closestMineralField.Tag)
		},

		Next: func(b *bot.Bot) bool {
			return quantity == 0 || b.FindTownHalls().Len() >= quantity
		},
	}
}
