package main

import (
	"log"
	"math"
	"math/rand"

	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/enums/ability"
	"github.com/aiseeq/s2l/protocol/enums/terran"
)

func (b *Bot) BuildSupplyDepot() {
	if !b.CanBuy(ability.Build_SupplyDepot) {
		return
	}

	currSupply := b.Obs.PlayerCommon.FoodUsed
	maxSupply := b.Obs.PlayerCommon.FoodCap
	supplyLeft := maxSupply - currSupply

	townHalls := b.findTownHalls()
	if townHalls.Empty() {
		return
	}

	// Kinda assume that all command centers are building SCVs at the same time
	// even though they're not, but that's just to estimate the production per
	// bases.
	timeForScv := float64(BuildTimeSCV) / float64(townHalls.Len())
	scvDuringDepots := uint32(math.Ceil(float64(BuildTimeSupplyDepot) / timeForScv))

	depotsOrdered := b.findWorkers().Filter(IsOrderedTo(ability.Build_SupplyDepot)).Len() >= 1
	depotsInProgress := b.Units.My.OfType(terran.SupplyDepot).Filter(IsInProgress).Len() >= 1
	shouldBuildDepot := maxSupply < 200 && supplyLeft <= scvDuringDepots && !depotsOrdered && !depotsInProgress

	if !shouldBuildDepot {
		return
	}

	// From this point onward, it's relatively sure that we're going to build the
	// supply depot unless we're actively losing.
	randomTownHall := townHalls[rand.Intn(len(townHalls))]

	// Find a good position for the supply depot
	pos := b.whereToBuild(randomTownHall.Point(), scl.S2x2, terran.SupplyDepot, ability.Build_SupplyDepot)
	if pos == nil {
		log.Printf("No valid position found for supply depot")
		return
	}

	builder := b.findIdleOrGatheringWorkers().ClosestTo(pos)
	if builder == nil {
		log.Printf("No idle or gathering worker found to build supply depot")
		return
	}

	// Go back to the closest resource after building the supply depot
	if resource := b.findResourcesNearTownHalls(townHalls).ClosestTo(pos); resource != nil {
		log.Printf("Building supply depot at %v and queuing to gather at %v", *pos, resource.Point())

		builder.CommandPos(ability.Build_SupplyDepot, pos)
		builder.CommandTagQueue(ability.Smart, resource.Tag)

		if resource.IsMineral() {
			b.Miners.MineralForMiner[builder.Tag] = resource.Tag
		}

		if resource.IsGeyser() {
			b.Miners.GasForMiner[builder.Tag] = resource.Tag
		}
	} else {
		log.Printf("Building supply depot at %v", *pos)
		builder.CommandPos(ability.Build_SupplyDepot, pos)
	}

	b.DeductResources(ability.Build_SupplyDepot)
}
