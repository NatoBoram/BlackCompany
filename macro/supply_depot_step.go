package macro

import (
	"math"
	"math/rand"
	"slices"

	"github.com/NatoBoram/BlackCompany/bot"
	"github.com/NatoBoram/BlackCompany/filter"
	"github.com/NatoBoram/BlackCompany/log"
	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/enums/ability"
	"github.com/aiseeq/s2l/protocol/enums/terran"
)

// supplyDepotStep builds supply depots.
//
// TODO: Continuously make one supply depot after the other. Don't wait until
// we're supply blocked.
var supplyDepotStep = bot.BuildStep{
	Name: "Supply Depot",
	Predicate: func(b *bot.Bot) bool {
		if !b.CanBuy(ability.Build_SupplyDepot) {
			return false
		}

		currSupply := b.Obs.PlayerCommon.FoodUsed
		maxSupply := b.Obs.PlayerCommon.FoodCap
		supplyLeft := maxSupply - currSupply

		if maxSupply >= 200 {
			return false
		}

		townHalls := b.FindTownHalls()
		production := b.FindProductionStructures()

		structures := slices.Concat(townHalls, production)
		if structures.Empty() {
			return false
		}

		// Calculate how much supply we'll use during depot construction
		timeForScv := float64(bot.BuildTimeSCV) / float64(structures.Len())
		scvDuringDepots := uint32(math.Ceil(float64(bot.BuildTimeSupplyDepot) / timeForScv))

		// Don't build if we have enough supply or if already building
		depotsOrdered := b.FindWorkers().Filter(filter.IsOrderedTo(ability.Build_SupplyDepot)).Len() >= 1
		depotsInProgress := b.Units.My.OfType(terran.SupplyDepot).Filter(filter.IsInProgress).Len() >= 1

		return supplyLeft <= scvDuringDepots && !depotsOrdered && !depotsInProgress
	},

	Execute: func(b *bot.Bot) {
		townHalls := b.FindTownHalls().Filter(filter.IsCcAtExpansion(b.State.CcForExp))
		if townHalls.Empty() {
			return
		}

		randomTownHall := townHalls[rand.Intn(len(townHalls))]

		// Find a good position for the supply depot
		pos := b.WhereToBuild(randomTownHall.Point(), scl.S2x2, terran.SupplyDepot, ability.Build_SupplyDepot)
		if pos == nil {
			return
		}

		builder := b.FindIdleOrGatheringWorkers().ClosestTo(pos)
		if builder == nil {
			return
		}

		// Go back to the closest resource after building
		if resource := b.FindResourcesNearTownHalls(townHalls).ClosestTo(pos); resource != nil {
			log.Info("Building supply depot at %v and queuing to gather at %v", *pos, resource.Point())

			builder.CommandPos(ability.Build_SupplyDepot, pos)
			builder.CommandTagQueue(ability.Smart, resource.Tag)

			if resource.IsMineral() {
				b.Miners.MineralForMiner[builder.Tag] = resource.Tag
			}

			if resource.IsGeyser() {
				b.Miners.GasForMiner[builder.Tag] = resource.Tag
			}
		} else {
			log.Info("Building supply depot at %v", *pos)
			builder.CommandPos(ability.Build_SupplyDepot, pos)
		}

		b.DeductResources(ability.Build_SupplyDepot)
	},

	Next: func(b *bot.Bot) bool {
		return b.Units.My.OfType(terran.SupplyDepot).Len() >= 1 ||
			b.FindWorkers().Filter(filter.IsOrderedTo(ability.Build_SupplyDepot)).Exists()
	},
}
