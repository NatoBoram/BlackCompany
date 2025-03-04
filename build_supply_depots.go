package main

import (
	"math"
	"math/rand"

	"github.com/aiseeq/s2l/lib/point"
	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/enums/ability"
	"github.com/aiseeq/s2l/protocol/enums/terran"
)

func filterInProgress(unit *scl.Unit) bool {
	return unit.BuildProgress < 1
}

func filterGatheringOrIdle(unit *scl.Unit) bool {
	if len(unit.Orders) == 0 {
		return true
	}
	return unit.Orders[0].AbilityId == ability.Harvest_Gather_SCV
}

func filterBuildSupplyDepotOrder(unit *scl.Unit) bool {
	if len(unit.Orders) == 0 {
		return false
	}

	for _, order := range unit.Orders {
		if order.AbilityId == ability.Build_SupplyDepot {
			return true
		}
	}

	return false
}

func (b *Bot) BuildSupplyDepot() {
	currSupply := b.Obs.PlayerCommon.FoodUsed
	maxSupply := b.Obs.PlayerCommon.FoodCap
	supplyLeft := maxSupply - currSupply

	ccs := b.Units.My.OfType(terran.CommandCenter, terran.OrbitalCommand, terran.PlanetaryFortress)
	if ccs.Empty() {
		return
	}

	// Kinda assume that all command centers are building SCVs at the same time
	// even though they're not, but that's just to estimate the production per
	// bases.
	timeForScv := float64(BuildTimeSCV) / float64(ccs.Len())
	scvDuringDepots := uint32(math.Ceil(float64(BuildTimeSupplyDepot) / timeForScv))

	depotsOrdered := b.Units.My.OfType(terran.SCV).Filter(filterBuildSupplyDepotOrder).Len()
	depotsInProgress := b.Units.My.OfType(terran.SupplyDepot).Filter(filterInProgress).Len()
	shouldBuildDepot := maxSupply < 200 && supplyLeft <= scvDuringDepots && depotsOrdered+depotsInProgress == 0

	if !shouldBuildDepot || !b.CanBuy(ability.Build_SupplyDepot) {
		return
	}

	workers := b.Units.My[terran.SCV].Filter(filterGatheringOrIdle)
	if len(workers) == 0 {
		return
	}

	builder := workers.First()
	if builder == nil {
		return
	}

	cc := ccs[rand.Intn(len(ccs))]

	// Find a good position for the supply depot
	pos := b.findSupplyDepotPlacement(cc)
	if pos == nil {
		return
	}

	builder.CommandPosQueue(ability.Build_SupplyDepot, pos)
	b.DeductResources(ability.Build_SupplyDepot)
}

// findSupplyDepotPlacement finds a valid position to place a supply depot near the given command center
func (b *Bot) findSupplyDepotPlacement(cc *scl.Unit) *point.Point {
	return nil
}
