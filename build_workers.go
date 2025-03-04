package main

import (
	"github.com/aiseeq/s2l/lib/scl"
	"github.com/aiseeq/s2l/protocol/enums/ability"
	"github.com/aiseeq/s2l/protocol/enums/terran"
)

func (b *Bot) BuildWorker() {
	ccs := b.Units.My.OfType(terran.CommandCenter, terran.OrbitalCommand, terran.PlanetaryFortress)
	workers := b.Units.My[terran.SCV]

	for _, cc := range ccs.Filter(scl.Ready, scl.Idle) {
		if !b.CanBuy(ability.Train_SCV) {
			break
		}

		saturated := b.IsCCSaturated(cc)
		if len(workers) < 80 && !saturated {
			cc.Command(ability.Train_SCV)
			b.DeductResources(ability.Train_SCV)
		}
	}
}

// IsMineralSaturated checks if the mineral fields near a command center are saturated with workers.
func (b *Bot) IsMineralSaturated(cc *scl.Unit) bool {
	// Mineral fields near the CC
	mineralFields := b.Units.Minerals.All().CloserThan(scl.ResourceSpreadDistance, cc)
	if mineralFields.Empty() {
		return true
	}

	// Count miners assigned to this CC
	minersAssigned := 0
	for _, scv := range b.Units.My[terran.SCV] {
		// Check if this SCV has the CC's tag as its assigned CC in the mining data
		if b.Miners.CCForMiner[scv.Tag] == cc.Tag {
			minersAssigned++
		}
	}

	// Optimal saturation is 2 workers per mineral field, max is 3
	optimalSaturation := mineralFields.Len() * 2
	return minersAssigned >= optimalSaturation
}

// IsGasSaturated checks if the refineries near a command center are fully
// saturated.
func (b *Bot) IsGasSaturated(cc *scl.Unit) bool {
	refineries := b.Units.My[terran.Refinery].CloserThan(scl.ResourceSpreadDistance, cc)
	if refineries.Empty() {
		return true
	}

	for _, refinery := range refineries {
		if refinery.AssignedHarvesters < refinery.IdealHarvesters {
			return false
		}
	}

	return true
}

// IsCCSaturated checks if a command center is saturated with workers.
func (b *Bot) IsCCSaturated(cc *scl.Unit) bool {
	return b.IsMineralSaturated(cc) && b.IsGasSaturated(cc)
}

func getTownHalls(b *Bot) []*scl.Unit {
	var cc []*scl.Unit
	cc = append(cc, b.Units.My[terran.CommandCenter]...)
	cc = append(cc, b.Units.My[terran.OrbitalCommand]...)
	cc = append(cc, b.Units.My[terran.PlanetaryFortress]...)
	return cc
}
